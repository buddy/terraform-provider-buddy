package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
	"time"
)

type testAccPipelineExpectedAttributes struct {
	Name                      string
	On                        string
	AlwaysFromScratch         bool
	FailOnPrepareEnvWarning   bool
	FetchAllRefs              bool
	AutoClearCache            bool
	NoSkipToMostRecent        bool
	DoNotCreateCommitStatus   bool
	StartDate                 string
	Delay                     int
	Cron                      string
	Paused                    bool
	Priority                  string
	IgnoreFailOnProjectStatus bool
	ExecutionMessageTemplate  string
	TargetSiteUrl             string
	Creator                   *api.Profile
	Project                   *api.Project
	Ref                       string
	Event                     *api.PipelineEvent
	TriggerConditions         []*api.PipelineTriggerCondition
}

func TestAccPipeline_schedule(t *testing.T) {
	var pipeline api.Pipeline
	var project api.Project
	var profile api.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	startDate := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	newStartDate := time.Now().UTC().Add(5 * time.Hour).Format(time.RFC3339)
	priority := api.PipelinePriorityLow
	newPriority := api.PipelinePriorityHigh
	delay := 5
	newDelay := 7
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigSchedule(domain, projectName, name, startDate, delay, true, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						On:                      api.PipelineOnSchedule,
						Project:                 &project,
						Creator:                 &profile,
						StartDate:               startDate,
						Delay:                   delay,
						Priority:                priority,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  true,
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigSchedule(domain, projectName, newName, newStartDate, newDelay, false, newPriority),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						On:                      api.PipelineOnSchedule,
						Project:                 &project,
						Creator:                 &profile,
						StartDate:               newStartDate,
						Delay:                   newDelay,
						Priority:                newPriority,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  false,
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_date"},
			},
		},
	})
}

func TestAccPipeline_schedule_cron(t *testing.T) {
	var pipeline api.Pipeline
	var project api.Project
	var profile api.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	cron := "15 14 1 * *"
	newCron := "0 22 * * 1-5"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigScheduleCron(domain, projectName, name, cron, true),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						On:                      api.PipelineOnSchedule,
						Project:                 &project,
						Creator:                 &profile,
						Cron:                    cron,
						Priority:                api.PipelinePriorityNormal,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  true,
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigScheduleCron(domain, projectName, newName, newCron, false),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						On:                      api.PipelineOnSchedule,
						Project:                 &project,
						Creator:                 &profile,
						Cron:                    newCron,
						Priority:                api.PipelinePriorityNormal,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  false,
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_date"},
			},
		},
	})
}

func testAccPipelineConfigScheduleCron(domain string, projectName string, name string, cron string, paused bool) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "bar" {
	domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
    on = "SCHEDULE"
	cron = "%s"
	paused = %t
}
`, domain, projectName, name, cron, paused)
}

func testAccPipelineConfigSchedule(domain string, projectName string, name string, startDate string, delay int, paused bool, priority string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "bar" {
	domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
    on = "SCHEDULE"
	start_date = "%s"
	delay = %d
	paused = %t
	priority = "%s"
}
`, domain, projectName, name, startDate, delay, paused, priority)
}

func TestAccPipeline_event(t *testing.T) {
	var pipeline api.Pipeline
	var project api.Project
	var profile api.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	ref := util.RandString(10)
	newRef := util.RandString(10)
	tcChangePath := "/path"
	newTcChangePath := "/abc"
	tcVarKey := util.RandString(10)
	newTcVarKey := util.RandString(10)
	tcVarValue := util.RandString(10)
	newTcVarValue := util.RandString(10)
	tcHours := 5
	newTcHours := 10
	tcDays := 1
	newTcDays := 3
	tcZoneId := "America/Monterrey"
	newTcZoneId := "Europe/Amsterdam"
	eventType := api.PipelineEventTypePush
	newEventType := api.PipelineEventTypeCreateRef
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigEvent(domain, projectName, name, eventType, ref, tcChangePath, tcVarKey, tcVarValue, tcHours, tcDays, tcZoneId),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						On:                      api.PipelineOnEvent,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Priority:                api.PipelinePriorityNormal,
						Event: &api.PipelineEvent{
							Type: eventType,
							Refs: []string{ref},
						},
						TriggerConditions: []*api.PipelineTriggerCondition{
							{
								TriggerCondition: api.PipelineTriggerConditionOnChange,
							},
							{
								TriggerCondition:      api.PipelineTriggerConditionOnChangeAtPath,
								TriggerConditionPaths: []string{tcChangePath},
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarIs,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarIsNot,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarContains,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarNotContains,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition: api.PipelineTriggerConditionDateTime,
								TriggerHours:     []int{tcHours},
								TriggerDays:      []int{tcDays},
								ZoneId:           tcZoneId,
							},
						},
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigEvent(domain, projectName, newName, newEventType, newRef, newTcChangePath, newTcVarKey, newTcVarValue, newTcHours, newTcDays, newTcZoneId),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						On:                      api.PipelineOnEvent,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Priority:                api.PipelinePriorityNormal,
						Event: &api.PipelineEvent{
							Type: newEventType,
							Refs: []string{newRef},
						},
						TriggerConditions: []*api.PipelineTriggerCondition{
							{
								TriggerCondition: api.PipelineTriggerConditionOnChange,
							},
							{
								TriggerCondition:      api.PipelineTriggerConditionOnChangeAtPath,
								TriggerConditionPaths: []string{newTcChangePath},
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarIs,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarIsNot,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarContains,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition:     api.PipelineTriggerConditionVarNotContains,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition: api.PipelineTriggerConditionDateTime,
								TriggerHours:     []int{newTcHours},
								TriggerDays:      []int{newTcDays},
								ZoneId:           newTcZoneId,
							},
						},
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_date"},
			},
		},
	})
}

func TestAccPipeline_click(t *testing.T) {
	var pipeline api.Pipeline
	var project api.Project
	var profile api.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	ref := util.RandString(10)
	newRef := util.RandString(10)
	msgTemplate := util.RandString(10)
	newMsgTemplate := util.RandString(10)
	targetUrl := "https://127.0.0.1"
	newTargetUrl := "https://192.168.1.1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigClick(domain, projectName, name, true, false, true, false, true, false, true, msgTemplate, targetUrl, ref),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						On:                        api.PipelineOnClick,
						AlwaysFromScratch:         true,
						AutoClearCache:            false,
						NoSkipToMostRecent:        true,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: true,
						FetchAllRefs:              true,
						FailOnPrepareEnvWarning:   false,
						Priority:                  api.PipelinePriorityNormal,
						TargetSiteUrl:             targetUrl,
						ExecutionMessageTemplate:  msgTemplate,
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigClick(domain, projectName, newName, false, true, false, true, false, true, false, newMsgTemplate, newTargetUrl, newRef),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      newName,
						On:                        api.PipelineOnClick,
						AlwaysFromScratch:         false,
						AutoClearCache:            true,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   true,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   true,
						Priority:                  api.PipelinePriorityNormal,
						TargetSiteUrl:             newTargetUrl,
						ExecutionMessageTemplate:  newMsgTemplate,
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       newRef,
					}),
				),
			},
			// import pipeline
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_date"},
			},
		},
	})
}

func testAccPipelineAttributes(n string, pipeline *api.Pipeline, want *testAccPipelineExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPipelineId, _ := strconv.Atoi(attrs["pipeline_id"])
		attrsAlwaysFromScratch, _ := strconv.ParseBool(attrs["always_from_scratch"])
		attrsFailOnPrepareEnvWarning, _ := strconv.ParseBool(attrs["fail_on_prepare_env_warning"])
		attrsFetchAllRefs, _ := strconv.ParseBool(attrs["fetch_all_refs"])
		attrsAutoClearCache, _ := strconv.ParseBool(attrs["auto_clear_cache"])
		attrsNoSkipToMostRecent, _ := strconv.ParseBool(attrs["no_skip_to_most_recent"])
		attrsDoNotCreateCommitStatus, _ := strconv.ParseBool(attrs["do_not_create_commit_status"])
		attrsIgnoreFailOnProjectStatus, _ := strconv.ParseBool(attrs["ignore_fail_on_project_status"])
		attrsDelay, _ := strconv.Atoi(attrs["delay"])
		attrsPaused, _ := strconv.ParseBool(attrs["paused"])
		attrsCreatorMemberId, _ := strconv.Atoi(attrs["creator.0.member_id"])
		if err := util.CheckFieldEqualAndSet("Name", pipeline.Name, want.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], want.Name); err != nil {
			return err
		}
		if err := util.CheckFieldSet("HtmlUrl", pipeline.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], pipeline.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Priority", pipeline.Priority, want.Priority); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("priority", attrs["priority"], want.Priority); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("Id", pipeline.Id); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("pipeline_id", attrsPipelineId, pipeline.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("On", pipeline.On, want.On); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("on", attrs["on"], want.On); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AlwaysFromScratch", pipeline.AlwaysFromScratch, want.AlwaysFromScratch); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("always_from_scratch", attrsAlwaysFromScratch, want.AlwaysFromScratch); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("fail_on_prepare_env_warning", attrsFailOnPrepareEnvWarning, want.FailOnPrepareEnvWarning); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("FailOnPrepareEnvWarning", pipeline.FailOnPrepareEnvWarning, want.FailOnPrepareEnvWarning); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("fetch_all_refs", attrsFetchAllRefs, want.FetchAllRefs); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("FetchAllRefs", pipeline.FetchAllRefs, want.FetchAllRefs); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AutoClearCache", pipeline.AutoClearCache, want.AutoClearCache); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("auto_clear_cache", attrsAutoClearCache, want.AutoClearCache); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("NoSkipToMostRecent", pipeline.NoSkipToMostRecent, want.NoSkipToMostRecent); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("no_skip_to_most_recent", attrsNoSkipToMostRecent, want.NoSkipToMostRecent); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("DoNotCreateCommitStatus", pipeline.DoNotCreateCommitStatus, want.DoNotCreateCommitStatus); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("do_not_create_commit_status", attrsDoNotCreateCommitStatus, want.DoNotCreateCommitStatus); err != nil {
			return err
		}
		if want.StartDate != "" {
			if err := util.CheckFieldEqualAndSet("start_date", attrs["start_date"], want.StartDate); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("StartDate", pipeline.StartDate, want.StartDate); err != nil {
				return err
			}
		}
		if want.Delay != 0 {
			if err := util.CheckIntFieldEqualAndSet("delay", attrsDelay, want.Delay); err != nil {
				return err
			}
			if err := util.CheckIntFieldEqualAndSet("Delay", pipeline.Delay, want.Delay); err != nil {
				return err
			}
		}
		if want.Cron != "" {
			if err := util.CheckFieldEqualAndSet("cron", attrs["cron"], want.Cron); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Cron", pipeline.Cron, want.Cron); err != nil {
				return err
			}
		}
		if err := util.CheckBoolFieldEqual("paused", attrsPaused, want.Paused); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Paused", pipeline.Paused, want.Paused); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("ignore_fail_on_project_status", attrsIgnoreFailOnProjectStatus, want.IgnoreFailOnProjectStatus); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("IgnoreFailOnProjectStatus", pipeline.IgnoreFailOnProjectStatus, want.IgnoreFailOnProjectStatus); err != nil {
			return err
		}
		if want.ExecutionMessageTemplate != "" {
			if err := util.CheckFieldEqualAndSet("execution_message_template", attrs["execution_message_template"], want.ExecutionMessageTemplate); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("ExecutionMessageTemplate", pipeline.ExecutionMessageTemplate, want.ExecutionMessageTemplate); err != nil {
				return err
			}
		}
		if want.TargetSiteUrl != "" {
			if err := util.CheckFieldEqualAndSet("target_site_url", attrs["target_site_url"], want.TargetSiteUrl); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("TargetSiteUrl", pipeline.TargetSiteUrl, want.TargetSiteUrl); err != nil {
				return err
			}
		}
		if err := util.CheckFieldSet("create_date", attrs["create_date"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("CreateDate", pipeline.CreateDate); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Creator.Name", pipeline.Creator.Name, want.Creator.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("creator.0.name", attrs["creator.0.name"], want.Creator.Name); err != nil {
			return err
		}
		if err := util.CheckFieldSet("Creator.HtmlUrl", pipeline.Creator.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldSet("creator.0.html_url", attrs["creator.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Creator.AvatarUrl", pipeline.Creator.AvatarUrl, want.Creator.AvatarUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("creator.0.avatar_url", attrs["creator.0.avatar_url"], want.Creator.AvatarUrl); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("Creator.Id", pipeline.Creator.Id, want.Creator.Id); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("creator.0.member_id", attrsCreatorMemberId, want.Creator.Id); err != nil {
			return err
		}
		if err := util.CheckFieldSet("Creator.Email", pipeline.Creator.Email); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Project.HtmlUrl", pipeline.Project.HtmlUrl, want.Project.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project.0.html_url", attrs["project.0.html_url"], want.Project.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Project.Name", pipeline.Project.Name, want.Project.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project.0.name", attrs["project.0.name"], want.Project.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Project.DisplayName", pipeline.Project.DisplayName, want.Project.DisplayName); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project.0.display_name", attrs["project.0.display_name"], want.Project.DisplayName); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Project.Status", pipeline.Project.Status, want.Project.Status); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project.0.status", attrs["project.0.status"], want.Project.Status); err != nil {
			return err
		}
		if want.Ref != "" {
			if err := util.CheckFieldEqualAndSet("Refs[0]", pipeline.Refs[0], want.Ref); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("refs.0", attrs["refs.0"], want.Ref); err != nil {
				return err
			}
		}
		if want.Event != nil {
			if err := util.CheckFieldEqualAndSet("Events[0].Type", pipeline.Events[0].Type, want.Event.Type); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("event.0.type", attrs["event.0.type"], want.Event.Type); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Events[0].Refs[0]", pipeline.Events[0].Refs[0], want.Event.Refs[0]); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("event.0.refs.0", attrs["event.0.refs.0"], want.Event.Refs[0]); err != nil {
				return err
			}
		}
		if len(want.TriggerConditions) > 0 {
			for i, triggerCondition := range want.TriggerConditions {
				if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerCondition", i), pipeline.TriggerConditions[i].TriggerCondition, triggerCondition.TriggerCondition); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.condition", i), attrs[fmt.Sprintf("trigger_condition.%d.condition", i)], triggerCondition.TriggerCondition); err != nil {
					return err
				}
				if triggerCondition.TriggerCondition == api.PipelineTriggerConditionOnChangeAtPath {
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerConditionPaths[0]", i), pipeline.TriggerConditions[i].TriggerConditionPaths[0], triggerCondition.TriggerConditionPaths[0]); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.paths.0", i), attrs[fmt.Sprintf("trigger_condition.%d.paths.0", i)], triggerCondition.TriggerConditionPaths[0]); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == api.PipelineTriggerConditionVarIs ||
					triggerCondition.TriggerCondition == api.PipelineTriggerConditionVarIsNot ||
					triggerCondition.TriggerCondition == api.PipelineTriggerConditionVarContains ||
					triggerCondition.TriggerCondition == api.PipelineTriggerConditionVarNotContains {
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerVariableKey", i), pipeline.TriggerConditions[i].TriggerVariableKey, triggerCondition.TriggerVariableKey); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.variable_key", i), attrs[fmt.Sprintf("trigger_condition.%d.variable_key", i)], triggerCondition.TriggerVariableKey); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerVariableValue", i), pipeline.TriggerConditions[i].TriggerVariableValue, triggerCondition.TriggerVariableValue); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.variable_value", i), attrs[fmt.Sprintf("trigger_condition.%d.variable_value", i)], triggerCondition.TriggerVariableValue); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == api.PipelineTriggerConditionDateTime {
					attrsTriggerConditionsHours, _ := strconv.Atoi(attrs[fmt.Sprintf("trigger_condition.%d.hours.0", i)])
					attrsTriggerConditionsDays, _ := strconv.Atoi(attrs[fmt.Sprintf("trigger_condition.%d.days.0", i)])
					if err := util.CheckIntFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerHours[0]", i), pipeline.TriggerConditions[i].TriggerHours[0], triggerCondition.TriggerHours[0]); err != nil {
						return err
					}
					if err := util.CheckIntFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.hours[0]", i), attrsTriggerConditionsHours, triggerCondition.TriggerHours[0]); err != nil {
						return err
					}
					if err := util.CheckIntFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerDays[0]", i), pipeline.TriggerConditions[i].TriggerDays[0], triggerCondition.TriggerDays[0]); err != nil {
						return err
					}
					if err := util.CheckIntFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.days[0]", i), attrsTriggerConditionsDays, triggerCondition.TriggerDays[0]); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].ZoneId", i), pipeline.TriggerConditions[i].ZoneId, triggerCondition.ZoneId); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.zone_id", i), attrs[fmt.Sprintf("trigger_condition.%d.zone_id", i)], triggerCondition.ZoneId); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == api.PipelineTriggerConditionSuccessPipeline {
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerProjectName", i), pipeline.TriggerConditions[i].TriggerProjectName, triggerCondition.TriggerProjectName); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.project_name", i), attrs[fmt.Sprintf("trigger_condition.%d.project_name", i)], triggerCondition.TriggerProjectName); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerPipelineName", i), pipeline.TriggerConditions[i].TriggerPipelineName, triggerCondition.TriggerPipelineName); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.pipeline_name", i), attrs[fmt.Sprintf("trigger_condition.%d.pipeline_name", i)], triggerCondition.TriggerPipelineName); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

func testAccPipelineGet(n string, pipeline *api.Pipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, projectName, pid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		pipelineId, err := strconv.Atoi(pid)
		if err != nil {
			return err
		}
		p, _, err := acc.ApiClient.PipelineService.Get(domain, projectName, pipelineId)
		if err != nil {
			return err
		}
		*pipeline = *p
		return nil
	}
}

func testAccPipelineConfigEvent(domain string, projectName string, name string, eventType string, ref string, tcChangePath string, tcVarKey string, tcVarValue string, tcHours int, tcDays int, tcZoneId string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
    on = "EVENT" 
    event {
        type = "%s"
        refs = ["%s"]
    }
    trigger_condition {
        condition = "ON_CHANGE"
    }
    trigger_condition {
        condition = "ON_CHANGE_AT_PATH"
        paths = ["%s"]
    }
    trigger_condition {
        condition = "VAR_IS"
        variable_key = "%s"
        variable_value = "%s"
    }
    trigger_condition {
        condition = "VAR_IS_NOT"
        variable_key = "%s"
        variable_value = "%s"
    }
    trigger_condition {
        condition = "VAR_CONTAINS"
        variable_key = "%s"
        variable_value = "%s"
    }
    trigger_condition {
        condition = "VAR_NOT_CONTAINS"
        variable_key = "%s"
        variable_value = "%s"
    }
    trigger_condition {
        condition = "DATETIME"
        hours = [%d]
        days = [%d]
        zone_id = "%s"
    }
}
`, domain, projectName, name, eventType, ref, tcChangePath, tcVarKey, tcVarValue, tcVarKey, tcVarValue, tcVarKey, tcVarValue, tcVarKey, tcVarValue, tcHours, tcDays, tcZoneId)
}

func testAccPipelineConfigClick(domain string, projectName string, name string, alwaysFromScratch bool, failOnPrepareEnvWarning bool, fetchAllRefs bool, autoClearCache bool, noSkipToMostRecent bool, doNotCreateCommitStatus bool, ignoreFailOnProjectStatus bool, executionMessageTemplate string, targetSiteUrl string, ref string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
    on = "CLICK"
    always_from_scratch = %t
	fail_on_prepare_env_warning = %t
	fetch_all_refs = %t
    auto_clear_cache = %t
    no_skip_to_most_recent = %t
    do_not_create_commit_status = %t
    ignore_fail_on_project_status = %t
    execution_message_template = "%s"
    target_site_url = "%s"
    refs = ["%s"]
}
`, domain, projectName, name, alwaysFromScratch, failOnPrepareEnvWarning, fetchAllRefs, autoClearCache, noSkipToMostRecent, doNotCreateCommitStatus, ignoreFailOnProjectStatus, executionMessageTemplate, targetSiteUrl, ref)
}

func testAccPipelineCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_pipeline" {
			continue
		}
		domain, projectName, pid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		pipelineId, err := strconv.Atoi(pid)
		if err != nil {
			return err
		}
		pipeline, resp, err := acc.ApiClient.PipelineService.Get(domain, projectName, pipelineId)
		if err == nil && pipeline != nil {
			return util.ErrorResourceExists()
		}
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}
