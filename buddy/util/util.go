package util

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"golang.org/x/crypto/ssh"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CharSetAlpha = "abcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func CheckFieldEqual(field string, got string, want string) error {
	if got != want {
		return ErrorFieldFormatted(field, got, want)
	}
	return nil
}

func GenerateRsaKeyPair() (error, string, string) {
	privateKey, err := rsa.GenerateKey(crand.Reader, 4096)
	if err != nil {
		return err, "", ""
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyBytesEncoded := pem.EncodeToMemory(privateKeyBlock)
	sshPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err, "", ""
	}
	sshPublicKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	return nil, strings.TrimSpace(string(sshPublicKeyBytes)), strings.TrimSpace(string(privateKeyBytesEncoded))
}

func CheckDateFieldEqual(field string, got string, want string) error {
	gotDate, err := time.Parse(time.RFC3339, got)
	if err != nil {
		return err
	}
	wantDate, err := time.Parse(time.RFC3339, want)
	if err != nil {
		return err
	}
	gotDate = gotDate.Truncate(time.Second)
	wantDate = wantDate.Truncate(time.Second)
	if gotDate.Equal(wantDate) {
		return nil
	}
	return ErrorFieldFormatted(field, gotDate.String(), wantDate.String())
}

func CheckFieldEqualAndSet(field string, got string, want string) error {
	if err := CheckFieldEqual(field, got, want); err != nil {
		return err
	}
	return CheckFieldSet(field, got)
}

func CheckFieldSet(field string, got string) error {
	if got == "" {
		return ErrorFieldEmpty(field)
	}
	return nil
}

func CheckBoolFieldEqual(field string, got bool, want bool) error {
	if got != want {
		return ErrorFieldFormatted(field, strconv.FormatBool(got), strconv.FormatBool(want))
	}
	return nil
}

func CheckIntFieldEqual(field string, got int, want int) error {
	if got != want {
		return ErrorFieldFormatted(field, strconv.Itoa(got), strconv.Itoa(want))
	}
	return nil
}

func CheckIntFieldEqualAndSet(field string, got int, want int) error {
	if err := CheckIntFieldEqual(field, got, want); err != nil {
		return err
	}
	return CheckIntFieldSet(field, got)
}

func CheckIntFieldSet(field string, got int) error {
	if got == 0 {
		return ErrorFieldEmpty(field)
	}
	return nil
}

func ErrorFieldFormatted(field string, got string, want string) error {
	return fmt.Errorf("got %q %q; want %q", field, got, want)
}

func ErrorFieldEmpty(field string) error {
	return fmt.Errorf("expected %q not to be empty", field)
}

func ErrorResourceExists() error {
	return errors.New("resource still exists")
}

func ComposeDoubleId(a, b string) string {
	return fmt.Sprintf("%s:%s", a, b)
}

func ComposeTripleId(a, b, c string) string {
	return fmt.Sprintf("%s:%s:%s", a, b, c)
}

func StringToPointer(p string) *string {
	s := new(string)
	*s = p
	return s
}

func StringSetToStringSlice(stringSet *schema.Set) []string {
	ret := make([]string, 0)
	if stringSet == nil {
		return ret
	}
	for _, envVal := range stringSet.List() {
		ret = append(ret, envVal.(string))
	}
	return ret
}

func InterfaceStringToPointer(i interface{}) *string {
	return StringToPointer(i.(string))
}

func InterfaceStringSetToStringSlice(i interface{}) []string {
	return StringSetToStringSlice(i.(*schema.Set))
}

func InterfaceStringSetToPointer(i interface{}) *[]string {
	a := StringSetToStringSlice(i.(*schema.Set))
	return &a
}

func InterfaceIntToPointer(i interface{}) *int {
	return IntToPointer(i.(int))
}

func InterfaceBoolToPointer(i interface{}) *bool {
	return BoolToPointer(i.(bool))
}

func IsBoolPointerSet(i interface{}) bool {
	return i != nil
}

func BoolToPointer(p bool) *bool {
	b := new(bool)
	*b = p
	return b
}

func IntToPointer(p int) *int {
	i := new(int)
	*i = p
	return i
}

func IntSetToIntSlice(intSet *schema.Set) []int {
	var ret []int
	if intSet == nil {
		return ret
	}
	for _, envVal := range intSet.List() {
		ret = append(ret, envVal.(int))
	}
	return ret
}

func DecomposeDoubleId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("wrong id format %q", id)
	}
	return parts[0], parts[1], nil
}

func DecomposeTripleId(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("wrong id format %q", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func RandStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}

func RandInt() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func RandString(strlen int) string {
	return RandStringFromCharSet(strlen, CharSetAlpha)
}

func RandEmail() string {
	return fmt.Sprintf("%s@0zxc.com", UniqueString())
}

func UniqueString() string {
	return fmt.Sprintf("%s%d", RandString(5), time.Now().UnixNano())
}

func ValidateDomain(v interface{}, _ string) (we []string, err []error) {
	value := v.(string)
	length := len(value)
	if length < 4 {
		err = append(err, errors.New("domain must have at least 4 characters"))
	} else if length > 100 {
		err = append(err, errors.New("domain cannot be longer than 100 characters"))
	}
	match, _ := regexp.MatchString("^[a-z0-9][a-z0-9\\-_]+[a-z0-9]$", value)
	if !match {
		err = append(err, errors.New("domain must be lowercase and contain only letters, numbers or dash ( - ) and footer ( _ ) characters. It must start and end with a letter or number"))
	}
	return
}

func TestSleep(ms int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	}
}

func ValidateEmail(v interface{}, _ string) (we []string, err []error) {
	value := v.(string)
	match, _ := regexp.MatchString("(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,}$", value)
	if !match {
		err = append(err, errors.New("email is not valid"))
	}
	return
}

func InterfacePipelinePermissionToPointer(i interface{}) *buddy.PipelinePermissions {
	var result *buddy.PipelinePermissions
	if i != nil {
		l := i.([]interface{})
		if len(l) != 1 {
			return result
		}
		result = &buddy.PipelinePermissions{}
		m := l[0].(map[string]interface{})
		if m["others"] != nil && m["others"] != "" {
			result.Others = m["others"].(string)
		} else {
			result.Others = buddy.PipelinePermissionDefault
		}
		result.Users = []*buddy.PipelineResourcePermission{}
		result.Groups = []*buddy.PipelineResourcePermission{}
		if m["user"] != nil {
			for _, v := range m["user"].([]interface{}) {
				u := v.(map[string]interface{})
				ua := buddy.PipelineResourcePermission{
					Id:          u["id"].(int),
					AccessLevel: u["access_level"].(string),
				}
				result.Users = append(result.Users, &ua)
			}
		}
		if m["group"] != nil {
			for _, v := range m["group"].([]interface{}) {
				g := v.(map[string]interface{})
				ga := buddy.PipelineResourcePermission{
					Id:          g["id"].(int),
					AccessLevel: g["access_level"].(string),
				}
				result.Groups = append(result.Groups, &ga)
			}
		}
	}
	return result
}

func MapTriggerConditionsToApi(l interface{}) *[]*buddy.PipelineTriggerCondition {
	var expanded []*buddy.PipelineTriggerCondition
	for _, v := range l.([]interface{}) {
		m := v.(map[string]interface{})
		c := buddy.PipelineTriggerCondition{
			TriggerCondition: m["condition"].(string),
		}
		if m["paths"] != nil {
			c.TriggerConditionPaths = InterfaceStringSetToStringSlice(m["paths"])
		}
		if m["variable_key"] != nil {
			c.TriggerVariableKey = m["variable_key"].(string)
		}
		if m["variable_value"] != nil {
			c.TriggerVariableValue = m["variable_value"].(string)
		}
		if m["hours"] != nil {
			c.TriggerHours = IntSetToIntSlice(m["hours"].(*schema.Set))
		}
		if m["days"] != nil {
			c.TriggerDays = IntSetToIntSlice(m["days"].(*schema.Set))
		}
		if m["zone_id"] != nil {
			c.ZoneId = m["zone_id"].(string)
		}
		if m["project_name"] != nil {
			c.TriggerProjectName = m["project_name"].(string)
		}
		if m["pipeline_name"] != nil {
			c.TriggerPipelineName = m["pipeline_name"].(string)
		}
		expanded = append(expanded, &c)
	}
	return &expanded
}

func MapPipelineEventsToApi(l interface{}) *[]*buddy.PipelineEvent {
	var expanded []*buddy.PipelineEvent
	for _, v := range l.([]interface{}) {
		m := v.(map[string]interface{})
		e := buddy.PipelineEvent{
			Type: m["type"].(string),
			Refs: InterfaceStringSetToStringSlice(m["refs"]),
		}
		expanded = append(expanded, &e)
	}
	return &expanded
}

func MapPipelineRemoteParametersToApi(l interface{}) *[]*buddy.PipelineRemoteParameter {
	var expanded []*buddy.PipelineRemoteParameter
	for _, v := range l.([]interface{}) {
		m := v.(map[string]interface{})
		e := buddy.PipelineRemoteParameter{
			Key:   m["key"].(string),
			Value: m["value"].(string),
		}
		expanded = append(expanded, &e)
	}
	return &expanded
}

func ApiShortIntegrationToMap(i *buddy.Integration) map[string]interface{} {
	if i == nil {
		return nil
	}
	integration := map[string]interface{}{}
	integration["html_url"] = i.HtmlUrl
	integration["integration_id"] = i.HashId
	integration["name"] = i.Name
	integration["type"] = i.Type
	return integration
}

func ApiTriggerConditionToMap(c *buddy.PipelineTriggerCondition) map[string]interface{} {
	if c == nil {
		return nil
	}
	condition := map[string]interface{}{}
	condition["condition"] = c.TriggerCondition
	condition["paths"] = c.TriggerConditionPaths
	condition["variable_key"] = c.TriggerVariableKey
	condition["variable_value"] = c.TriggerVariableValue
	condition["hours"] = c.TriggerHours
	condition["days"] = c.TriggerDays
	condition["zone_id"] = c.ZoneId
	condition["project_name"] = c.TriggerProjectName
	condition["pipeline_name"] = c.TriggerPipelineName
	return condition
}

func ApiPipelineResourcePermissionToMap(p *buddy.PipelineResourcePermission) map[string]interface{} {
	if p == nil {
		return nil
	}
	permission := map[string]interface{}{}
	permission["id"] = p.Id
	permission["access_level"] = p.AccessLevel
	return permission
}

func ApiPipelineResourcePermissionsToMap(p []*buddy.PipelineResourcePermission) []interface{} {
	var list []interface{}
	for _, v := range p {
		list = append(list, ApiPipelineResourcePermissionToMap(v))
	}
	return list
}

func ApiPipelinePermissionsToMap(p *buddy.PipelinePermissions) []interface{} {
	if p == nil {
		return nil
	}
	permissions := map[string]interface{}{}
	permissions["others"] = p.Others
	permissions["user"] = ApiPipelineResourcePermissionsToMap(p.Users)
	permissions["group"] = ApiPipelineResourcePermissionsToMap(p.Groups)
	list := []interface{}{
		permissions,
	}
	return list
}

func ApiPipelineTriggerConditionsToMap(l []*buddy.PipelineTriggerCondition) []interface{} {
	if l == nil {
		return nil
	}
	var list []interface{}
	for _, c := range l {
		list = append(list, ApiTriggerConditionToMap(c))
	}
	return list
}

func ApiPipelineRemoteParameterToMap(e *buddy.PipelineRemoteParameter) map[string]interface{} {
	if e == nil {
		return nil
	}
	param := map[string]interface{}{}
	param["key"] = e.Key
	param["value"] = e.Value
	return param
}

func ApiPipelineEventToMap(e *buddy.PipelineEvent) map[string]interface{} {
	if e == nil {
		return nil
	}
	event := map[string]interface{}{}
	event["type"] = e.Type
	event["refs"] = e.Refs
	return event
}

func ApiPipelineEventsToMap(l []*buddy.PipelineEvent) []interface{} {
	if l == nil {
		return nil
	}
	var list []interface{}
	for _, e := range l {
		list = append(list, ApiPipelineEventToMap(e))
	}
	return list
}

func ApiPipelineRemoteParametersToMap(l []*buddy.PipelineRemoteParameter) []interface{} {
	if l == nil {
		return nil
	}
	var list []interface{}
	for _, e := range l {
		list = append(list, ApiPipelineRemoteParameterToMap(e))
	}
	return list
}

func ApiPipelineToResourceData(domain string, projectName string, pipeline *buddy.Pipeline, d *schema.ResourceData, short bool) error {
	d.SetId(ComposeTripleId(domain, projectName, strconv.Itoa(pipeline.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("project_name", projectName)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("always_from_scratch", pipeline.AlwaysFromScratch)
		if err != nil {
			return err
		}
		err = d.Set("fail_on_prepare_env_warning", pipeline.FailOnPrepareEnvWarning)
		if err != nil {
			return err
		}
		err = d.Set("fetch_all_refs", pipeline.FetchAllRefs)
		if err != nil {
			return err
		}
		err = d.Set("auto_clear_cache", pipeline.AutoClearCache)
		if err != nil {
			return err
		}
		err = d.Set("no_skip_to_most_recent", pipeline.NoSkipToMostRecent)
		if err != nil {
			return err
		}
		err = d.Set("do_not_create_commit_status", pipeline.DoNotCreateCommitStatus)
		if err != nil {
			return err
		}
		err = d.Set("start_date", pipeline.StartDate)
		if err != nil {
			return err
		}
		err = d.Set("delay", pipeline.Delay)
		if err != nil {
			return err
		}
		err = d.Set("clone_depth", pipeline.CloneDepth)
		if err != nil {
			return err
		}
		err = d.Set("cron", pipeline.Cron)
		if err != nil {
			return err
		}
		err = d.Set("paused", pipeline.Paused)
		if err != nil {
			return err
		}
		err = d.Set("ignore_fail_on_project_status", pipeline.IgnoreFailOnProjectStatus)
		if err != nil {
			return err
		}
		err = d.Set("execution_message_template", pipeline.ExecutionMessageTemplate)
		if err != nil {
			return err
		}
		err = d.Set("worker", pipeline.Worker)
		if err != nil {
			return err
		}
		err = d.Set("target_site_url", pipeline.TargetSiteUrl)
		if err != nil {
			return err
		}
		err = d.Set("create_date", pipeline.CreateDate)
		if err != nil {
			return err
		}
		err = d.Set("creator", []interface{}{ApiShortMemberToMap(pipeline.Creator)})
		if err != nil {
			return err
		}
		err = d.Set("project", []interface{}{ApiShortProjectToMap(pipeline.Project)})
		if err != nil {
			return err
		}
		err = d.Set("trigger_condition", ApiPipelineTriggerConditionsToMap(pipeline.TriggerConditions))
		if err != nil {
			return err
		}
		err = d.Set("permissions", ApiPipelinePermissionsToMap(pipeline.Permissions))
		if err != nil {
			return err
		}
	}
	err = d.Set("priority", pipeline.Priority)
	if err != nil {
		return err
	}
	err = d.Set("html_url", pipeline.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("pipeline_id", pipeline.Id)
	if err != nil {
		return err
	}
	err = d.Set("name", pipeline.Name)
	if err != nil {
		return err
	}
	err = d.Set("on", pipeline.On)
	if err != nil {
		return err
	}
	err = d.Set("last_execution_status", pipeline.LastExecutionStatus)
	if err != nil {
		return err
	}
	err = d.Set("last_execution_revision", pipeline.LastExecutionRevision)
	if err != nil {
		return err
	}
	err = d.Set("refs", pipeline.Refs)
	if err != nil {
		return err
	}
	err = d.Set("tags", pipeline.Tags)
	if err != nil {
		return err
	}
	definitionSource := pipeline.DefinitionSource
	if definitionSource == "" {
		definitionSource = buddy.PipelineDefinitionSourceLocal
	}
	err = d.Set("definition_source", definitionSource)
	if err != nil {
		return err
	}
	err = d.Set("remote_path", pipeline.RemotePath)
	if err != nil {
		return err
	}
	err = d.Set("remote_branch", pipeline.RemoteBranch)
	if err != nil {
		return err
	}
	err = d.Set("remote_project_name", pipeline.RemoteProjectName)
	if err != nil {
		return err
	}
	err = d.Set("remote_parameter", ApiPipelineRemoteParametersToMap(pipeline.RemoteParameters))
	if err != nil {
		return err
	}
	err = d.Set("disabled", pipeline.Disabled)
	if err != nil {
		return err
	}
	err = d.Set("disabling_reason", pipeline.DisabledReason)
	if err != nil {
		return err
	}
	return d.Set("event", ApiPipelineEventsToMap(pipeline.Events))
}

func ApiIntegrationToResourceData(domain string, i *buddy.Integration, d *schema.ResourceData, short bool) error {
	d.SetId(ComposeDoubleId(domain, i.HashId))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("name", i.Name)
	if err != nil {
		return err
	}
	err = d.Set("type", i.Type)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("scope", i.Scope)
		if err != nil {
			return err
		}
		err = d.Set("project_name", i.ProjectName)
		if err != nil {
			return err
		}
		err = d.Set("group_id", i.GroupId)
		if err != nil {
			return err
		}
	}
	err = d.Set("integration_id", i.HashId)
	if err != nil {
		return err
	}
	return d.Set("html_url", i.HtmlUrl)
}

func MapRoleAssumptionsToApi(l interface{}) *[]*buddy.RoleAssumption {
	var expanded []*buddy.RoleAssumption
	for _, v := range l.([]interface{}) {
		m := v.(map[string]interface{})
		r := buddy.RoleAssumption{
			Arn: m["arn"].(string),
		}
		if m["external_id"] != nil {
			r.ExternalId = m["external_id"].(string)
		}
		if m["duration"] != nil {
			r.Duration = m["duration"].(int)
		}
		expanded = append(expanded, &r)
	}
	return &expanded
}

func ApiShortGroupToMap(g *buddy.Group) map[string]interface{} {
	if g == nil {
		return nil
	}
	group := map[string]interface{}{}
	group["html_url"] = g.HtmlUrl
	group["group_id"] = g.Id
	group["name"] = g.Name
	return group
}

func ApiShortVariableToMap(v *buddy.Variable) map[string]interface{} {
	if v == nil {
		return nil
	}
	variable := map[string]interface{}{}
	variable["key"] = v.Key
	variable["encrypted"] = v.Encrypted
	variable["settable"] = v.Settable
	variable["description"] = v.Description
	variable["value"] = v.Value
	variable["variable_id"] = v.Id
	return variable
}

func ApiShortWebhookToMap(w *buddy.Webhook) map[string]interface{} {
	if w == nil {
		return nil
	}
	webhook := map[string]interface{}{}
	webhook["target_url"] = w.TargetUrl
	webhook["webhook_id"] = w.Id
	webhook["html_url"] = w.HtmlUrl
	return webhook
}

func ApiShortVariableSshKeyToMap(v *buddy.Variable) map[string]interface{} {
	if v == nil {
		return nil
	}
	variable := map[string]interface{}{}
	variable["key"] = v.Key
	variable["encrypted"] = v.Encrypted
	variable["settable"] = v.Settable
	variable["description"] = v.Description
	variable["value"] = v.Value
	variable["variable_id"] = v.Id
	variable["public_value"] = v.PublicValue
	variable["key_fingerprint"] = v.KeyFingerprint
	variable["checksum"] = v.Checksum
	variable["file_chmod"] = v.FileChmod
	variable["file_path"] = v.FilePath
	variable["file_place"] = v.FilePlace
	return variable
}

func ApiShortMemberToMap(m *buddy.Member) map[string]interface{} {
	if m == nil {
		return nil
	}
	member := map[string]interface{}{}
	member["html_url"] = m.HtmlUrl
	member["name"] = m.Name
	member["email"] = m.Email
	member["avatar_url"] = m.AvatarUrl
	member["member_id"] = m.Id
	member["admin"] = m.Admin
	member["workspace_owner"] = m.WorkspaceOwner
	return member
}

func ApiShortGroupMemberToMap(m *buddy.Member) map[string]interface{} {
	member := ApiShortMemberToMap(m)
	if member != nil {
		member["status"] = m.Status
	}
	return member
}

func ApiShortWorkspaceToMap(w *buddy.Workspace) map[string]interface{} {
	if w == nil {
		return nil
	}
	workspace := map[string]interface{}{}
	workspace["html_url"] = w.HtmlUrl
	workspace["workspace_id"] = w.Id
	workspace["name"] = w.Name
	workspace["domain"] = w.Domain
	return workspace
}

func ApiShortProjectToMap(p *buddy.Project) map[string]interface{} {
	if p == nil {
		return nil
	}
	project := map[string]interface{}{}
	project["html_url"] = p.HtmlUrl
	project["name"] = p.Name
	project["display_name"] = p.DisplayName
	project["status"] = p.Status
	return project
}

func ApiShortPermissionToMap(permission *buddy.Permission) map[string]interface{} {
	if permission == nil {
		return nil
	}
	permissionMap := map[string]interface{}{}
	permissionMap["name"] = permission.Name
	permissionMap["pipeline_access_level"] = permission.PipelineAccessLevel
	permissionMap["repository_access_level"] = permission.RepositoryAccessLevel
	permissionMap["sandbox_access_level"] = permission.SandboxAccessLevel
	permissionMap["project_team_access_level"] = permission.ProjectTeamAccessLevel
	permissionMap["permission_id"] = permission.Id
	permissionMap["html_url"] = permission.HtmlUrl
	permissionMap["type"] = permission.Type
	return permissionMap
}

func ApiShortPipelineToMap(p *buddy.Pipeline) map[string]interface{} {
	if p == nil {
		return nil
	}
	pipeline := map[string]interface{}{}
	pipeline["name"] = p.Name
	pipeline["pipeline_id"] = p.Id
	pipeline["html_url"] = p.HtmlUrl
	pipeline["on"] = p.On
	pipeline["priority"] = p.Priority
	pipeline["last_execution_status"] = p.LastExecutionStatus
	pipeline["last_execution_revision"] = p.LastExecutionRevision
	pipeline["disabled"] = p.Disabled
	pipeline["disabling_reason"] = p.DisabledReason
	pipeline["refs"] = p.Refs
	pipeline["tags"] = p.Tags
	pipeline["event"] = ApiPipelineEventsToMap(p.Events)
	definitionSource := p.DefinitionSource
	if definitionSource == "" {
		definitionSource = buddy.PipelineDefinitionSourceLocal
	}
	pipeline["definition_source"] = definitionSource
	pipeline["remote_path"] = p.RemotePath
	pipeline["remote_branch"] = p.RemoteBranch
	pipeline["remote_project_name"] = p.RemoteProjectName
	pipeline["remote_parameter"] = ApiPipelineRemoteParametersToMap(p.RemoteParameters)
	return pipeline
}

func ApiVariableSshKeyToResourceData(domain string, variable *buddy.Variable, d *schema.ResourceData, useValueProcessed bool) error {
	d.SetId(ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("key", variable.Key)
	if err != nil {
		return err
	}
	if useValueProcessed {
		err = d.Set("value_processed", variable.Value)
	} else {
		err = d.Set("value", variable.Value)
	}
	if err != nil {
		return err
	}
	err = d.Set("encrypted", variable.Encrypted)
	if err != nil {
		return err
	}
	err = d.Set("settable", variable.Settable)
	if err != nil {
		return err
	}
	err = d.Set("description", variable.Description)
	if err != nil {
		return err
	}
	err = d.Set("file_place", variable.FilePlace)
	if err != nil {
		return err
	}
	err = d.Set("file_path", variable.FilePath)
	if err != nil {
		return err
	}
	err = d.Set("file_chmod", variable.FileChmod)
	if err != nil {
		return err
	}
	err = d.Set("variable_id", variable.Id)
	if err != nil {
		return err
	}
	err = d.Set("checksum", variable.Checksum)
	if err != nil {
		return err
	}
	err = d.Set("key_fingerprint", variable.KeyFingerprint)
	if err != nil {
		return err
	}
	return d.Set("public_value", variable.PublicValue)
}

func ApiVariableToResourceData(domain string, variable *buddy.Variable, d *schema.ResourceData, useValueProcessed bool) error {
	d.SetId(ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("key", variable.Key)
	if err != nil {
		return err
	}
	if useValueProcessed {
		err = d.Set("value_processed", variable.Value)
	} else {
		err = d.Set("value", variable.Value)
	}
	if err != nil {
		return err
	}
	err = d.Set("encrypted", variable.Encrypted)
	if err != nil {
		return err
	}
	err = d.Set("settable", variable.Settable)
	if err != nil {
		return err
	}
	err = d.Set("description", variable.Description)
	if err != nil {
		return err
	}
	return d.Set("variable_id", variable.Id)
}

func ApiWebhookToResourceData(domain string, webhook *buddy.Webhook, d *schema.ResourceData, short bool) error {
	d.SetId(ComposeDoubleId(domain, strconv.Itoa(webhook.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("target_url", webhook.TargetUrl)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("secret_key", webhook.SecretKey)
		if err != nil {
			return err
		}
		err = d.Set("projects", webhook.Projects)
		if err != nil {
			return err
		}
		err = d.Set("events", webhook.Events)
		if err != nil {
			return err
		}
	}
	err = d.Set("webhook_id", webhook.Id)
	if err != nil {
		return err
	}
	return d.Set("html_url", webhook.HtmlUrl)
}

func ApiProjectGroupToResourceData(domain string, projectName string, group *buddy.ProjectGroup, d *schema.ResourceData, setParentPermissionId bool) error {
	d.SetId(ComposeTripleId(domain, projectName, strconv.Itoa(group.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("project_name", projectName)
	if err != nil {
		return err
	}
	err = d.Set("group_id", group.Id)
	if err != nil {
		return err
	}
	if setParentPermissionId {
		err = d.Set("permission_id", group.PermissionSet.Id)
		if err != nil {
			return err
		}
	}
	err = d.Set("html_url", group.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("name", group.Name)
	if err != nil {
		return err
	}
	return d.Set("permission", []interface{}{ApiShortPermissionToMap(group.PermissionSet)})
}

func ApiProjectMemberToResourceData(domain string, projectName string, member *buddy.ProjectMember, d *schema.ResourceData, setParentPermissionId bool) error {
	d.SetId(ComposeTripleId(domain, projectName, strconv.Itoa(member.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("project_name", projectName)
	if err != nil {
		return err
	}
	err = d.Set("member_id", member.Id)
	if err != nil {
		return err
	}
	if setParentPermissionId {
		err = d.Set("permission_id", member.PermissionSet.Id)
		if err != nil {
			return err
		}
	}
	err = d.Set("html_url", member.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("name", member.Name)
	if err != nil {
		return err
	}
	err = d.Set("email", member.Email)
	if err != nil {
		return err
	}
	err = d.Set("avatar_url", member.AvatarUrl)
	if err != nil {
		return err
	}
	err = d.Set("admin", member.Admin)
	if err != nil {
		return err
	}
	err = d.Set("workspace_owner", member.WorkspaceOwner)
	if err != nil {
		return err
	}
	return d.Set("permission", []interface{}{ApiShortPermissionToMap(member.PermissionSet)})
}

func ApiProjectToResourceData(domain string, project *buddy.Project, d *schema.ResourceData, short bool) error {
	d.SetId(ComposeDoubleId(domain, project.Name))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("html_url", project.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("name", project.Name)
	if err != nil {
		return err
	}
	err = d.Set("display_name", project.DisplayName)
	if err != nil {
		return err
	}
	err = d.Set("status", project.Status)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("create_date", project.CreateDate)
		if err != nil {
			return err
		}
		err = d.Set("created_by", []interface{}{ApiShortMemberToMap(project.CreatedBy)})
		if err != nil {
			return err
		}
		err = d.Set("http_repository", project.HttpRepository)
		if err != nil {
			return err
		}
		err = d.Set("ssh_repository", project.SshRepository)
		if err != nil {
			return err
		}
		err = d.Set("update_default_branch_from_external", project.UpdateDefaultBranchFromExternal)
		if err != nil {
			return err
		}
		err = d.Set("allow_pull_requests", project.AllowPullRequests)
		if err != nil {
			return err
		}
		err = d.Set("access", project.Access)
		if err != nil {
			return err
		}
		err = d.Set("fetch_submodules", project.FetchSubmodules)
		if err != nil {
			return err
		}
		err = d.Set("fetch_submodules_env_key", project.FetchSubmodulesEnvKey)
		if err != nil {
			return err
		}
		return d.Set("default_branch", project.DefaultBranch)
	}
	return nil
}

func ApiPermissionToResourceData(domain string, p *buddy.Permission, d *schema.ResourceData) error {
	d.SetId(ComposeDoubleId(domain, strconv.Itoa(p.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("name", p.Name)
	if err != nil {
		return err
	}
	err = d.Set("permission_id", p.Id)
	if err != nil {
		return err
	}
	err = d.Set("pipeline_access_level", p.PipelineAccessLevel)
	if err != nil {
		return err
	}
	err = d.Set("repository_access_level", p.RepositoryAccessLevel)
	if err != nil {
		return err
	}
	err = d.Set("project_team_access_level", p.ProjectTeamAccessLevel)
	if err != nil {
		return err
	}
	err = d.Set("sandbox_access_level", p.SandboxAccessLevel)
	if err != nil {
		return err
	}
	err = d.Set("description", p.Description)
	if err != nil {
		return err
	}
	err = d.Set("html_url", p.HtmlUrl)
	if err != nil {
		return err
	}
	return d.Set("type", p.Type)
}

func ApiWorkspaceToResourceData(workspace *buddy.Workspace, d *schema.ResourceData, short bool) error {
	d.SetId(workspace.Domain)
	err := d.Set("domain", workspace.Domain)
	if err != nil {
		return err
	}
	err = d.Set("workspace_id", workspace.Id)
	if err != nil {
		return err
	}
	err = d.Set("html_url", workspace.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("name", workspace.Name)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("owner_id", workspace.OwnerId)
		if err != nil {
			return err
		}
		err = d.Set("frozen", workspace.Frozen)
		if err != nil {
			return err
		}
		return d.Set("create_date", workspace.CreateDate)
	}
	return nil
}

func ApiProfileEmailToResourceData(p *buddy.ProfileEmail, d *schema.ResourceData) error {
	d.SetId(p.Email)
	err := d.Set("email", p.Email)
	if err != nil {
		return err
	}
	return d.Set("confirmed", p.Confirmed)
}

func ApiPublicKeyToResourceData(k *buddy.PublicKey, d *schema.ResourceData) error {
	d.SetId(strconv.Itoa(k.Id))
	err := d.Set("content", k.Content)
	if err != nil {
		return err
	}
	err = d.Set("html_url", k.HtmlUrl)
	if err != nil {
		return err
	}
	return d.Set("title", k.Title)
}

func ApiProfileToResourceData(p *buddy.Profile, d *schema.ResourceData) error {
	d.SetId("me")
	err := d.Set("member_id", p.Id)
	if err != nil {
		return err
	}
	err = d.Set("html_url", p.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("name", p.Name)
	if err != nil {
		return err
	}
	return d.Set("avatar_url", p.AvatarUrl)
}

func ApiMemberToResourceData(domain string, m *buddy.Member, d *schema.ResourceData, short bool) error {
	d.SetId(ComposeDoubleId(domain, strconv.Itoa(m.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("name", m.Name)
	if err != nil {
		return err
	}
	err = d.Set("member_id", m.Id)
	if err != nil {
		return err
	}
	err = d.Set("email", m.Email)
	if err != nil {
		return err
	}
	err = d.Set("html_url", m.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("avatar_url", m.AvatarUrl)
	if err != nil {
		return err
	}
	err = d.Set("admin", m.Admin)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("auto_assign_to_new_projects", m.AutoAssignToNewProjects)
		if err != nil {
			return err
		}
		if m.AutoAssignToNewProjects {
			err = d.Set("auto_assign_permission_set_id", m.AutoAssignPermissionSetId)
			if err != nil {
				return err
			}
		}
	}
	return d.Set("workspace_owner", m.WorkspaceOwner)
}

func ApiGroupToResourceData(domain string, g *buddy.Group, d *schema.ResourceData, short bool) error {
	d.SetId(ComposeDoubleId(domain, strconv.Itoa(g.Id)))
	err := d.Set("name", g.Name)
	if err != nil {
		return err
	}
	err = d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("group_id", g.Id)
	if err != nil {
		return err
	}
	err = d.Set("html_url", g.HtmlUrl)
	if err != nil {
		return err
	}
	if !short {
		err = d.Set("auto_assign_to_new_projects", g.AutoAssignToNewProjects)
		if err != nil {
			return err
		}
		if g.AutoAssignToNewProjects {
			err = d.Set("auto_assign_permission_set_id", g.AutoAssignPermissionSetId)
			if err != nil {
				return err
			}
		}
	}
	return d.Set("description", g.Description)
}

func ApiGroupMemberToResourceData(domain string, groupId int, m *buddy.Member, d *schema.ResourceData) error {
	d.SetId(ComposeTripleId(domain, strconv.Itoa(groupId), strconv.Itoa(m.Id)))
	err := d.Set("domain", domain)
	if err != nil {
		return err
	}
	err = d.Set("group_id", groupId)
	if err != nil {
		return err
	}
	err = d.Set("member_id", m.Id)
	if err != nil {
		return err
	}
	err = d.Set("html_url", m.HtmlUrl)
	if err != nil {
		return err
	}
	err = d.Set("status", m.Status)
	if err != nil {
		return err
	}
	err = d.Set("name", m.Name)
	if err != nil {
		return err
	}
	err = d.Set("email", m.Email)
	if err != nil {
		return err
	}
	err = d.Set("avatar_url", m.AvatarUrl)
	if err != nil {
		return err
	}
	err = d.Set("admin", m.Admin)
	if err != nil {
		return err
	}
	return d.Set("workspace_owner", m.WorkspaceOwner)
}

func IsResourceNotFound(resp *http.Response, err error) bool {
	if resp.StatusCode == http.StatusNotFound {
		return true
	}
	if resp.StatusCode == http.StatusForbidden && strings.Contains(err.Error(), "Only active workspace have access to API") {
		return true
	}
	return false
}
