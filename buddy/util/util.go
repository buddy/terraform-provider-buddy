package util

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"golang.org/x/crypto/ssh"
	"math/big"
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

func NewDiagnosticApiError(method string, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic("Buddy API error occured", fmt.Sprintf("Unable to %s:\n%s", method, err.Error()))
}

func NewDiagnosticApiNotFound(resource string) diag.Diagnostic {
	return diag.NewErrorDiagnostic("Buddy API error - not found", fmt.Sprintf("%s not found", resource))
}

func NewDiagnosticDecomposeError(resource string, err error) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path.Root("id"),
		"Unknown id",
		fmt.Sprintf("The provider cannot decode id of the %s:\n%s", resource, err.Error()),
	)
}

func CheckFieldEqual(field string, got string, want string) error {
	if got != want {
		return ErrorFieldFormatted(field, got, want)
	}
	return nil
}

func GenerateCertificate() (error, string) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Company, INC."},
			Country:      []string{"US"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(crand.Reader, 4096)
	if err != nil {
		return err, ""
	}
	certBytes, err := x509.CreateCertificate(crand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return err, ""
	}
	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return err, ""
	}
	return nil, certPEM.String()
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

func ErrorFieldSet(field string) error {
	return fmt.Errorf("expected %q to be empty", field)
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

func PointerInt(v int64) *int {
	p := int(v)
	return &p
}

func StringValidatorsDomain() []validator.String {
	return []validator.String{
		stringvalidator.LengthAtLeast(4),
		stringvalidator.LengthAtMost(100),
		stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9][a-z0-9\-_]+[a-z0-9]$`), "domain must be lowercase and contain only letters, numbers or dash ( - ) and footer ( _ ) characters. It must start and end with a letter or number"),
	}
}

func StringValidatorIdentifier() []validator.String {
	return []validator.String{
		stringvalidator.RegexMatches(regexp.MustCompile(`(?i)^[a-z]\w*$`), "identifier is not valid"),
	}
}

func StringValidatorsEmail() []validator.String {
	return []validator.String{
		stringvalidator.RegexMatches(regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`), "email is not valid"),
	}
}

func TestSleep(ms int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	}
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

func ArrayInt64ToInt(arr *[]int64) *[]int {
	res := make([]int, len(*arr))
	for i, v := range *arr {
		res[i] = int(v)
	}
	return &res
}

func StringSetToApi(ctx context.Context, s *types.Set) (*[]string, diag.Diagnostics) {
	var arr []string
	d := s.ElementsAs(ctx, &arr, false)
	return &arr, d
}

func Int64SetToApi(ctx context.Context, s *types.Set) (*[]int, diag.Diagnostics) {
	var arr []int64
	d := s.ElementsAs(ctx, &arr, false)
	return ArrayInt64ToInt(&arr), d
}

func GetPipelineDefinitionSource(pipeline *buddy.Pipeline) string {
	ds := pipeline.DefinitionSource
	if ds == "" {
		return buddy.PipelineDefinitionSourceLocal
	}
	return ds
}
