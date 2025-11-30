package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"snitchdns-tf/internal/client"
)

// Create implements the resource create logic
func (r *RecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create timeout context
	createTimeout, diags := data.Timeouts.Create(ctx, 5*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// Convert data map to map[string]interface{}
	dataMap := make(map[string]interface{})
	for key, value := range data.Data.Elements() {
		strVal, ok := value.(types.String)
		if ok {
			dataMap[key] = strVal.ValueString()
		}
	}

	// Convert conditional_data map if present
	conditionalDataMap := make(map[string]interface{})
	if !data.ConditionalData.IsNull() {
		for key, value := range data.ConditionalData.Elements() {
			strVal, ok := value.(types.String)
			if ok {
				conditionalDataMap[key] = strVal.ValueString()
			}
		}
	}

	// Create record via API
	createReq := client.CreateRecordRequest{
		Active:           data.Active.ValueBool(),
		Class:            data.Class.ValueString(),
		Type:             data.Type.ValueString(),
		TTL:              int(data.TTL.ValueInt64()),
		Data:             dataMap,
		IsConditional:    data.IsConditional.ValueBoolPointer() != nil && data.IsConditional.ValueBool(),
		ConditionalCount: int(data.ConditionalCount.ValueInt64()),
		ConditionalLimit: int(data.ConditionalLimit.ValueInt64()),
		ConditionalReset: data.ConditionalReset.ValueBoolPointer() != nil && data.ConditionalReset.ValueBool(),
		ConditionalData:  conditionalDataMap,
	}

	record, err := r.client.CreateRecord(data.ZoneID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating record",
			fmt.Sprintf("Could not create record: %s", err),
		)
		return
	}

	// Update data model from API response
	data.ID = types.StringValue(fmt.Sprintf("%d", record.ID))
	data.ZoneID = types.StringValue(fmt.Sprintf("%d", record.ZoneID))
	data.Active = types.BoolValue(record.Active)
	data.Class = types.StringValue(record.Class)
	data.Type = types.StringValue(record.Type)
	data.TTL = types.Int64Value(int64(record.TTL))
	data.IsConditional = types.BoolValue(record.IsConditional)
	data.ConditionalCount = types.Int64Value(int64(record.ConditionalCount))
	data.ConditionalLimit = types.Int64Value(int64(record.ConditionalLimit))
	data.ConditionalReset = types.BoolValue(record.ConditionalReset)

	// Convert data map to types.Map
	dataElements := make(map[string]types.String)
	for key, value := range record.Data {
		dataElements[key] = types.StringValue(fmt.Sprintf("%v", value))
	}
	dataValue, diags := types.MapValueFrom(ctx, types.StringType, dataElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Data = dataValue

	// Convert conditional_data map if present
	if len(record.ConditionalData) > 0 {
		condDataElements := make(map[string]types.String)
		for key, value := range record.ConditionalData {
			condDataElements[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
		condDataValue, diags := types.MapValueFrom(ctx, types.StringType, condDataElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ConditionalData = condDataValue
	} else {
		data.ConditionalData = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements the resource read logic
func (r *RecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create timeout context
	readTimeout, diags := data.Timeouts.Read(ctx, 2*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	tflog.Debug(ctx, "Reading record", map[string]any{
		"zone_id":   data.ZoneID.ValueString(),
		"record_id": data.ID.ValueString(),
	})

	// Get record from API
	record, err := r.client.GetRecord(data.ZoneID.ValueString(), data.ID.ValueString())
	if err != nil {
		// Check if this is a 404 - resource was deleted outside Terraform
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			tflog.Warn(ctx, "Record not found, removing from state", map[string]any{
				"zone_id":   data.ZoneID.ValueString(),
				"record_id": data.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading DNS Record",
			fmt.Sprintf("Could not read record ID %s in zone %s: %s",
				data.ID.ValueString(), data.ZoneID.ValueString(), err),
		)
		return
	}

	// Update data model from API response
	data.ID = types.StringValue(fmt.Sprintf("%d", record.ID))
	data.ZoneID = types.StringValue(fmt.Sprintf("%d", record.ZoneID))
	data.Active = types.BoolValue(record.Active)
	data.Class = types.StringValue(record.Class)
	data.Type = types.StringValue(record.Type)
	data.TTL = types.Int64Value(int64(record.TTL))
	data.IsConditional = types.BoolValue(record.IsConditional)
	data.ConditionalCount = types.Int64Value(int64(record.ConditionalCount))
	data.ConditionalLimit = types.Int64Value(int64(record.ConditionalLimit))
	data.ConditionalReset = types.BoolValue(record.ConditionalReset)

	// Convert data map to types.Map
	dataElements := make(map[string]types.String)
	for key, value := range record.Data {
		dataElements[key] = types.StringValue(fmt.Sprintf("%v", value))
	}
	dataValue, diags := types.MapValueFrom(ctx, types.StringType, dataElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Data = dataValue

	// Convert conditional_data map if present
	if len(record.ConditionalData) > 0 {
		condDataElements := make(map[string]types.String)
		for key, value := range record.ConditionalData {
			condDataElements[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
		condDataValue, diags := types.MapValueFrom(ctx, types.StringType, condDataElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ConditionalData = condDataValue
	} else {
		data.ConditionalData = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements the resource update logic
func (r *RecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create timeout context
	updateTimeout, diags := data.Timeouts.Update(ctx, 2*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	// Convert data map to map[string]interface{}
	dataMap := make(map[string]interface{})
	for key, value := range data.Data.Elements() {
		strVal, ok := value.(types.String)
		if ok {
			dataMap[key] = strVal.ValueString()
		}
	}

	// Convert conditional_data map if present
	var conditionalDataMap map[string]interface{}
	if !data.ConditionalData.IsNull() {
		conditionalDataMap = make(map[string]interface{})
		for key, value := range data.ConditionalData.Elements() {
			strVal, ok := value.(types.String)
			if ok {
				conditionalDataMap[key] = strVal.ValueString()
			}
		}
	}

	// Update record via API
	active := data.Active.ValueBool()
	cls := data.Class.ValueString()
	typ := data.Type.ValueString()
	ttl := int(data.TTL.ValueInt64())
	isConditional := data.IsConditional.ValueBool()
	conditionalCount := int(data.ConditionalCount.ValueInt64())
	conditionalLimit := int(data.ConditionalLimit.ValueInt64())
	conditionalReset := data.ConditionalReset.ValueBool()

	updateReq := client.UpdateRecordRequest{
		Active:           &active,
		Class:            &cls,
		Type:             &typ,
		TTL:              &ttl,
		Data:             dataMap,
		IsConditional:    &isConditional,
		ConditionalCount: &conditionalCount,
		ConditionalLimit: &conditionalLimit,
		ConditionalReset: &conditionalReset,
		ConditionalData:  conditionalDataMap,
	}

	record, err := r.client.UpdateRecord(data.ZoneID.ValueString(), data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating record",
			fmt.Sprintf("Could not update record ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	// Update data model from API response
	data.ID = types.StringValue(fmt.Sprintf("%d", record.ID))
	data.ZoneID = types.StringValue(fmt.Sprintf("%d", record.ZoneID))
	data.Active = types.BoolValue(record.Active)
	data.Class = types.StringValue(record.Class)
	data.Type = types.StringValue(record.Type)
	data.TTL = types.Int64Value(int64(record.TTL))
	data.IsConditional = types.BoolValue(record.IsConditional)
	data.ConditionalCount = types.Int64Value(int64(record.ConditionalCount))
	data.ConditionalLimit = types.Int64Value(int64(record.ConditionalLimit))
	data.ConditionalReset = types.BoolValue(record.ConditionalReset)

	// Convert data map to types.Map
	dataElements := make(map[string]types.String)
	for key, value := range record.Data {
		dataElements[key] = types.StringValue(fmt.Sprintf("%v", value))
	}
	dataValue, diags := types.MapValueFrom(ctx, types.StringType, dataElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Data = dataValue

	// Convert conditional_data map if present
	if len(record.ConditionalData) > 0 {
		condDataElements := make(map[string]types.String)
		for key, value := range record.ConditionalData {
			condDataElements[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
		condDataValue, diags := types.MapValueFrom(ctx, types.StringType, condDataElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ConditionalData = condDataValue
	} else {
		data.ConditionalData = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements the resource delete logic
func (r *RecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create timeout context
	deleteTimeout, diags := data.Timeouts.Delete(ctx, 3*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	// Delete record via API
	err := r.client.DeleteRecord(data.ZoneID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting record",
			fmt.Sprintf("Could not delete record ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
}

// ImportState implements the resource import logic
func (r *RecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: "zone_id:record_id"
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected import ID format 'zone_id:record_id', got: %s", req.ID),
		)
		return
	}

	zoneID := parts[0]
	recordID := parts[1]

	// Validate they are numeric
	if _, err := strconv.Atoi(zoneID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid zone ID",
			fmt.Sprintf("Zone ID must be numeric, got: %s", zoneID),
		)
		return
	}
	if _, err := strconv.Atoi(recordID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid record ID",
			fmt.Sprintf("Record ID must be numeric, got: %s", recordID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), recordID)...)
}
