package source

// todo project groups
//import (
//	"buddy-terraform/buddy/util"
//	"context"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
//	"regexp"
//)
//
//func ProjectGroups() *schema.Resource {
//	return &schema.Resource{
//		Description: "List project groups and optionally filter them by name\n\n" +
//			"Token scope required: `WORKSPACE`",
//		ReadContext: readContextProjectGroups,
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
//			"project_name": {
//				Description: "The project's name",
//				Type:        schema.TypeString,
//				Required:    true,
//			},
//			"name_regex": {
//				Description:  "The group's name regular expression to match",
//				Type:         schema.TypeString,
//				Optional:     true,
//				ValidateFunc: validation.StringIsValidRegExp,
//			},
//			"groups": {
//				Description: "List of groups",
//				Type:        schema.TypeList,
//				Computed:    true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"html_url": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"group_id": {
//							Type:     schema.TypeInt,
//							Computed: true,
//						},
//						"name": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func readContextProjectGroups(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	var nameRegex *regexp.Regexp
//	domain := d.Get("domain").(string)
//	projectName := d.Get("project_name").(string)
//	groups, _, err := c.ProjectGroupService.GetProjectGroups(domain, projectName)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	var result []interface{}
//	if name, ok := d.GetOk("name_regex"); ok {
//		nameRegex = regexp.MustCompile(name.(string))
//	}
//	for _, g := range groups.Groups {
//		if nameRegex != nil && !nameRegex.MatchString(g.Name) {
//			continue
//		}
//		result = append(result, util.ApiShortGroupToMap(g))
//	}
//	d.SetId(util.UniqueString())
//	err = d.Set("groups", result)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}
