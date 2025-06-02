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

func TestAccSourcePipelines(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name1 := "aaaa" + util.RandString(10)
	name2 := util.RandString(10)
	ref := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourcePipelinesConfig(domain, projectName, name1, name2, ref),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePipelinesAttributes("data.buddy_pipelines.all", 2, "", ""),
					testAccSourcePipelinesAttributes("data.buddy_pipelines.name", 1, name1, ref),
				),
			},
		},
	})
}

func testAccSourcePipelinesAttributes(n string, count int, name string, ref string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPipelinesCount, _ := strconv.Atoi(attrs["pipelines.#"])
		attrsPipelineId, _ := strconv.Atoi(attrs["pipelines.0.pipeline_id"])
		attrsDescRequired, _ := strconv.ParseBool(attrs["pipelines.0.description_required"])
		attrsConcurrentRun, _ := strconv.ParseBool(attrs["pipelines.0.concurrent_pipeline_runs"])
		if err := util.CheckIntFieldEqual("pipelines.#", attrsPipelinesCount, count); err != nil {
			return err
		}
		if count > 0 {
			if name != "" {
				if err := util.CheckFieldEqualAndSet("pipelines.0.name", attrs["pipelines.0.name"], name); err != nil {
					return err
				}
			} else {
				if err := util.CheckFieldSet("pipelines.0.name", attrs["pipelines.0.name"]); err != nil {
					return err
				}
			}
			if ref != "" {
				if err := util.CheckFieldEqualAndSet("pipelines.0.refs.0", attrs["pipelines.0.refs.0"], ref); err != nil {
					return err
				}
			}
			if err := util.CheckFieldEqualAndSet("pipelines.0.git_changeset_base", attrs["pipelines.0.git_changeset_base"], buddy.PipelineGitChangeSetBaseLatestRun); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("pipelines.0.filesystem_changeset_base", attrs["pipelines.0.filesystem_changeset_base"], buddy.PipelineFilesystemChangeSetBaseDateModified); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("pipelines.0.description_required", attrsDescRequired, false); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("pipelines.0.concurrent_pipeline_runs", attrsConcurrentRun, false); err != nil {
				return err
			}
			if err := util.CheckFieldSet("pipelines.0.html_url", attrs["pipelines.0.html_url"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("pipelines.0.identifier", attrs["pipelines.0.identifier"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("pipelines.0.last_execution_status", attrs["pipelines.0.last_execution_status"]); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("pipelines.0.priority", attrs["pipelines.0.priority"], buddy.PipelinePriorityNormal); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("pipelines.0.cpu", attrs["pipelines.0.cpu"], buddy.PipelineCpuX64); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("pipelines.0.git_config_ref", attrs["pipelines.0.git_config_ref"], buddy.PipelineGitConfigRefNone); err != nil {
				return err
			}
			if err := util.CheckIntFieldSet("pipelines.0.pipeline_id", attrsPipelineId); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccSourcePipelinesConfig(domain string, projectName string, name1 string, name2 string, ref string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_pipeline" "a" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
   refs = ["%s"]
}

resource "buddy_pipeline" "b" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
   event {
       type = "PUSH"
       refs = ["%s"]
   }
}

data "buddy_pipelines" "all" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   depends_on = [buddy_pipeline.a, buddy_pipeline.b]
}

data "buddy_pipelines" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name_regex = "^aaaa"
   depends_on = [buddy_pipeline.a, buddy_pipeline.b]
}
`, domain, projectName, name1, ref, name2, ref)
}
