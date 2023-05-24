package test

//
//import (
//	"buddy-terraform/buddy/acc"
//	"buddy-terraform/buddy/util"
//	"fmt"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"strconv"
//	"testing"
//)
//
//func TestAccSourcePipeline(t *testing.T) {
//	domain := util.UniqueString()
//	projectName := util.UniqueString()
//	name := util.RandString(10)
//	ref := util.RandString(10)
//	reason := util.RandString(10)
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			// click
//			{
//				Config: testAccSourcePipelineConfigClick(domain, projectName, name, ref, buddy.PipelinePriorityHigh),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnClick, ref, "", buddy.PipelinePriorityHigh, false, ""),
//					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnClick, ref, "", buddy.PipelinePriorityHigh, false, ""),
//				),
//			},
//			// click disabled
//			{
//				Config: testAccSourcePipelineConfigClickDisabled(domain, projectName, name, ref, buddy.PipelinePriorityHigh, reason),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnClick, ref, "", buddy.PipelinePriorityHigh, true, reason),
//					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnClick, ref, "", buddy.PipelinePriorityHigh, true, reason),
//				),
//			},
//			// event
//			{
//				Config: testAccSourcePipelineConfigEvent(domain, projectName, name, ref, buddy.PipelinePriorityLow),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourcePipelineAttributes("data.buddy_pipeline.name", name, buddy.PipelineOnEvent, "", ref, buddy.PipelinePriorityLow, false, ""),
//					testAccSourcePipelineAttributes("data.buddy_pipeline.id", name, buddy.PipelineOnEvent, "", ref, buddy.PipelinePriorityLow, false, ""),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourcePipelineAttributes(n string, name string, on string, ref string, eventRef string, priority string, disabled bool, disabledReason string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsPipelineId, _ := strconv.Atoi(attrs["pipeline_id"])
//		attrsDisabled, _ := strconv.ParseBool(attrs["disabled"])
//		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckIntFieldSet("pipeline_id", attrsPipelineId); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("on", attrs["on"], on); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("last_execution_status", attrs["last_execution_status"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("priority", attrs["priority"], priority); err != nil {
//			return err
//		}
//		if ref != "" {
//			if err := util.CheckFieldEqualAndSet("refs.0", attrs["refs.0"], ref); err != nil {
//				return err
//			}
//		}
//		if eventRef != "" {
//			if err := util.CheckFieldEqualAndSet("event.0.refs.0", attrs["event.0.refs.0"], eventRef); err != nil {
//				return err
//			}
//		}
//		if disabled {
//			if err := util.CheckBoolFieldEqual("disabled", attrsDisabled, disabled); err != nil {
//				return err
//			}
//			if err := util.CheckFieldEqual("disabling_reason", attrs["disabling_reason"], disabledReason); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//}
//
//func testAccSourcePipelineConfigEvent(domain string, projectName string, name string, ref string, priority string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_project" "proj" {
//    domain = "${buddy_workspace.foo.domain}"
//    display_name = "%s"
//}
//
//resource "buddy_pipeline" "bar" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    name = "%s"
//    on = "EVENT"
//    event {
//        type = "PUSH"
//        refs = ["%s"]
//    }
//	priority = "%s"
//}
//
//data "buddy_pipeline" "name" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    name = "${buddy_pipeline.bar.name}"
//}
//
//data "buddy_pipeline" "id" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
//}
//`, domain, projectName, name, ref, priority)
//}
//
//func testAccSourcePipelineConfigClick(domain string, projectName string, name string, ref string, priority string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_project" "proj" {
//    domain = "${buddy_workspace.foo.domain}"
//    display_name = "%s"
//}
//
//resource "buddy_pipeline" "bar" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    name = "%s"
//    on = "CLICK"
//    refs = ["%s"]
//	priority = "%s"
//}
//
//data "buddy_pipeline" "name" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    name = "${buddy_pipeline.bar.name}"
//}
//
//data "buddy_pipeline" "id" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
//}
//`, domain, projectName, name, ref, priority)
//}
//
//func testAccSourcePipelineConfigClickDisabled(domain string, projectName string, name string, ref string, priority string, reason string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_project" "proj" {
//    domain = "${buddy_workspace.foo.domain}"
//    display_name = "%s"
//}
//
//resource "buddy_pipeline" "bar" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    name = "%s"
//    on = "CLICK"
//    refs = ["%s"]
//	priority = "%s"
//	disabled = true
//	disabling_reason = "%s"
//}
//
//data "buddy_pipeline" "name" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    name = "${buddy_pipeline.bar.name}"
//}
//
//data "buddy_pipeline" "id" {
//    domain = "${buddy_workspace.foo.domain}"
//    project_name = "${buddy_project.proj.name}"
//    pipeline_id = "${buddy_pipeline.bar.pipeline_id}"
//}
//`, domain, projectName, name, ref, priority, reason)
//}
