package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"snitchdns-tf/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RecordResource{}
var _ resource.ResourceWithImportState = &RecordResource{}

func NewRecordResource() resource.Resource {
	return &RecordResource{}
}

// RecordResource defines the resource implementation.
type RecordResource struct {
	client *client.Client
}

// RecordResourceModel describes the resource data model.
type RecordResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	ZoneID           types.String   `tfsdk:"zone_id"`
	Active           types.Bool     `tfsdk:"active"`
	Class            types.String   `tfsdk:"cls"`
	Type             types.String   `tfsdk:"type"`
	TTL              types.Int64    `tfsdk:"ttl"`
	Data             types.Map      `tfsdk:"data"`
	IsConditional    types.Bool     `tfsdk:"is_conditional"`
	ConditionalCount types.Int64    `tfsdk:"conditional_count"`
	ConditionalLimit types.Int64    `tfsdk:"conditional_limit"`
	ConditionalReset types.Bool     `tfsdk:"conditional_reset"`
	ConditionalData  types.Map      `tfsdk:"conditional_data"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

func (r *RecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

func (r *RecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNS record within a SnitchDNS zone. Records define the actual DNS responses for queries and support all standard DNS record types (A, AAAA, CNAME, MX, TXT, etc.) as well as conditional responses.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the DNS record. Assigned by the API upon creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the zone this record belongs to. Records are always associated with a specific zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"active": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the record is active and will respond to DNS queries. Set to `false` to temporarily disable without deleting.",
			},
			"cls": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DNS class for the record. Typically `IN` (Internet), but can also be `CH` (Chaos) or `HS` (Hesiod).",
				Validators: []validator.String{
					stringvalidator.OneOf("IN", "CH", "HS"),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DNS record type. Supported types: A, AAAA, AFSDB, CAA, CNAME, DNAME, HINFO, MX, NAPTR, NS, PTR, RP, SOA, SPF, SRV, SSHFP, TSIG, TXT.",
				Validators: []validator.String{
					stringvalidator.OneOf("A", "AAAA", "AFSDB", "CAA", "CNAME", "DNAME", "HINFO", "MX", "NAPTR", "NS", "PTR", "RP", "SOA", "SPF", "SRV", "SSHFP", "TSIG", "TXT"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ttl": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Time to live in seconds. Determines how long DNS resolvers should cache this record. Typical values range from 60 (1 minute) to 86400 (1 day).",
				Validators: []validator.Int64{
					int64validator.Between(1, 2147483647), // Max 32-bit int for TTL
				},
			},
			"data": schema.MapAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Record-specific data as key-value pairs. The required fields depend on the record type. For A records: `{address = \"192.168.1.1\"}`. For CNAME: `{name = \"target.example.com\"}`. For MX: `{priority = \"10\", hostname = \"mail.example.com\"}`.",
			},
			"is_conditional": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enable conditional responses based on query count. When enabled, the record can return different data based on how many times it has been queried.",
			},
			"conditional_count": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Current query count for conditional logic. Automatically incremented by SnitchDNS when the record is queried.",
			},
			"conditional_limit": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Query limit for conditional responses. When `conditional_count` reaches this limit, the conditional behavior triggers.",
			},
			"conditional_reset": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reset the query counter when the limit is reached. If `true`, the counter resets to 0; if `false`, it remains at the limit.",
			},
			"conditional_data": schema.MapAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Alternative data to return when conditional limit is reached. Uses the same format as the `data` attribute.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *RecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}
