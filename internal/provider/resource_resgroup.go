package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure resgroupResource satisfies the resource.Resource interface.
var _ resource.Resource = &resgroupResource{}

// NewResgroupResource returns a new reservation group resource.
func NewResgroupResource() resource.Resource {
	return &resgroupResource{}
}

// resgroupResource manages a CloudLab reservation group.
type resgroupResource struct {
	client *Client
}

// resgroupNodeModel maps a node type entry in a reservation group.
type resgroupNodeModel struct {
	NodeType  types.String `tfsdk:"node_type"`
	Aggregate types.String `tfsdk:"aggregate"`
	Count     types.Int64  `tfsdk:"count"`
}

// resgroupResourceModel maps the resource schema data.
type resgroupResourceModel struct {
	ID        types.String        `tfsdk:"id"`
	Project   types.String        `tfsdk:"project"`
	Reason    types.String        `tfsdk:"reason"`
	StartAt   types.String        `tfsdk:"start_at"`
	ExpiresAt types.String        `tfsdk:"expires_at"`
	NodeTypes []resgroupNodeModel `tfsdk:"node_types"`
	Creator   types.String        `tfsdk:"creator"`
	Status    types.String        `tfsdk:"status"`
	CreatedAt types.String        `tfsdk:"created_at"`
}

// Metadata returns the resource type name.
func (r *resgroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resgroup"
}

// Schema defines the schema for the resource.
func (r *resgroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a CloudLab reservation group. Reservation groups allow you to pre-reserve " +
			"specific hardware resources on CloudLab for a defined time window, ensuring availability " +
			"when you need to run experiments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier (UUID) of the reservation group assigned by CloudLab.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project": schema.StringAttribute{
				Description: "The CloudLab project for this reservation group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reason": schema.StringAttribute{
				Description: "A description of why you need to reserve these resources.",
				Required:    true,
			},
			"start_at": schema.StringAttribute{
				Description: "The time the reservation should start (RFC3339 format). " +
					"If omitted, the reservation starts immediately.",
				Optional: true,
			},
			"expires_at": schema.StringAttribute{
				Description: "The time the reservation expires (RFC3339 format).",
				Optional:    true,
				Computed:    true,
			},
			"node_types": schema.ListNestedAttribute{
				Description: "The list of node types and counts to reserve.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node_type": schema.StringAttribute{
							Description: "The hardware node type to reserve (e.g., xl170, m400).",
							Required:    true,
						},
						"aggregate": schema.StringAttribute{
							Description: "The CloudLab site/aggregate (e.g., utah.cloudlab.us).",
							Required:    true,
						},
						"count": schema.Int64Attribute{
							Description: "The number of nodes of this type to reserve.",
							Required:    true,
						},
					},
				},
			},
			"creator": schema.StringAttribute{
				Description: "The CloudLab username who created the reservation group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "The current status of the reservation group.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the reservation group was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure sets the provider-configured client on the resource.
func (r *resgroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *provider.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the reservation group and sets the initial Terraform state.
func (r *resgroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resgroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &ResgroupCreateRequest{
		Project: plan.Project.ValueString(),
		Reason:  plan.Reason.ValueString(),
	}

	if !plan.StartAt.IsNull() && !plan.StartAt.IsUnknown() {
		v := plan.StartAt.ValueString()
		createReq.StartAt = &v
	}
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		v := plan.ExpiresAt.ValueString()
		createReq.ExpiresAt = &v
	}

	for _, n := range plan.NodeTypes {
		createReq.NodeTypes = append(createReq.NodeTypes, ResgroupNode{
			NodeType:  n.NodeType.ValueString(),
			Aggregate: n.Aggregate.ValueString(),
			Count:     n.Count.ValueInt64(),
		})
	}

	tflog.Info(ctx, "Creating CloudLab reservation group", map[string]any{
		"project": createReq.Project,
	})

	rg, err := r.client.CreateResgroup(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Reservation Group", err.Error())
		return
	}

	plan = mapResgroupResponseToModel(rg, plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *resgroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resgroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rg, err := r.client.GetResgroup(state.ID.ValueString())
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			tflog.Warn(ctx, "Reservation group not found, removing from state", map[string]any{"id": state.ID.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Reservation Group", err.Error())
		return
	}

	state = mapResgroupResponseToModel(rg, state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates mutable reservation group attributes.
func (r *resgroupResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"CloudLab reservation groups cannot be updated in-place. Please delete and recreate the reservation group.",
	)
}

// Delete deletes the reservation group.
func (r *resgroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resgroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting CloudLab reservation group", map[string]any{"id": state.ID.ValueString()})

	if err := r.client.DeleteResgroup(state.ID.ValueString()); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError("Error Deleting Reservation Group", err.Error())
		return
	}
}

// mapResgroupResponseToModel maps an API response to the Terraform model.
func mapResgroupResponseToModel(rg *ResgroupResponse, model resgroupResourceModel) resgroupResourceModel {
	model.ID = types.StringValue(rg.ID)
	model.Project = types.StringValue(rg.Project)
	model.Reason = types.StringValue(rg.Reason)
	model.Creator = types.StringValue(rg.Creator)
	model.Status = types.StringValue(rg.Status)
	model.CreatedAt = types.StringValue(rg.CreatedAt)

	if rg.StartAt != nil {
		model.StartAt = types.StringValue(*rg.StartAt)
	} else {
		model.StartAt = types.StringNull()
	}

	if rg.ExpiresAt != nil {
		model.ExpiresAt = types.StringValue(*rg.ExpiresAt)
	} else {
		model.ExpiresAt = types.StringNull()
	}

	return model
}
