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
func (r *ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ZoneResourceModel

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

	tflog.Debug(ctx, "Creating zone", map[string]any{
		"domain": data.Domain.ValueString(),
	})

	// Convert tags list to comma-separated string
	var tags []string
	if !data.Tags.IsNull() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	tagsStr := strings.Join(tags, ",")

	// Create zone via API
	createReq := client.CreateZoneRequest{
		Domain:     data.Domain.ValueString(),
		Active:     data.Active.ValueBool(),
		CatchAll:   data.CatchAll.ValueBool(),
		Forwarding: data.Forwarding.ValueBool(),
		Regex:      data.Regex.ValueBool(),
		Master:     false, // Always false for user-created zones
		Tags:       tagsStr,
	}

	zone, err := r.client.CreateZone(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating zone",
			fmt.Sprintf("Could not create zone: %s", err),
		)
		return
	}

	// Map response to data model
	data.ID = types.StringValue(strconv.Itoa(zone.ID))
	data.UserID = types.Int64Value(int64(zone.UserID))
	data.Master = types.BoolValue(zone.Master)
	data.CreatedAt = types.StringValue(zone.CreatedAt)
	data.UpdatedAt = types.StringValue(zone.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements the resource read logic
func (r *ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ZoneResourceModel

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

	// Get zone from API
	zone, err := r.client.GetZone(data.ID.ValueString())
	if err != nil {
		// Check if this is a 404 - resource was deleted outside Terraform
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			tflog.Warn(ctx, "Zone not found, removing from state", map[string]any{
				"id": data.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading zone",
			fmt.Sprintf("Could not read zone ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	// Update data model from API response
	data.UserID = types.Int64Value(int64(zone.UserID))
	data.Domain = types.StringValue(zone.Domain)
	data.Active = types.BoolValue(zone.Active)
	data.CatchAll = types.BoolValue(zone.CatchAll)
	data.Forwarding = types.BoolValue(zone.Forwarding)
	data.Regex = types.BoolValue(zone.Regex)
	data.Master = types.BoolValue(zone.Master)
	data.CreatedAt = types.StringValue(zone.CreatedAt)
	data.UpdatedAt = types.StringValue(zone.UpdatedAt)

	// Convert tags array to list
	if len(zone.Tags) > 0 {
		tagsValue, diags := types.ListValueFrom(ctx, types.StringType, zone.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements the resource update logic
func (r *ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ZoneResourceModel

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

	// Convert tags list to comma-separated string
	var tags []string
	if !data.Tags.IsNull() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	tagsStr := strings.Join(tags, ",")

	// Update zone via API
	domain := data.Domain.ValueString()
	active := data.Active.ValueBool()
	catchAll := data.CatchAll.ValueBool()
	forwarding := data.Forwarding.ValueBool()
	regex := data.Regex.ValueBool()

	updateReq := client.UpdateZoneRequest{
		Domain:     &domain,
		Active:     &active,
		CatchAll:   &catchAll,
		Forwarding: &forwarding,
		Regex:      &regex,
		Tags:       &tagsStr,
	}

	zone, err := r.client.UpdateZone(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating zone",
			fmt.Sprintf("Could not update zone ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}

	// Update data model from API response
	data.ID = types.StringValue(fmt.Sprintf("%d", zone.ID))
	data.UserID = types.Int64Value(int64(zone.UserID))
	data.Domain = types.StringValue(zone.Domain)
	data.Active = types.BoolValue(zone.Active)
	data.CatchAll = types.BoolValue(zone.CatchAll)
	data.Forwarding = types.BoolValue(zone.Forwarding)
	data.Regex = types.BoolValue(zone.Regex)
	data.Master = types.BoolValue(zone.Master)
	data.CreatedAt = types.StringValue(zone.CreatedAt)
	data.UpdatedAt = types.StringValue(zone.UpdatedAt)

	// Convert tags array to list
	if len(zone.Tags) > 0 {
		tagsValue, diags := types.ListValueFrom(ctx, types.StringType, zone.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements the resource delete logic
func (r *ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ZoneResourceModel

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

	// Delete zone via API
	err := r.client.DeleteZone(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting zone",
			fmt.Sprintf("Could not delete zone ID %s: %s", data.ID.ValueString(), err),
		)
		return
	}
}

// ImportState implements the resource import logic
func (r *ZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID from the import request
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
