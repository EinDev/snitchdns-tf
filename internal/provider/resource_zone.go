package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
var _ resource.Resource = &ZoneResource{}
var _ resource.ResourceWithImportState = &ZoneResource{}

func NewZoneResource() resource.Resource {
	return &ZoneResource{}
}

// ZoneResource defines the resource implementation.
type ZoneResource struct {
	client *client.Client
}

// ZoneResourceModel describes the resource data model.
type ZoneResourceModel struct {
	ID         types.String   `tfsdk:"id"`
	UserID     types.Int64    `tfsdk:"user_id"`
	Domain     types.String   `tfsdk:"domain"`
	Active     types.Bool     `tfsdk:"active"`
	CatchAll   types.Bool     `tfsdk:"catch_all"`
	Forwarding types.Bool     `tfsdk:"forwarding"`
	Regex      types.Bool     `tfsdk:"regex"`
	Master     types.Bool     `tfsdk:"master"`
	Tags       types.List     `tfsdk:"tags"`
	CreatedAt  types.String   `tfsdk:"created_at"`
	UpdatedAt  types.String   `tfsdk:"updated_at"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

func (r *ZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (r *ZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNS zone in SnitchDNS. Zones are containers for DNS records and can be configured with various options like catch-all, forwarding, and regex matching.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the zone. Assigned by the API upon creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the user who owns this zone. Automatically set by the API based on authentication.",
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain name for this zone (e.g., `example.com`). This is the base domain that will be used for DNS queries.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the zone is active and will respond to DNS queries. Set to `false` to disable the zone without deleting it.",
				Required:            true,
			},
			"catch_all": schema.BoolAttribute{
				MarkdownDescription: "Enable catch-all DNS queries for this zone. When enabled, the zone will respond to queries for any subdomain, even if no specific record exists.",
				Required:            true,
			},
			"forwarding": schema.BoolAttribute{
				MarkdownDescription: "Enable DNS forwarding to upstream DNS servers. When enabled, unmatched queries will be forwarded to a configured upstream resolver.",
				Required:            true,
			},
			"regex": schema.BoolAttribute{
				MarkdownDescription: "Use regular expression matching for the domain name. When enabled, the domain field can contain a regex pattern instead of a literal domain.",
				Required:            true,
			},
			"master": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this is a master zone. Master zones have special privileges and cannot be modified via the API. This is set automatically during creation.",
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of tags to organize and categorize zones. Tags can be used for filtering and grouping zones.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the zone was created in RFC3339 format.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the zone was last updated in RFC3339 format.",
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

func (r *ZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// CRUD methods are implemented in resource_zone_impl.go
