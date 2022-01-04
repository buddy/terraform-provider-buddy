package util

import (
	"buddy-terraform/buddy/api"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
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

func InterfaceStringToPointer(i interface{}) *string {
	return StringToPointer(i.(string))
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
