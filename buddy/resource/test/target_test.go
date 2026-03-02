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

var targetIgnoreImportVerify = []string{
	"auth",
	"proxy",
	"permissions",
	"project_name",
	"pipeline_id",
	"environment_id",
	"allowed_pipeline",
}

func TestAccTarget_ftp(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "33"
	username := util.RandString(10)
	password := util.RandString(10)
	secure := true
	disabled := true
	allPipelinesAllowed := true

	newName := util.RandString(10)
	newIdentifier := util.UniqueString()
	newHost := "2.2.2.2"
	newPort := "44"
	newUsername := util.RandString(10)
	newPassword := util.RandString(10)
	newSecure := false
	newDisabled := false
	newAllPipelinesAllowed := false

	typ := buddy.TargetTypeFtp

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetFtpsConfig(domain, name, identifier, host, port, username, password, true, true, allPipelinesAllowed),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:                &name,
						Identifier:          &identifier,
						Type:                &typ,
						Host:                &host,
						Port:                &port,
						Secure:              &secure,
						Disabled:            &disabled,
						AllPipelinesAllowed: &allPipelinesAllowed,
						Auth: &buddy.TargetAuth{
							Username: username,
							Password: password,
						},
					}),
				),
			},
			{
				Config: testAccTargetFtpsConfig(domain, newName, newIdentifier, newHost, newPort, newUsername, newPassword, false, false, newAllPipelinesAllowed),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:                &newName,
						Identifier:          &newIdentifier,
						Type:                &typ,
						Host:                &newHost,
						Port:                &newPort,
						Secure:              &newSecure,
						Disabled:            &newDisabled,
						AllPipelinesAllowed: &newAllPipelinesAllowed,
						Auth: &buddy.TargetAuth{
							Username: newUsername,
							Password: newPassword,
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshPassword(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	username := util.RandString(10)
	password := util.RandString(10)
	tag := util.RandString(3)
	typ := buddy.TargetTypeSsh

	newName := util.RandString(10)
	newHost := "2.2.2.2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetSshPasswordConfig(domain, name, identifier, tag, host, port, path, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Tags:       &[]string{tag},
						Host:       &host,
						Port:       &port,
						Path:       &path,
						Type:       &typ,
						Auth: &buddy.TargetAuth{
							Method:   buddy.TargetAuthMethodPassword,
							Username: username,
							Password: password,
						},
					}),
				),
			},
			{
				Config: testAccTargetSshPasswordConfig(domain, newName, identifier, tag, newHost, port, path, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &newName,
						Identifier: &identifier,
						Tags:       &[]string{tag},
						Host:       &newHost,
						Port:       &port,
						Path:       &path,
						Type:       &typ,
						Auth: &buddy.TargetAuth{
							Method:   buddy.TargetAuthMethodPassword,
							Username: username,
							Password: password,
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshProxyCredentials(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	proxyName := util.RandString(10)
	proxyUsername := util.RandString(10)
	proxyPassword := util.RandString(10)
	proxyHost := "2.2.2.2"
	proxyPort := "55"
	typ := buddy.TargetTypeSsh

	newProxyName := util.RandString(10)
	newProxyUsername := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetSshProxyConfig(domain, name, identifier, host, port, path, proxyName, proxyHost, proxyPort, proxyUsername, proxyPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Host:       &host,
						Port:       &port,
						Path:       &path,
						Type:       &typ,
						Auth: &buddy.TargetAuth{
							Method: buddy.TargetAuthMethodProxyCredentials,
						},
						Proxy: &buddy.TargetProxy{
							Host: proxyHost,
							Port: proxyPort,
							Name: proxyName,
							Auth: &buddy.TargetAuth{
								Method:   buddy.TargetAuthMethodPassword,
								Username: proxyUsername,
								Password: proxyPassword,
							},
						},
					}),
				),
			},
			{
				Config: testAccTargetSshProxyConfig(domain, name, identifier, host, port, path, newProxyName, proxyHost, proxyPort, newProxyUsername, proxyPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Host:       &host,
						Port:       &port,
						Path:       &path,
						Type:       &typ,
						Auth: &buddy.TargetAuth{
							Method: buddy.TargetAuthMethodProxyCredentials,
						},
						Proxy: &buddy.TargetProxy{
							Host: proxyHost,
							Port: proxyPort,
							Name: newProxyName,
							Auth: &buddy.TargetAuth{
								Method:   buddy.TargetAuthMethodPassword,
								Username: newProxyUsername,
								Password: proxyPassword,
							},
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshKey(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	email := util.RandEmail()
	groupName := util.RandString(10)
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	username := util.RandString(10)
	key := util.RandString(10)
	passphrase := util.RandString(10)
	typ := buddy.TargetTypeSsh
	otherLevel := buddy.TargetPermissionManage
	userLevel := buddy.TargetPermissionUseOnly
	groupLevel := buddy.TargetPermissionManage
	pipelineIdentifier := util.UniqueString()
	projectName := util.UniqueString()

	newName := util.RandString(10)
	newHost := "2.2.2.2"
	newKey := util.RandString(10)
	newOtherLevel := buddy.TargetPermissionUseOnly
	newUserLevel := buddy.TargetPermissionManage
	newGroupLevel := buddy.TargetPermissionUseOnly
	newPipelineIdentifier := util.UniqueString()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetSshKeyConfig(domain, projectName, email, groupName, name, identifier, host, port, username, key, passphrase, otherLevel, userLevel, groupLevel, pipelineIdentifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Host:       &host,
						Port:       &port,
						Type:       &typ,
						Auth: &buddy.TargetAuth{
							Method:     buddy.TargetAuthMethodSshKey,
							Username:   username,
							Key:        key,
							Passphrase: passphrase,
						},
						Permissions: &buddy.TargetPermissions{
							Others: otherLevel,
							Users: []*buddy.TargetResourcePermission{{
								AccessLevel: userLevel,
							}},
							Groups: []*buddy.TargetResourcePermission{{
								AccessLevel: groupLevel,
							}},
						},
						AllowedPipelines: &[]*buddy.TargetAllowedPipeline{
							{
								Project:  projectName,
								Pipeline: pipelineIdentifier,
							},
						},
					}),
				),
			},
			{
				Config: testAccTargetSshKeyConfig(domain, projectName, email, groupName, newName, identifier, newHost, port, username, newKey, passphrase, newOtherLevel, newUserLevel, newGroupLevel, newPipelineIdentifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &newName,
						Identifier: &identifier,
						Host:       &newHost,
						Port:       &port,
						Type:       &typ,
						Auth: &buddy.TargetAuth{
							Method:     buddy.TargetAuthMethodSshKey,
							Username:   username,
							Key:        newKey,
							Passphrase: passphrase,
						},
						Permissions: &buddy.TargetPermissions{
							Others: newOtherLevel,
							Users: []*buddy.TargetResourcePermission{{
								AccessLevel: newUserLevel,
							}},
							Groups: []*buddy.TargetResourcePermission{{
								AccessLevel: newGroupLevel,
							}},
						},
						AllowedPipelines: &[]*buddy.TargetAllowedPipeline{
							{
								Project:  projectName,
								Pipeline: newPipelineIdentifier,
							},
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshAsset(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := "/a/b/c"
	username := util.RandString(10)
	asset := util.RandString(10)
	typ := buddy.TargetTypeSsh

	newPath := "/"
	newAsset := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetSshAssetConfig(domain, name, identifier, host, port, path, username, asset),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Host:       &host,
						Port:       &port,
						Type:       &typ,
						Path:       &path,
						Auth: &buddy.TargetAuth{
							Method:   buddy.TargetAuthMethodAssetsKey,
							Username: username,
							Asset:    asset,
						},
					}),
				),
			},
			{
				Config: testAccTargetSshAssetConfig(domain, name, identifier, host, port, newPath, username, newAsset),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Host:       &host,
						Port:       &port,
						Type:       &typ,
						Path:       &newPath,
						Auth: &buddy.TargetAuth{
							Method:   buddy.TargetAuthMethodAssetsKey,
							Username: username,
							Asset:    newAsset,
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_gitHttp(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	repository := "https://aa.com"
	username := util.RandString(10)
	password := util.RandString(10)
	typ := buddy.TargetTypeGit

	newRepository := "https://bb.com"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGitHttpConfig(domain, name, identifier, repository, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
						Repository: &repository,
						Auth: &buddy.TargetAuth{
							Method:   buddy.TargetAuthMethodHttp,
							Username: username,
							Password: password,
						},
					}),
				),
			},
			{
				Config: testAccTargetGitHttpConfig(domain, name, identifier, newRepository, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
						Repository: &newRepository,
						Auth: &buddy.TargetAuth{
							Method:   buddy.TargetAuthMethodHttp,
							Username: username,
							Password: password,
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_vultr(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	newName := util.RandString(10)
	typ := buddy.TargetTypeVultr
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetVultrConfig(domain, name, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
					}),
				),
			},
			{
				Config: testAccTargetVultrConfig(domain, newName, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &newName,
						Identifier: &identifier,
						Type:       &typ,
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_upcloud(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	newName := util.RandString(10)
	typ := buddy.TargetTypeUpcloud
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetUpcloudConfig(domain, name, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
					}),
				),
			},
			{
				Config: testAccTargetUpcloudConfig(domain, newName, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &newName,
						Identifier: &identifier,
						Type:       &typ,
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_digitalOcean(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	newName := util.RandString(10)
	typ := buddy.TargetTypeDigitalOcean
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetDigitalOceanConfig(domain, name, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
					}),
				),
			},
			{
				Config: testAccTargetDigitalOceanConfig(domain, newName, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &newName,
						Identifier: &identifier,
						Type:       &typ,
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_gitSsh(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	repository := "https://aa.com"
	key := util.RandString(10)
	typ := buddy.TargetTypeGit

	newKey := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGitSshConfig(domain, name, identifier, repository, key),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
						Repository: &repository,
						Auth: &buddy.TargetAuth{
							Method: buddy.TargetAuthMethodSshKey,
							Key:    key,
						},
					}),
				),
			},
			{
				Config: testAccTargetGitSshConfig(domain, name, identifier, repository, newKey),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, &buddy.TargetOps{
						Name:       &name,
						Identifier: &identifier,
						Type:       &typ,
						Repository: &repository,
						Auth: &buddy.TargetAuth{
							Method: buddy.TargetAuthMethodSshKey,
							Key:    newKey,
						},
					}),
				),
			},
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func testAccTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_target" {
			continue
		}
		domain, targetId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		target, resp, err := acc.ApiClient.TargetService.Get(domain, targetId)
		if err == nil && target != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}

func testAccTargetGet(n string, target *buddy.Target) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain := rs.Primary.Attributes["domain"]
		targetId := rs.Primary.Attributes["target_id"]
		t, _, err := acc.ApiClient.TargetService.Get(domain, targetId)
		if err != nil {
			return err
		}
		*target = *t
		return nil
	}
}

func testAccTargetAttributes(n string, target *buddy.Target, ops *buddy.TargetOps) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsSecure, _ := strconv.ParseBool(attrs["secure"])
		attrsAllPipelinesAllowed, _ := strconv.ParseBool(attrs["all_pipelines_allowed"])
		attrsDisabled, _ := strconv.ParseBool(attrs["disabled"])
		if err := util.CheckFieldSet("Id", target.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("target_id", attrs["target_id"], target.Id); err != nil {
			return err
		}
		if err := util.CheckFieldSet("HtmlUrl", target.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], target.HtmlUrl); err != nil {
			return err
		}

		if ops.Identifier != nil {
			if err := util.CheckFieldEqual("Identifier", target.Identifier, *ops.Identifier); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("identifier", attrs["identifier"], *ops.Identifier); err != nil {
				return err
			}
		}
		if ops.Name != nil {
			if err := util.CheckFieldEqualAndSet("Name", target.Name, *ops.Name); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Name", attrs["name"], *ops.Name); err != nil {
				return err
			}
		}
		if ops.Type != nil {
			if err := util.CheckFieldEqual("Type", target.Type, *ops.Type); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("type", attrs["type"], *ops.Type); err != nil {
				return err
			}
		}
		if ops.Tags != nil {
			if err := util.CheckIntFieldEqual("Tags", len(target.Tags), len(*ops.Tags)); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Tags[0]", target.Tags[0], (*ops.Tags)[0]); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("tags.0", attrs["tags.0"], (*ops.Tags)[0]); err != nil {
				return err
			}
		}

		if ops.Scope != nil {
			if err := util.CheckFieldEqual("Scope", target.Scope, *ops.Scope); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("scope", attrs["scope"], *ops.Scope); err != nil {
				return err
			}
		}

		if ops.Host != nil {
			if err := util.CheckFieldEqual("Host", target.Host, *ops.Host); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("host", attrs["host"], *ops.Host); err != nil {
				return err
			}
		}

		if ops.Repository != nil {
			if err := util.CheckFieldEqual("Repository", target.Repository, *ops.Repository); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("repository", attrs["repository"], *ops.Repository); err != nil {
				return err
			}
		}

		if ops.Port != nil {
			if err := util.CheckFieldEqual("Port", target.Port, *ops.Port); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("port", attrs["port"], *ops.Port); err != nil {
				return err
			}
		}

		if ops.Path != nil {
			if err := util.CheckFieldEqual("Path", target.Path, *ops.Path); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("path", attrs["path"], *ops.Path); err != nil {
				return err
			}
		}

		if ops.Secure != nil {
			if err := util.CheckBoolFieldEqual("Secure", target.Secure, *ops.Secure); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("secure", attrsSecure, *ops.Secure); err != nil {
				return err
			}
		}

		if ops.AllPipelinesAllowed != nil {
			if err := util.CheckBoolFieldEqual("AllPipelinesAllowed", target.AllPipelinesAllowed, *ops.AllPipelinesAllowed); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("all_pipelines_allowed", attrsAllPipelinesAllowed, *ops.AllPipelinesAllowed); err != nil {
				return err
			}
		}

		if ops.AllowedPipelines != nil {
			if err := util.CheckIntFieldEqualAndSet("AllowedPipelines", len(target.AllowedPipelines), len(*ops.AllowedPipelines)); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("AllowedPipelines[0].Project", target.AllowedPipelines[0].Project, (*ops.AllowedPipelines)[0].Project); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("AllowedPipelines[0].Pipeline", target.AllowedPipelines[0].Pipeline, (*ops.AllowedPipelines)[0].Pipeline); err != nil {
				return err
			}
		}

		if ops.Disabled != nil {
			if err := util.CheckBoolFieldEqual("Disabled", target.Disabled, *ops.Disabled); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("disabled", attrsDisabled, *ops.Disabled); err != nil {
				return err
			}
		}

		if ops.Permissions != nil {
			if err := util.CheckFieldEqualAndSet("Permissions.Others", target.Permissions.Others, ops.Permissions.Others); err != nil {
				return err
			}
			if len(ops.Permissions.Users) > 0 {
				if err := util.CheckFieldEqualAndSet("Permissions.Users[0].AccessLevel", target.Permissions.Users[0].AccessLevel, ops.Permissions.Users[0].AccessLevel); err != nil {
					return err
				}
			} else {
				if err := util.CheckIntFieldEqual("len(Permissions.Users)", len(target.Permissions.Users), 0); err != nil {
					return err
				}
			}
			if len(ops.Permissions.Groups) > 0 {
				if err := util.CheckFieldEqualAndSet("Permissions.Groups[0].AccessLevel", target.Permissions.Groups[0].AccessLevel, ops.Permissions.Groups[0].AccessLevel); err != nil {
					return err
				}
			} else {
				if err := util.CheckIntFieldEqual("len(Permissions.Groups)", len(target.Permissions.Groups), 0); err != nil {
					return err
				}
			}
		} else {
			if err := util.CheckFieldEqualAndSet("Permissions.Others", target.Permissions.Others, buddy.TargetPermissionUseOnly); err != nil {
				return err
			}
		}
		if ops.Auth != nil {
			if ops.Auth.Method != "" {
				if err := util.CheckFieldEqualAndSet("Auth.Method", target.Auth.Method, ops.Auth.Method); err != nil {
					return err
				}
			}
			if ops.Auth.Username != "" {
				if err := util.CheckFieldEqualAndSet("Auth.Username", target.Auth.Username, ops.Auth.Username); err != nil {
					return err
				}
			}
			if ops.Auth.Password != "" {
				if err := util.CheckFieldSet("Auth.Password", target.Auth.Password); err != nil {
					return err
				}
			}
			if ops.Auth.Asset != "" {
				if err := util.CheckFieldEqualAndSet("Auth.Asset", target.Auth.Asset, ops.Auth.Asset); err != nil {
					return err
				}
			}
			if ops.Auth.Passphrase != "" {
				if err := util.CheckFieldSet("Auth.Passphrase", target.Auth.Passphrase); err != nil {
					return err
				}
			}
			if ops.Auth.Key != "" {
				if err := util.CheckFieldSet("Auth.Key", target.Auth.Key); err != nil {
					return err
				}
			}
			if ops.Auth.KeyPath != "" {
				if err := util.CheckFieldEqualAndSet("Auth.KeyPath", target.Auth.KeyPath, ops.Auth.KeyPath); err != nil {
					return err
				}
			}
		}
		if ops.Proxy != nil {
			if ops.Proxy.Name != "" {
				if err := util.CheckFieldEqualAndSet("Proxy.Name", target.Proxy.Name, ops.Proxy.Name); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func testAccTargetFtpsConfig(domain string, name string, identifier string, host string, port string, username string, password string, secure bool, disabled bool, allPipelinesAllowed bool) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "FTP"
    host       = "%s"
    port       = "%s"
    secure     = %t
    disabled   = %t
    auth {
        username = "%s"
        password = "%s"
    }
		all_pipelines_allowed = %t
}`, domain, name, identifier, host, port, secure, disabled, username, password, allPipelinesAllowed)
}

func testAccTargetSshPasswordConfig(domain string, name string, identifier string, tag string, host string, port string, path string, username string, password string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "abcdef" 
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "SSH"
    tags         = ["%s"]
    host         = "%s"
    port         = "%s"
    path         = "%s"
    auth {
        method   = "PASSWORD"
        username = "%s"
        password = "%s"
    }
}`, domain, name, identifier, tag, host, port, path, username, password)
}

func testAccTargetGitHttpConfig(domain string, name string, identifier string, repository string, username string, password string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "abcdef" 
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "GIT"
    repository   = "%s"
    auth {
        method   = "HTTP"
        username = "%s"
        password = "%s"
    }
}`, domain, name, identifier, repository, username, password)
}

func testAccTargetGitSshConfig(domain string, name string, identifier string, repository string, key string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "abcdef" 
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "GIT"
    repository   = "%s"
    auth {
        method   = "SSH_KEY"
        key = "%s"
    }
}`, domain, name, identifier, repository, key)
}

func testAccTargetVultrConfig(domain string, name string, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_integration" "test" {
    domain = buddy_workspace.test.domain
    type   = "VULTR"
    name   = "test"
    token  = "abcdef"
    scope  = "WORKSPACE"
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain 
    name         = "%s"
    identifier   = "%s"
    type         = "VULTR"
    host         = "1.1.1.1"
    port         = "22"
    integration  = "${buddy_integration.test.identifier}"
    auth {
        method   = "PASSWORD"
        username = "user"
        password = "pass"
    }
}`, domain, name, identifier)
}

func testAccTargetUpcloudConfig(domain string, name string, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_integration" "test" {
    domain    = buddy_workspace.test.domain
    type      = "UPCLOUD"
    name      = "test"
    username  = "user"
    password  = "pass"
    scope     = "WORKSPACE"
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain 
    name         = "%s"
    identifier   = "%s"
    type         = "UPCLOUD"
    host         = "1.1.1.1"
    port         = "22"
    integration  = "${buddy_integration.test.identifier}"
    auth {
        method   = "PASSWORD"
        username = "user"
        password = "pass"
    }
}`, domain, name, identifier)
}

func testAccTargetDigitalOceanConfig(domain string, name string, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_integration" "test" {
    domain = buddy_workspace.test.domain
    type   = "DIGITAL_OCEAN"
    name   = "test"
    token  = "abcdef"
    scope  = "WORKSPACE"
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain 
    name         = "%s"
    identifier   = "%s"
    type         = "DIGITAL_OCEAN"
    host         = "1.1.1.1"
    port         = "22"
    integration  = "${buddy_integration.test.identifier}"
    auth {
        method   = "PASSWORD"
        username = "user"
        password = "pass"
    }
}`, domain, name, identifier)
}

func testAccTargetSshAssetConfig(domain string, name string, identifier string, host string, port string, path string, username string, asset string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "abcdef" 
}

resource "buddy_environment" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "dev"
    identifier   = "dev"
}

resource "buddy_target" "test" {
    domain         = buddy_workspace.test.domain
    project_name   = buddy_project.test.name
    environment_id = buddy_environment.test.environment_id
    name           = "%s"
    identifier     = "%s"
    type           = "SSH" 
    host           = "%s"
    port           = "%s"
    path           = "%s"
    auth {
        method   = "ASSETS_KEY"
        username = "%s"
        asset = "%s"
    }
}`, domain, name, identifier, host, port, path, username, asset)
}

func testAccTargetSshProxyConfig(domain string, name string, identifier string, host string, port string, path string, proxyName string, proxyHost string, proxyPort string, proxyUser string, proxyPass string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain         = buddy_workspace.test.domain 
    name           = "%s"
    identifier     = "%s"
    type           = "SSH" 
    host           = "%s"
    port           = "%s"
    path           = "%s"
    auth {
        method   = "PROXY_CREDENTIALS"
    }
    proxy {
        name = "%s"
        host = "%s"
        port = "%s"
        auth {
           method = "PASSWORD"
           username = "%s"
           password = "%s"
        }
    }
}`, domain, name, identifier, host, port, path, proxyName, proxyHost, proxyPort, proxyUser, proxyPass)
}

func testAccTargetSshKeyConfig(domain string, projectName string, email string, groupName string, name string, identifier string, host string, port string, username string, key string, passphrase string, othersLevel string, userLevel string, groupLevel string, pipelineIdentifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "%s"
}

resource "buddy_pipeline" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
		identifier   = "%s"
}

resource "buddy_member" "test" {
    domain = "${buddy_workspace.test.domain}"
    email = "%s"
}

resource "buddy_group" "test" {
	domain = "${buddy_workspace.test.domain}"
	name = "%s"
}

resource "buddy_permission" "test" {
    domain = "${buddy_workspace.test.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	  sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "test" {
	domain = "${buddy_workspace.test.domain}"
	project_name = "${buddy_project.test.name}"
	member_id = "${buddy_member.test.member_id}"
	permission_id = "${buddy_permission.test.permission_id}"
}

resource "buddy_project_group" "test" {
	domain = "${buddy_workspace.test.domain}"
	project_name = "${buddy_project.test.name}"
	group_id = "${buddy_group.test.group_id}"
	permission_id = "${buddy_permission.test.permission_id}"
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    pipeline_id  = buddy_pipeline.test.pipeline_id
    name         = "%s"
    identifier   = "%s"
    type         = "SSH" 
    host         = "%s"
    port         = "%s"
    auth {
        method   = "SSH_KEY"
        username = "%s"
        key      = "%s"
        passphrase = "%s"
    }
    permissions {
        others = "%s"
        user {
           id = "${buddy_project_member.test.member_id}"
           access_level = "%s"
        }
        group {
           id = "${buddy_project_group.test.group_id}"
           access_level = "%s"
        }
    }
		allowed_pipeline {
      project = buddy_project.test.name
      pipeline = "%s"
		}
}`, domain, projectName, pipelineIdentifier, pipelineIdentifier, email, groupName, name, identifier, host, port, username, key, passphrase, othersLevel, userLevel, groupLevel, pipelineIdentifier)
}
