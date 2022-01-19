package util

import (
	"buddy-terraform/buddy/api"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CharSetAlpha = "abcdefghijklmnopqrstuvwxyz"
	SshKey       = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEA573Yk5eUJh2mVKYByEQqFt6mSQsES6HhlT9H4iDjQWtK2g78vDm1
Ni9MmUfMyJ1IYcNhOZW5CfPlueqXsUobZyo1eAMI1/0Ju5bSwZwRkdR2htydCJOX0mG0KZ
zT6NKNf7lg0R3bHoQcCcbHnvzE1F8wIXWsDJyeo0iXb8kilNr0kFAkMmjZGKjeXbRKR374
48TcrSaeGvSlOpsAa6YvBGoRnfhSkFysE+FxTFRhF7iknrvhjLqwbK7BR5Pf3j8hifiy3i
tMJM230COXwCAg2BoHyzH6xefP4TE6Po2qVfAcNmUzp+ktVbqf2HH44aFiZJgZYXJCTct/
RopyD0Uq3QAAA9D3L7tx9y+7cQAAAAdzc2gtcnNhAAABAQDnvdiTl5QmHaZUpgHIRCoW3q
ZJCwRLoeGVP0fiIONBa0raDvy8ObU2L0yZR8zInUhhw2E5lbkJ8+W56pexShtnKjV4AwjX
/Qm7ltLBnBGR1HaG3J0Ik5fSYbQpnNPo0o1/uWDRHdsehBwJxsee/MTUXzAhdawMnJ6jSJ
dvySKU2vSQUCQyaNkYqN5dtEpHfvjjxNytJp4a9KU6mwBrpi8EahGd+FKQXKwT4XFMVGEX
uKSeu+GMurBsrsFHk9/ePyGJ+LLeK0wkzbfQI5fAICDYGgfLMfrF58/hMTo+japV8Bw2ZT
On6S1Vup/YcfjhoWJkmBlhckJNy39GinIPRSrdAAAAAwEAAQAAAQAG860BkHSDSDRrKae4
CENy+C7o1gnE8xA/V+yiHfZzSfKu4/A0/U4wV+7mUj8UbZN0S1YpUhKA9+4WS7FNQjncOG
nuNbkYMaEPHZEo+bOVOlhr50ZWsYbGauPqs6evvlE8WaVL4KdoHPJyYKIwZMjKzig1eMA2
iKRBpbXVRqVg7bn0+opBUdv5FpsDkSa+ijJKLA7szSpM7yq03sZfx9u1/WTvh1Qa375mFC
O/8NUJoGii5Tacp+QeIHi2IEl+38eBLx1mal0AM68mkeJvU6Pfa+f/aUG4Z9RfQByD2/ql
JmrLVzaEv1Jy3n4lFnTsDT11dLmdQ7WjCNa/NKvVOrVJAAAAgCX4+x7xjcKAxEdprXpKkN
m4L+4Ciy1I03x/fk0GNDi316IiEUylTzTCy7HFAg0RUaxyN3iyFbwe1kN6DPV4r877WiOV
lt7v7hIjS1eeOwSGpsxijPxbWEIzOnu35t23YUROCOWaYXP9EAk1YEqEB34UKVq9jeB2t+
7zh/0ORkbwAAAAgQD8c9R/g6+Jx5c09MY0nX8IiTp1YBws6WVXDkaPPY//t1nKSfYEZy4p
ojOgWAS8vDMd7g24gr0Fm0lEKBXhgPVJTuClIH1P6IKRacLNtpp9uzVjBn/hi3f60GXo8y
A8ZbizwLkSpECLjFUer+ZlbrlcZ2Oq6o8DufiIjFAiqH2RMwAAAIEA6v+Cxku/IZ5XLrnm
wItoQx8eLu/Ly88wCSe1rgwzazgV4zQz1Y6B8ONqY03SgnE1im0/zaVmhMsRyE8FNGHPZP
8dclO3b8qod57BiY3PLpBKi9swhNQlxQn1zDhaF5cDEcTXYNL6xXC6iD6aa0bffjdBkuzz
A10SiVeaAxv3c68AAAAYTWljaGFsQE1pY2hhcy1pTWFjLmxvY2FsAQID
-----END OPENSSH PRIVATE KEY-----`
	SshKey2 = `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`
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
	var ret []string
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

func InterfaceIntToPointer(i interface{}) *int {
	return IntToPointer(i.(int))
}

func InterfaceBoolToPointer(i interface{}) *bool {
	return BoolToPointer(i.(bool))
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

func ApiShortGroupToMap(g *api.Group) map[string]interface{} {
	if g == nil {
		return nil
	}
	group := map[string]interface{}{}
	group["html_url"] = g.HtmlUrl
	group["group_id"] = g.Id
	group["name"] = g.Name
	return group
}

func ApiShortVariableToMap(v *api.Variable) map[string]interface{} {
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

func ApiShortWebhookToMap(w *api.Webhook) map[string]interface{} {
	if w == nil {
		return nil
	}
	webhook := map[string]interface{}{}
	webhook["target_url"] = w.TargetUrl
	webhook["webhook_id"] = w.Id
	webhook["html_url"] = w.HtmlUrl
	return webhook
}

func ApiShortVariableSshKeyToMap(v *api.Variable) map[string]interface{} {
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
	variable["display_name"] = v.FileName
	return variable
}

func ApiShortMemberToMap(m *api.Member) map[string]interface{} {
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

func ApiShortWorkspaceToMap(w *api.Workspace) map[string]interface{} {
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

func ApiShortProjectToMap(p *api.Project) map[string]interface{} {
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

func ApiShortPermissionToMap(permission *api.Permission) map[string]interface{} {
	if permission == nil {
		return nil
	}
	permissionMap := map[string]interface{}{}
	permissionMap["name"] = permission.Name
	permissionMap["pipeline_access_level"] = permission.PipelineAccessLevel
	permissionMap["repository_access_level"] = permission.RepositoryAccessLevel
	permissionMap["sandbox_access_level"] = permission.SandboxAccessLevel
	permissionMap["permission_id"] = permission.Id
	permissionMap["html_url"] = permission.HtmlUrl
	permissionMap["type"] = permission.Type
	return permissionMap
}

func ApiVariableSshKeyToResourceData(domain string, variable *api.Variable, d *schema.ResourceData, useValueProcessed bool) error {
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
	err = d.Set("display_name", variable.FileName)
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

func ApiVariableToResourceData(domain string, variable *api.Variable, d *schema.ResourceData, useValueProcessed bool) error {
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

func ApiWebhookToResourceData(domain string, webhook *api.Webhook, d *schema.ResourceData, short bool) error {
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

func ApiProjectGroupToResourceData(domain string, projectName string, group *api.ProjectGroup, d *schema.ResourceData, setParentPermissionId bool) error {
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

func ApiProjectMemberToResourceData(domain string, projectName string, member *api.ProjectMember, d *schema.ResourceData, setParentPermissionId bool) error {
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

func ApiProjectToResourceData(domain string, project *api.Project, d *schema.ResourceData, short bool) error {
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
		err = d.Set("ssh_public_key", project.SshPublicKey)
		if err != nil {
			return err
		}
		err = d.Set("key_fingerprint", project.KeyFingerprint)
		if err != nil {
			return err
		}
		return d.Set("default_branch", project.DefaultBranch)
	}
	return nil
}

func ApiPermissionToResourceData(domain string, p *api.Permission, d *schema.ResourceData) error {
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

func ApiWorkspaceToResourceData(workspace *api.Workspace, d *schema.ResourceData, short bool) error {
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

func ApiProfileEmailToResourceData(p *api.ProfileEmail, d *schema.ResourceData) error {
	d.SetId(p.Email)
	err := d.Set("email", p.Email)
	if err != nil {
		return err
	}
	return d.Set("confirmed", p.Confirmed)
}

func ApiPublicKeyToResourceData(k *api.PublicKey, d *schema.ResourceData) error {
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

func ApiProfileToResourceData(p *api.Profile, d *schema.ResourceData) error {
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

func ApiMemberToResourceData(domain string, m *api.Member, d *schema.ResourceData) error {
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
	return d.Set("workspace_owner", m.WorkspaceOwner)
}

func ApiGroupToResourceData(domain string, g *api.Group, d *schema.ResourceData) error {
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
	return d.Set("description", g.Description)
}

func ApiGroupMemberToResourceData(domain string, groupId int, m *api.Member, d *schema.ResourceData) error {
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
