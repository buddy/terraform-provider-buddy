package source

// todo project
//import (
//	"buddy-terraform/buddy/util"
//	"context"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//func Project() *schema.Resource {
//	return &schema.Resource{
//		Description: "Get project by name or display_name\n\n" +
//			"Token scope required: `WORKSPACE`",
//		ReadContext: readContextProject,
//		Schema: map[string]*schema.Schema{
//			"id": {
//				Description: "The Terraform resource identifier for this item",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"domain": {
//				Description:  "The workspace's URL handle",
//				Type:         schema.TypeString,
//				Required:     true,
//				ValidateFunc: util.ValidateDomain,
//			},
//			"display_name": {
//				Description: "The project's display name",
//				Type:        schema.TypeString,
//				Computed:    true,
//				Optional:    true,
//				ExactlyOneOf: []string{
//					"display_name",
//					"name",
//				},
//			},
//			"name": {
//				Description: "The project's unique name ID",
//				Type:        schema.TypeString,
//				Computed:    true,
//				Optional:    true,
//				ExactlyOneOf: []string{
//					"display_name",
//					"name",
//				},
//			},
//			"html_url": {
//				Description: "The project's URL",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"status": {
//				Description: "The project's status",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//		},
//	}
//}
//
//func readContextProject(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	var project *buddy.Project
//	var err error
//	domain := d.Get("domain").(string)
//	if name, ok := d.GetOk("name"); ok {
//		project, _, err = c.ProjectService.Get(domain, name.(string))
//		if err != nil {
//			return diag.FromErr(err)
//		}
//	} else {
//		displayName := d.Get("display_name").(string)
//		projects, _, err := c.ProjectService.GetListAll(domain, nil)
//		if err != nil {
//			return diag.FromErr(err)
//		}
//		for _, p := range projects.Projects {
//			if p.DisplayName == displayName {
//				project = p
//				break
//			}
//		}
//		if project == nil {
//			return diag.Errorf("Project not found")
//		}
//	}
//	err = util.ApiProjectToResourceData(domain, project, d, true)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}
