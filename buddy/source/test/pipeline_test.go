package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourcePipeline(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	ref := util.RandString(10)
	reason := util.RandString(10)
	prBranch := util.RandString(10)
	prEvent := buddy.PipelinePullRequestEventClosed
	gitChangeSet := buddy.PipelineGitChangeSetBasePullRequest
	filesystemChangeSet := buddy.PipelineFilesystemChangeSetBaseContents
	newGitChangeSet := buddy.PipelineGitChangeSetBaseLatestRunMatchingRef
	newFilesystemChangeSet := buddy.PipelineFilesystemChangeSetBaseDateModified
	cpu := buddy.PipelineCpuArm
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			// click
			{
				Config: testAccSourcePipelineConfigClick(domain, projectName, name, ref, buddy.PipelinePriorityHigh, cpu, true, true, gitChangeSet, filesystemChangeSet),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnClick, cpu, ref, "", "", "", buddy.PipelinePriorityHigh, true, true, gitChangeSet, filesystemChangeSet, false, "", buddy.PipelineGitConfigRefFixed, &buddy.PipelineGitConfig{
						Project: projectName,
						Branch:  "main",
						Path:    "def.yml",
					}),
					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnClick, cpu, ref, "", "", "", buddy.PipelinePriorityHigh, true, true, gitChangeSet, filesystemChangeSet, false, "", buddy.PipelineGitConfigRefFixed, &buddy.PipelineGitConfig{
						Project: projectName,
						Branch:  "main",
						Path:    "def.yml",
					}),
				),
			},
			// click disabled
			{
				Config: testAccSourcePipelineConfigClickDisabled(domain, projectName, name, ref, buddy.PipelinePriorityHigh, reason),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnClick, cpu, ref, "", "", "", buddy.PipelinePriorityHigh, true, true, buddy.PipelineGitChangeSetBaseLatestRun, filesystemChangeSet, true, reason, buddy.PipelineGitConfigRefNone, nil),
					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnClick, cpu, ref, "", "", "", buddy.PipelinePriorityHigh, true, true, buddy.PipelineGitChangeSetBaseLatestRun, filesystemChangeSet, true, reason, buddy.PipelineGitConfigRefNone, nil),
				),
			},
			// event
			{
				Config: testAccSourcePipelineConfigEvent(domain, projectName, name, ref, buddy.PipelinePriorityLow),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnEvent, cpu, "", ref, "", "", buddy.PipelinePriorityLow, true, true, buddy.PipelineGitChangeSetBaseLatestRun, filesystemChangeSet, false, "", buddy.PipelineGitConfigRefDynamic, nil),
					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnEvent, cpu, "", ref, "", "", buddy.PipelinePriorityLow, true, true, buddy.PipelineGitChangeSetBaseLatestRun, filesystemChangeSet, false, "", buddy.PipelineGitConfigRefDynamic, nil),
				),
			},
			// pr event
			{
				Config: testAccSourcePipelineConfigPullRequest(domain, projectName, name, prBranch, prEvent, false, false, newGitChangeSet, newFilesystemChangeSet),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnEvent, cpu, "", "", prBranch, prEvent, buddy.PipelinePriorityLow, false, false, newGitChangeSet, newFilesystemChangeSet, false, "", buddy.PipelineGitConfigRefNone, nil),
					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnEvent, cpu, "", "", prBranch, prEvent, buddy.PipelinePriorityLow, false, false, newGitChangeSet, newFilesystemChangeSet, false, "", buddy.PipelineGitConfigRefNone, nil),
				),
			},
		},
	})
}

func testAccPipelineGitConfig(attrs map[string]string, gitConfigRef string, gitConfig *buddy.PipelineGitConfig) error {
	if gitConfigRef != "" {
		if err := util.CheckFieldEqualAndSet("git_config_ref", attrs["git_config_ref"], gitConfigRef); err != nil {
			return err
		}
		attrsGitConfigProject, attrsGitConfigProjectExists := attrs["git_config.project"]
		attrsGitConfigBranch, attrsGitConfigBranchExists := attrs["git_config.branch"]
		attrsGitConfigPath, attrsGitConfigPathExists := attrs["git_config.path"]
		if gitConfig == nil {
			if attrsGitConfigProjectExists {
				return util.ErrorFieldSet("git_config.project")
			}
			if attrsGitConfigBranchExists {
				return util.ErrorFieldSet("git_config.branch")
			}
			if attrsGitConfigPathExists {
				return util.ErrorFieldSet("git_config.path")
			}
		} else {
			if err := util.CheckFieldEqualAndSet("git_config.project", attrsGitConfigProject, gitConfig.Project); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("git_config.branch", attrsGitConfigBranch, gitConfig.Branch); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("git_config.path", attrsGitConfigPath, gitConfig.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

func testAccSourcePipelineAttributes(n string, name string, on string, cpu string, ref string, eventRef string, prBranch string, prEvent string, priority string, descriptionRequired bool, concurrentPipelineRuns bool, gitChangesetBase string, filesystemChangesetBase string, disabled bool, disabledReason string, gitConfigRef string, gitConfig *buddy.PipelineGitConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPipelineId, _ := strconv.Atoi(attrs["pipeline_id"])
		attrsDisabled, _ := strconv.ParseBool(attrs["disabled"])
		attrsDescRequired, _ := strconv.ParseBool(attrs["description_required"])
		attrsConcurrentRuns, _ := strconv.ParseBool(attrs["concurrent_pipeline_runs"])
		if err := util.CheckBoolFieldEqual("description_required", attrsDescRequired, descriptionRequired); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("concurrent_pipeline_runs", attrsConcurrentRuns, concurrentPipelineRuns); err != nil {
			return err
		}
		if gitChangesetBase != "" {
			if err := util.CheckFieldEqualAndSet("git_changeset_base", attrs["git_changeset_base"], gitChangesetBase); err != nil {
				return err
			}
		}
		if filesystemChangesetBase != "" {
			if err := util.CheckFieldEqualAndSet("filesystem_changeset_base", attrs["filesystem_changeset_base"], filesystemChangesetBase); err != nil {
				return err
			}
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("pipeline_id", attrsPipelineId); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("on", attrs["on"], on); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("cpu", attrs["cpu"], cpu); err != nil {
			return err
		}
		if err := util.CheckFieldSet("last_execution_status", attrs["last_execution_status"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("priority", attrs["priority"], priority); err != nil {
			return err
		}
		if ref != "" {
			if err := util.CheckFieldEqualAndSet("refs.0", attrs["refs.0"], ref); err != nil {
				return err
			}
		}
		if eventRef != "" {
			if err := util.CheckFieldEqualAndSet("event.0.refs.0", attrs["event.0.refs.0"], eventRef); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("event.0.type", attrs["event.0.type"], buddy.PipelineEventTypePush); err != nil {
				return err
			}
		}
		if prBranch != "" && prEvent != "" {
			if err := util.CheckFieldEqualAndSet("event.0.branches.0", attrs["event.0.branches.0"], prBranch); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("event.0.events.0", attrs["event.0.events.0"], prEvent); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("event.0.type", attrs["event.0.type"], buddy.PipelineEventTypePullRequest); err != nil {
				return err
			}
		}
		if disabled {
			if err := util.CheckBoolFieldEqual("disabled", attrsDisabled, disabled); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("disabling_reason", attrs["disabling_reason"], disabledReason); err != nil {
				return err
			}
		}
		if err := testAccPipelineGitConfig(attrs, gitConfigRef, gitConfig); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourcePipelineConfigPullRequest(domain string, projectName string, name string, branch string, event string, concurrentPipelineRuns bool, descriptionRequired bool, gitChangesetBase string, filesystemChangesetBase string) string {
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
       type = "PULL_REQUEST"
       branches = ["%s"]
			 events = ["%s"]
   }
   concurrent_pipeline_runs = "%t"
   description_required = "%t"
   git_changeset_base = "%s"
   filesystem_changeset_base = "%s"
}

data "buddy_pipeline" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "${buddy_pipeline.bar.name}"
}

data "buddy_pipeline" "id" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
}
`, domain, projectName, name, branch, event, concurrentPipelineRuns, descriptionRequired, gitChangesetBase, filesystemChangesetBase)
}

func testAccSourcePipelineConfigEvent(domain string, projectName string, name string, ref string, priority string) string {
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
   git_config_ref = "DYNAMIC"
   event {
       type = "PUSH"
       refs = ["%s"]
   }
	priority = "%s"
}

data "buddy_pipeline" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "${buddy_pipeline.bar.name}"
}

data "buddy_pipeline" "id" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
}
`, domain, projectName, name, ref, priority)
}

func testAccSourcePipelineConfigClick(domain string, projectName string, name string, ref string, priority string, cpu string, concurrentPipelineRuns bool, descriptionRequired bool, gitChangesetBase string, filesystemChangesetBase string) string {
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
   cpu = "%s"
   refs = ["%s"]
	 git_config_ref = "FIXED"
   git_config = {
     project = "%s"
     branch = "main"
     path = "def.yml"
   }
	 priority = "%s"
   concurrent_pipeline_runs = "%t"
   description_required = "%t"
   git_changeset_base = "%s"
   filesystem_changeset_base = "%s"
}

data "buddy_pipeline" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "${buddy_pipeline.bar.name}"
}

data "buddy_pipeline" "id" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
}
`, domain, projectName, name, cpu, ref, projectName, priority, concurrentPipelineRuns, descriptionRequired, gitChangesetBase, filesystemChangesetBase)
}

func testAccSourcePipelineConfigClickDisabled(domain string, projectName string, name string, ref string, priority string, reason string) string {
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
   refs = ["%s"]
   git_config_ref = "NONE"
	 priority = "%s"
	 disabled = true
	 disabling_reason = "%s"
}

data "buddy_pipeline" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "${buddy_pipeline.bar.name}"
}

data "buddy_pipeline" "id" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
}
`, domain, projectName, name, ref, priority, reason)
}
