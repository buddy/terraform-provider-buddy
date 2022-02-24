package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/http"
	"strconv"
	"strings"
)

func Pipeline() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a pipeline\n\n" +
			"Token scopes required: `WORKSPACE`, `EXECUTION_MANAGE`, `EXECUTION_INFO`",
		CreateContext: createContextPipeline,
		ReadContext:   readContextPipeline,
		UpdateContext: updateContextPipeline,
		DeleteContext: deleteContextPipeline,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "The workspace's URL handle",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"project_name": {
				Description: "The project's name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"html_url": {
				Description: "The pipeline's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pipeline_id": {
				Description: "The pipeline's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The pipeline's name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"definition_source": {
				Description: "The pipeline's definition source. Allowed: `LOCAL`, `REMOTE`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					api.PipelineDefinitionSourceLocal,
					api.PipelineDefinitionSourceRemote,
				}, false),
			},
			"remote_project_name": {
				Description: "The pipeline's remote definition project name. Set it if `definition_source: REMOTE`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"remote_branch": {
				Description: "The pipeline's remote definition branch name. Set it if `definition_source: REMOTE`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"remote_path": {
				Description: "The pipeline's remote definition path. Set it if `definition_source: REMOTE`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"remote_parameter": {
				Description: "The pipeline's remote definition parameters. Set it if `definition_source: REMOTE`",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"on": {
				Description: "The pipeline's trigger mode. Required when not using remote definition. Allowed: `CLICK`, `EVENT`, `SCHEDULE`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					api.PipelineOnClick,
					api.PipelineOnEvent,
					api.PipelineOnSchedule,
				}, false),
			},
			"priority": {
				Description: "The pipeline's priority. Allowed: `LOW`, `NORMAL`, `HIGH`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					api.PipelinePriorityHigh,
					api.PipelinePriorityNormal,
					api.PipelinePriorityLow,
				}, false),
			},
			"fetch_all_refs": {
				Description: "Defines either or not fetch all refs from repository",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"always_from_scratch": {
				Description: "Defines whether or not to upload everything from scratch on every run",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"fail_on_prepare_env_warning": {
				Description: "Defines either or not run should fail if any warning occurs in prepare environment",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"auto_clear_cache": {
				Description: "Defines whether or not to automatically clear cache before running the pipeline",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"no_skip_to_most_recent": {
				Description: "Defines whether or not to skip run to the most recent run",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"do_not_create_commit_status": {
				Description: "Defines whether or not to omit sending commit statuses to GitHub or GitLab upon execution",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"start_date": {
				Description: "The pipeline's start date. Required if the pipeline is set to `on: SCHEDULE` and no `cron` is specified. Format: `2016-11-18T12:38:16.000Z`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ConflictsWith: []string{
					"cron",
				},
				RequiredWith: []string{
					"delay",
				},
			},
			"delay": {
				Description: "The pipeline's runs interval (in minutes). Required if the pipeline is set to `on: SCHEDULE` and no `cron` is specified",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ConflictsWith: []string{
					"cron",
				},
				RequiredWith: []string{
					"start_date",
				},
			},
			"clone_depth": {
				Description: "The pipeline's filesystem clone depth. Creates a shallow clone with a history truncated to the specified number of commits",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"cron": {
				Description: "The pipeline's CRON expression. Required if the pipeline is set to `on: SCHEDULE` and neither `start_date` nor `delay` is specified",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ConflictsWith: []string{
					"delay",
					"start_date",
				},
			},
			"paused": {
				Description: "Is the pipeline's run paused. Restricted to `on: SCHEDULE`",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"ignore_fail_on_project_status": {
				Description: "If set to true the status of a given pipeline will be ignored on the projects' dashboard",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"execution_message_template": {
				Description: "The pipeline's run title. Default: `$BUDDY_EXECUTION_REVISION_SUBJECT`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"worker": {
				Description: "The pipeline's worker name. Only for `Buddy Enterprise`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"target_site_url": {
				Description: "The pipeline's website target URL",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"last_execution_status": {
				Description: "The pipeline's last run status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_execution_revision": {
				Description: "The pipeline's last run revision",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"create_date": {
				Description: "The pipeline's date of creation",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"creator": {
				Description: "The pipeline's creator",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"avatar_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"workspace_owner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"project": {
				Description: "The pipeline's project",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"refs": {
				Description: "The pipeline's list of refs. Set it if `on: CLICK`",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{
					"event",
				},
			},
			"tags": {
				Description: "The pipeline's list of tags. Only for `Buddy Enterprise`",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldVal, newVal := d.GetChange("tags")
					oldTags := *util.InterfaceStringSetToPointer(oldVal)
					newTags := *util.InterfaceStringSetToPointer(newVal)
					if len(oldTags) != len(newTags) {
						return false
					}
					for _, oldTag := range oldTags {
						found := false
						for _, newTag := range newTags {
							if strings.EqualFold(oldTag, newTag) {
								found = true
								break
							}
						}
						if !found {
							return false
						}
					}
					return true
				},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"event": {
				Description: "The pipeline's list of events. Set it if `on: EVENT`",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								api.PipelineEventTypePush,
								api.PipelineEventTypeCreateRef,
								api.PipelineEventTypeDeleteRef,
							}, false),
						},
						"refs": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				ConflictsWith: []string{
					"refs",
				},
			},
			"trigger_condition": {
				Description: "The pipeline's list of trigger conditions",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								api.PipelineTriggerConditionOnChange,
								api.PipelineTriggerConditionOnChangeAtPath,
								api.PipelineTriggerConditionVarIs,
								api.PipelineTriggerConditionVarIsNot,
								api.PipelineTriggerConditionVarContains,
								api.PipelineTriggerConditionVarNotContains,
								api.PipelineTriggerConditionDateTime,
								api.PipelineTriggerConditionSuccessPipeline,
							}, false),
						},
						"paths": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"variable_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"variable_value": {
							Type:      schema.TypeString,
							Optional:  true,
							Computed:  true,
							Sensitive: true,
						},
						"hours": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"days": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"pipeline_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func deleteContextPipeline(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, projectName, pid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	pipelineId, err := strconv.Atoi(pid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.PipelineService.Delete(domain, projectName, pipelineId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextPipeline(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain, projectName, pid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	pipelineId, err := strconv.Atoi(pid)
	if err != nil {
		return diag.FromErr(err)
	}
	opt := api.PipelineOperationOptions{
		Name: util.InterfaceStringToPointer(d.Get("name")),
	}
	if d.HasChange("on") {
		opt.On = util.InterfaceStringToPointer(d.Get("on"))
	}
	if d.HasChange("always_from_scratch") {
		opt.AlwaysFromScratch = util.InterfaceBoolToPointer(d.Get("always_from_scratch"))
	}
	if d.HasChange("priority") {
		opt.Priority = util.InterfaceStringToPointer(d.Get("priority"))
	}
	if d.HasChange("fail_on_prepare_env_warning") {
		opt.FailOnPrepareEnvWarning = util.InterfaceBoolToPointer(d.Get("fail_on_prepare_env_warning"))
	}
	if d.HasChange("fetch_all_refs") {
		opt.FetchAllRefs = util.InterfaceBoolToPointer(d.Get("fetch_all_refs"))
	}
	if d.HasChange("auto_clear_cache") {
		opt.AutoClearCache = util.InterfaceBoolToPointer(d.Get("auto_clear_cache"))
	}
	if d.HasChange("no_skip_to_most_recent") {
		opt.NoSkipToMostRecent = util.InterfaceBoolToPointer(d.Get("no_skip_to_most_recent"))
	}
	if d.HasChange("do_not_create_commit_status") {
		opt.DoNotCreateCommitStatus = util.InterfaceBoolToPointer(d.Get("do_not_create_commit_status"))
	}
	if d.HasChange("start_date") {
		opt.StartDate = util.InterfaceStringToPointer(d.Get("start_date"))
	}
	if d.HasChange("delay") {
		opt.Delay = util.InterfaceIntToPointer(d.Get("delay"))
	}
	if d.HasChange("clone_depth") {
		opt.CloneDepth = util.InterfaceIntToPointer(d.Get("clone_depth"))
	}
	if d.HasChange("cron") {
		opt.Cron = util.InterfaceStringToPointer(d.Get("cron"))
	}
	if d.HasChange("paused") {
		opt.Paused = util.InterfaceBoolToPointer(d.Get("paused"))
	}
	if d.HasChange("ignore_fail_on_project_status") {
		opt.IgnoreFailOnProjectStatus = util.InterfaceBoolToPointer(d.Get("ignore_fail_on_project_status"))
	}
	if d.HasChange("execution_message_template") {
		opt.ExecutionMessageTemplate = util.InterfaceStringToPointer(d.Get("execution_message_template"))
	}
	if d.HasChange("worker") {
		opt.Worker = util.InterfaceStringToPointer(d.Get("worker"))
	}
	if d.HasChange("target_site_url") {
		opt.TargetSiteUrl = util.InterfaceStringToPointer(d.Get("target_site_url"))
	}
	if d.HasChange("refs") {
		opt.Refs = util.InterfaceStringSetToPointer(d.Get("refs"))
	}
	if d.HasChange("tags") {
		opt.Tags = util.InterfaceStringSetToPointer(d.Get("tags"))
	}
	if d.HasChange("event") {
		opt.Events = util.MapPipelineEventsToApi(d.Get("event"))
	}
	if d.HasChange("trigger_condition") {
		opt.TriggerConditions = util.MapTriggerConditionsToApi(d.Get("trigger_condition"))
	}
	if d.HasChange("definition_source") {
		opt.DefinitionSource = util.InterfaceStringToPointer(d.Get("definition_source"))
	}
	if d.HasChange("remote_path") {
		opt.RemotePath = util.InterfaceStringToPointer(d.Get("remote_path"))
	}
	if d.HasChange("remote_branch") {
		opt.RemoteBranch = util.InterfaceStringToPointer(d.Get("remote_branch"))
	}
	if d.HasChange("remote_project_name") {
		opt.RemoteProjectName = util.InterfaceStringToPointer(d.Get("remote_project_name"))
	}
	if d.HasChange("remote_parameter") {
		opt.RemoteParameters = util.MapPipelineRemoteParametersToApi(d.Get("remote_parameter"))
	}
	_, _, err = c.PipelineService.Update(domain, projectName, pipelineId, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextPipeline(ctx, d, meta)
}

func readContextPipeline(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, projectName, pid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	pipelineId, err := strconv.Atoi(pid)
	if err != nil {
		return diag.FromErr(err)
	}
	pipeline, resp, err := c.PipelineService.Get(domain, projectName, pipelineId)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiPipelineToResourceData(domain, projectName, pipeline, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextPipeline(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	opt := api.PipelineOperationOptions{
		Name:                    util.InterfaceStringToPointer(d.Get("name")),
		FailOnPrepareEnvWarning: util.InterfaceBoolToPointer(d.Get("fail_on_prepare_env_warning")),
		FetchAllRefs:            util.InterfaceBoolToPointer(d.Get("fetch_all_refs")),
	}
	if on, ok := d.GetOk("on"); ok {
		opt.On = util.InterfaceStringToPointer(on)
	}
	if alwaysFromScratch, ok := d.GetOk("always_from_scratch"); ok {
		opt.AlwaysFromScratch = util.InterfaceBoolToPointer(alwaysFromScratch)
	}
	if priority, ok := d.GetOk("priority"); ok {
		opt.Priority = util.InterfaceStringToPointer(priority)
	}
	if autoClearCache, ok := d.GetOk("auto_clear_cache"); ok {
		opt.AutoClearCache = util.InterfaceBoolToPointer(autoClearCache)
	}
	if noSkipToMostRecent, ok := d.GetOk("no_skip_to_most_recent"); ok {
		opt.NoSkipToMostRecent = util.InterfaceBoolToPointer(noSkipToMostRecent)
	}
	if doNotCreateCommitStatus, ok := d.GetOk("do_not_create_commit_status"); ok {
		opt.DoNotCreateCommitStatus = util.InterfaceBoolToPointer(doNotCreateCommitStatus)
	}
	if startDate, ok := d.GetOk("start_date"); ok {
		opt.StartDate = util.InterfaceStringToPointer(startDate)
	}
	if delay, ok := d.GetOk("delay"); ok {
		opt.Delay = util.InterfaceIntToPointer(delay)
	}
	if cloneDepth, ok := d.GetOk("clone_depth"); ok {
		opt.CloneDepth = util.InterfaceIntToPointer(cloneDepth)
	}
	if cron, ok := d.GetOk("cron"); ok {
		opt.Cron = util.InterfaceStringToPointer(cron)
	}
	if paused, ok := d.GetOk("paused"); ok {
		opt.Paused = util.InterfaceBoolToPointer(paused)
	}
	if ignoreFailOnProjectStatus, ok := d.GetOk("ignore_fail_on_project_status"); ok {
		opt.IgnoreFailOnProjectStatus = util.InterfaceBoolToPointer(ignoreFailOnProjectStatus)
	}
	if executionMessageTemplate, ok := d.GetOk("execution_message_template"); ok {
		opt.ExecutionMessageTemplate = util.InterfaceStringToPointer(executionMessageTemplate)
	}
	if worker, ok := d.GetOk("worker"); ok {
		opt.Worker = util.InterfaceStringToPointer(worker)
	}
	if targetSiteUrl, ok := d.GetOk("target_site_url"); ok {
		opt.TargetSiteUrl = util.InterfaceStringToPointer(targetSiteUrl)
	}
	if refs, ok := d.GetOk("refs"); ok {
		opt.Refs = util.InterfaceStringSetToPointer(refs)
	}
	if tags, ok := d.GetOk("tags"); ok {
		opt.Tags = util.InterfaceStringSetToPointer(tags)
	}
	if events, ok := d.GetOk("event"); ok {
		opt.Events = util.MapPipelineEventsToApi(events)
	}
	if definitionSource, ok := d.GetOk("definition_source"); ok {
		opt.DefinitionSource = util.InterfaceStringToPointer(definitionSource)
	}
	if remotePath, ok := d.GetOk("remote_path"); ok {
		opt.RemotePath = util.InterfaceStringToPointer(remotePath)
	}
	if remoteBranch, ok := d.GetOk("remote_branch"); ok {
		opt.RemoteBranch = util.InterfaceStringToPointer(remoteBranch)
	}
	if remoteProjectName, ok := d.GetOk("remote_project_name"); ok {
		opt.RemoteProjectName = util.InterfaceStringToPointer(remoteProjectName)
	}
	if remoteParameter, ok := d.GetOk("remote_parameter"); ok {
		opt.RemoteParameters = util.MapPipelineRemoteParametersToApi(remoteParameter)
	}
	if triggerCondition, ok := d.GetOk("trigger_condition"); ok {
		opt.TriggerConditions = util.MapTriggerConditionsToApi(triggerCondition)
	}
	pipeline, _, err := c.PipelineService.Create(domain, projectName, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, projectName, strconv.Itoa(pipeline.Id)))
	return readContextPipeline(ctx, d, meta)
}
