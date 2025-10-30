package test

import (
	"encoding/base64"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"log"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
	"time"
)

type testAccPipelineExpectedAttributes struct {
	Name                      string
	Identifier                string
	AlwaysFromScratch         bool
	DescriptionRequired       bool
	ConcurrentPipelineRuns    bool
	GitChangesetBase          string
	FilesystemChangesetBase   string
	FailOnPrepareEnvWarning   bool
	FetchAllRefs              bool
	AutoClearCache            bool
	NoSkipToMostRecent        bool
	DoNotCreateCommitStatus   bool
	CloneDepth                int
	Paused                    bool
	PauseOnRepeatedFailures   int
	Priority                  string
	IgnoreFailOnProjectStatus bool
	ExecutionMessageTemplate  string
	TargetSiteUrl             string
	ManageVariablesByYaml     bool
	ManagePermissionsByYaml   bool
	Disabled                  bool
	DisablingReason           string
	Cpu                       string
	Creator                   *buddy.Profile
	Project                   *buddy.Project
	Ref                       string
	Loop                      string
	Event                     *buddy.PipelineEvent
	TriggerConditions         []*buddy.PipelineTriggerCondition
	Permissions               *buddy.PipelinePermissions
	PermissionUser            *buddy.Member
	PermissionGroup           *buddy.Group
}

type testAccPipelineRemoteExpectedAttributes struct {
	Name              string
	GitConfigRef      string
	GitConfig         *buddy.PipelineGitConfig
	DefinitionSource  string
	RemoteProjectName string
	RemoteRef         string
	RemoteBranch      string
	RemotePath        string
	RemoteParam       string
	Creator           *buddy.Profile
	Project           *buddy.Project
}

func TestAccPipeline_permissions(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	var member buddy.Member
	var group buddy.Group
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	ref := util.RandString(10)
	email := util.RandEmail()
	groupName := util.RandString(10)
	othersPerm1 := buddy.PipelinePermissionDenied
	userPerm1 := buddy.PipelinePermissionRunOnly
	othersPerm2 := buddy.PipelinePermissionRunOnly
	userPerm2 := buddy.PipelinePermissionDenied
	groupPerm1 := buddy.PipelinePermissionReadWrite
	othersPerm3 := buddy.PipelinePermissionReadWrite
	groupPerm2 := buddy.PipelinePermissionDefault
	loop := util.UniqueString()
	newLoop := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelinePermissionsEmpty(domain, projectName, name, ref),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						AlwaysFromScratch:         false,
						AutoClearCache:            false,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   false,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             "",
						ExecutionMessageTemplate:  "",
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
						Permissions:               nil,
					}),
				),
			},
			// update pipeline permission to user
			{
				Config: testAccPipelinePermissionsUser(domain, projectName, name, ref, loop, email, othersPerm1, userPerm1),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccMemberGet("buddy_member.a", &member),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						AlwaysFromScratch:         false,
						AutoClearCache:            false,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   false,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             "",
						ExecutionMessageTemplate:  "",
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
						Loop:                      loop,
						Permissions: &buddy.PipelinePermissions{
							Others: othersPerm1,
							Users: []*buddy.PipelineResourcePermission{
								{
									AccessLevel: userPerm1,
								},
							},
						},
						PermissionUser: &member,
					}),
				),
			},
			// update pipeline permission to group & user
			{
				Config: testAccPipelinePermissionsUserGroup(domain, projectName, name, ref, newLoop, email, groupName, othersPerm2, userPerm2, groupPerm1),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccMemberGet("buddy_member.a", &member),
					testAccGroupGet("buddy_group.g", &group),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						AlwaysFromScratch:         false,
						AutoClearCache:            false,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   false,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             "",
						ExecutionMessageTemplate:  "",
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
						Loop:                      newLoop,
						Permissions: &buddy.PipelinePermissions{
							Others: othersPerm2,
							Users: []*buddy.PipelineResourcePermission{
								{
									AccessLevel: userPerm2,
								},
							},
							Groups: []*buddy.PipelineResourcePermission{
								{
									AccessLevel: groupPerm1,
								},
							},
						},
						PermissionUser:  &member,
						PermissionGroup: &group,
					}),
				),
			},
			// update pipeline permission to group
			{
				Config: testAccPipelinePermissionsGroup(domain, projectName, name, ref, email, groupName, othersPerm3, groupPerm2),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccGroupGet("buddy_group.g", &group),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						AlwaysFromScratch:         false,
						AutoClearCache:            false,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   false,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             "",
						ExecutionMessageTemplate:  "",
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
						Permissions: &buddy.PipelinePermissions{
							Others: othersPerm3,
							Groups: []*buddy.PipelineResourcePermission{
								{
									AccessLevel: groupPerm2,
								},
							},
						},
						PermissionGroup: &group,
					}),
				),
			},
			// to empty
			{
				Config: testAccPipelinePermissionsBackToEmpty(domain, projectName, name, ref, email, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						AlwaysFromScratch:         false,
						AutoClearCache:            false,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   false,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             "",
						ExecutionMessageTemplate:  "",
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
						Permissions:               nil,
					}),
				),
			},
			// import pipeline
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permissions"},
			},
		},
	})
}

func TestAccPipeline_remoteBranch(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	remoteProjectName := util.UniqueString()
	remoteProjectName2 := util.UniqueString()
	gitConfigBranch := util.UniqueString()
	gitConfigPath := util.UniqueString()
	gitConfigYml := fmt.Sprintf(`
  git_config = {
    project = "%s"
    branch = "%s"
    path="%s"
  }
`, projectName, gitConfigBranch, gitConfigPath)
	name := util.RandString(10)
	remoteBranch := "master"
	remotePath := util.RandString(10)
	remotePath2 := util.RandString(10)
	yaml := testAccPipelineGetRemoteYaml(remoteBranch)
	cmd := "ls"
	cmd2 := "pwd"
	branch := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create workspace & projects
			{
				Config: testAccPipelineConfigRemoteInit(domain, projectName, remoteProjectName, remoteProjectName2),
			},
			// create yaml file & pipeline
			{
				PreConfig: func() {
					testAccPipelineCreateRemoteYaml(domain, remoteProjectName, remotePath, yaml)
					testAccPipelineCreateRemoteYaml(domain, remoteProjectName2, remotePath2, yaml)
				},
				Config: testAccPipelineConfigRemoteBranchMain(domain, projectName, name, buddy.PipelineGitConfigRefNone, "", remoteProjectName, remoteProjectName2, remoteProjectName, remoteBranch, remotePath, cmd),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:              name,
						GitConfigRef:      buddy.PipelineGitConfigRefNone,
						RemoteParam:       cmd,
						RemoteProjectName: remoteProjectName,
						RemoteBranch:      remoteBranch,
						RemotePath:        remotePath,
						DefinitionSource:  buddy.PipelineDefinitionSourceRemote,
						Creator:           &profile,
						Project:           &project,
					}),
				),
			},
			// update remote project
			{
				Config: testAccPipelineConfigRemoteBranchMain(domain, projectName, name, buddy.PipelineGitConfigRefDynamic, "", remoteProjectName, remoteProjectName2, remoteProjectName2, remoteBranch, remotePath2, cmd2),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:              name,
						GitConfigRef:      buddy.PipelineGitConfigRefDynamic,
						RemoteParam:       cmd2,
						RemoteProjectName: remoteProjectName2,
						RemoteBranch:      remoteBranch,
						RemotePath:        remotePath2,
						DefinitionSource:  buddy.PipelineDefinitionSourceRemote,
						Creator:           &profile,
						Project:           &project,
					}),
				),
			},
			// update to git config fixed
			{
				Config: testAccPipelineConfigRemoteBranchMain(domain, projectName, name, buddy.PipelineGitConfigRefFixed, gitConfigYml, remoteProjectName, remoteProjectName2, remoteProjectName2, remoteBranch, remotePath2, cmd2),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:         name,
						GitConfigRef: buddy.PipelineGitConfigRefFixed,
						GitConfig: &buddy.PipelineGitConfig{
							Project: projectName,
							Branch:  gitConfigBranch,
							Path:    gitConfigPath,
						},
						RemoteParam:       cmd2,
						RemoteProjectName: remoteProjectName2,
						RemoteBranch:      remoteBranch,
						RemotePath:        remotePath2,
						DefinitionSource:  buddy.PipelineDefinitionSourceRemote,
						Creator:           &profile,
						Project:           &project,
					}),
				),
			},
			// change to local
			{
				Config: testAccPipelineConfigRemoteToLocal(domain, projectName, remoteProjectName, remoteProjectName2, name, branch),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:             name,
						GitConfigRef:     buddy.PipelineGitConfigRefNone,
						DefinitionSource: buddy.PipelineDefinitionSourceLocal,
						Creator:          &profile,
						Project:          &project,
					}),
				),
			},
		},
	})
}

func TestAccPipeline_remote(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	remoteProjectName := util.UniqueString()
	remoteProjectName2 := util.UniqueString()
	gitConfigBranch := util.UniqueString()
	gitConfigPath := util.UniqueString()
	gitConfigYml := fmt.Sprintf(`
  git_config = {
    project = "%s"
    branch = "%s"
    path="%s"
  }
`, projectName, gitConfigBranch, gitConfigPath)
	name := util.RandString(10)
	remoteBranch := "master"
	remotePath := util.RandString(10)
	remotePath2 := util.RandString(10)
	yaml := testAccPipelineGetRemoteYaml(remoteBranch)
	cmd := "ls"
	cmd2 := "pwd"
	branch := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create workspace & projects
			{
				Config: testAccPipelineConfigRemoteInit(domain, projectName, remoteProjectName, remoteProjectName2),
			},
			// create yaml file & pipeline
			{
				PreConfig: func() {
					testAccPipelineCreateRemoteYaml(domain, remoteProjectName, remotePath, yaml)
					testAccPipelineCreateRemoteYaml(domain, remoteProjectName2, remotePath2, yaml)
				},
				Config: testAccPipelineConfigRemoteMain(domain, projectName, name, buddy.PipelineGitConfigRefNone, "", remoteProjectName, remoteProjectName2, remoteProjectName, remoteBranch, remotePath, cmd),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:              name,
						GitConfigRef:      buddy.PipelineGitConfigRefNone,
						RemoteParam:       cmd,
						RemoteProjectName: remoteProjectName,
						RemoteRef:         remoteBranch,
						RemotePath:        remotePath,
						DefinitionSource:  buddy.PipelineDefinitionSourceRemote,
						Creator:           &profile,
						Project:           &project,
					}),
				),
			},
			// update remote project
			{
				Config: testAccPipelineConfigRemoteMain(domain, projectName, name, buddy.PipelineGitConfigRefDynamic, "", remoteProjectName, remoteProjectName2, remoteProjectName2, remoteBranch, remotePath2, cmd2),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:              name,
						GitConfigRef:      buddy.PipelineGitConfigRefDynamic,
						RemoteParam:       cmd2,
						RemoteProjectName: remoteProjectName2,
						RemoteRef:         remoteBranch,
						RemotePath:        remotePath2,
						DefinitionSource:  buddy.PipelineDefinitionSourceRemote,
						Creator:           &profile,
						Project:           &project,
					}),
				),
			},
			// update to git config fixed
			{
				Config: testAccPipelineConfigRemoteMain(domain, projectName, name, buddy.PipelineGitConfigRefFixed, gitConfigYml, remoteProjectName, remoteProjectName2, remoteProjectName2, remoteBranch, remotePath2, cmd2),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:         name,
						GitConfigRef: buddy.PipelineGitConfigRefFixed,
						GitConfig: &buddy.PipelineGitConfig{
							Project: projectName,
							Branch:  gitConfigBranch,
							Path:    gitConfigPath,
						},
						RemoteParam:       cmd2,
						RemoteProjectName: remoteProjectName2,
						RemoteRef:         remoteBranch,
						RemotePath:        remotePath2,
						DefinitionSource:  buddy.PipelineDefinitionSourceRemote,
						Creator:           &profile,
						Project:           &project,
					}),
				),
			},
			// change to local
			{
				Config: testAccPipelineConfigRemoteToLocal(domain, projectName, remoteProjectName, remoteProjectName2, name, branch),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.remote", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineRemoteAttributes("buddy_pipeline.remote", &pipeline, &testAccPipelineRemoteExpectedAttributes{
						Name:             name,
						GitConfigRef:     buddy.PipelineGitConfigRefNone,
						DefinitionSource: buddy.PipelineDefinitionSourceLocal,
						Creator:          &profile,
						Project:          &project,
					}),
				),
			},
		},
	})
}

func TestAccPipeline_schedule(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	eventType := buddy.PipelineEventTypeSchedule
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	startDate := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	newStartDate := time.Now().UTC().Add(5 * time.Hour).Format(time.RFC3339)
	priority := buddy.PipelinePriorityLow
	newPriority := buddy.PipelinePriorityHigh
	pausedFailures := 78
	delay := 5
	newDelay := 7
	newPausedFailures := 1
	gitChangeSet := buddy.PipelineGitChangeSetBasePullRequest
	newGitChangeSet := buddy.PipelineGitChangeSetBaseLatestRunMatchingRef
	filesystemChangeSet := buddy.PipelineFilesystemChangeSetBaseDateModified
	newFilesystemChangeSet := buddy.PipelineFilesystemChangeSetBaseContents
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigSchedule(domain, projectName, name, startDate, delay, true, priority, pausedFailures, true, true, gitChangeSet, filesystemChangeSet),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:     name,
						Project:  &project,
						Creator:  &profile,
						Priority: priority,
						Event: &buddy.PipelineEvent{
							Type:      eventType,
							StartDate: startDate,
							Delay:     delay,
						},
						GitChangesetBase:        gitChangeSet,
						FilesystemChangesetBase: filesystemChangeSet,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  true,
						PauseOnRepeatedFailures: pausedFailures,
						DescriptionRequired:     true,
						ConcurrentPipelineRuns:  true,
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigSchedule(domain, projectName, newName, newStartDate, newDelay, false, newPriority, newPausedFailures, false, false, newGitChangeSet, newFilesystemChangeSet),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:    newName,
						Project: &project,
						Creator: &profile,
						Event: &buddy.PipelineEvent{
							Type:      eventType,
							StartDate: newStartDate,
							Delay:     newDelay,
						},
						Priority:                newPriority,
						GitChangesetBase:        newGitChangeSet,
						FilesystemChangesetBase: newFilesystemChangeSet,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  false,
						DescriptionRequired:     false,
						ConcurrentPipelineRuns:  false,
						PauseOnRepeatedFailures: newPausedFailures,
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"event"},
			},
		},
	})
}

func TestAccPipeline_schedule_cron(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	eventType := buddy.PipelineEventTypeSchedule
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	cron := "15 14 1 * *"
	newCron := "0 22 * * 1-5"
	newTimezone := "Europer/Warsaw"
	reason := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigScheduleCron(domain, projectName, name, cron, "", true),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:    name,
						Project: &project,
						Creator: &profile,
						Event: &buddy.PipelineEvent{
							Type: eventType,
							Cron: cron,
						},
						Priority:                buddy.PipelinePriorityNormal,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  true,
						PauseOnRepeatedFailures: 100,
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigScheduleCron(domain, projectName, newName, newCron, newTimezone, false),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:    newName,
						Project: &project,
						Creator: &profile,
						Event: &buddy.PipelineEvent{
							Type:     eventType,
							Cron:     newCron,
							Timezone: newTimezone,
						},
						Priority:                buddy.PipelinePriorityNormal,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  false,
						PauseOnRepeatedFailures: 100,
					}),
				),
			},
			// disable
			{
				Config: testAccPipelineConfigScheduleCronDisabled(domain, projectName, newName, newCron, false, true, reason),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:    newName,
						Project: &project,
						Creator: &profile,
						Event: &buddy.PipelineEvent{
							Type: eventType,
							Cron: newCron,
						},
						Priority:                buddy.PipelinePriorityNormal,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  false,
						Disabled:                true,
						DisablingReason:         reason,
					}),
				),
			},
			// enable
			{
				Config: testAccPipelineConfigScheduleCronDisabled(domain, projectName, newName, newCron, false, false, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:    newName,
						Project: &project,
						Creator: &profile,
						Event: &buddy.PipelineEvent{
							Type: eventType,
							Cron: newCron,
						},
						Priority:                buddy.PipelinePriorityNormal,
						FailOnPrepareEnvWarning: true,
						FetchAllRefs:            true,
						Paused:                  false,
						Disabled:                false,
						DisablingReason:         "",
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"event"},
			},
		},
	})
}

func testAccPipelineConfigScheduleCron(domain string, projectName string, name string, cron string, timezone string, paused bool) string {
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
  event {
    type = "SCHEDULE"
    cron = "%s"
    timezone = "%s"
  }	
	paused = %t
	fail_on_prepare_env_warning = true
	fetch_all_refs = true
}
`, domain, projectName, name, cron, timezone, paused)
}

func testAccPipelineConfigScheduleCronDisabled(domain string, projectName string, name string, cron string, paused bool, disabled bool, reason string) string {
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
  event {
    type = "SCHEDULE"
    cron = "%s"
  }	
	paused = %t
	fail_on_prepare_env_warning = true
	fetch_all_refs = true
  disabled = %t
	disabling_reason = "%s"	
}
`, domain, projectName, name, cron, paused, disabled, reason)
}

func testAccPipelineGetRemoteYaml(branch string) string {
	return fmt.Sprintf(`
- pipeline: "test"
  refs:
  - "refs/heads/%s"
  actions:
  - action: "Execute: ls"
    type: "BUILD"
    working_directory: "/buddy/test"
    docker_image_name: "library/ubuntu"
    docker_image_tag: "18.04"
    execute_commands:
    - "!{cmd}"
    volume_mappings:
    - "/:/buddy/test"
    cache_base_image: true
    shell: "BASH"
`, branch)
}

func testAccPipelineCreateRemoteYaml(domain string, projectName string, path string, yaml string) {
	content := base64.StdEncoding.EncodeToString([]byte(yaml))
	message := "test"
	_, _, err := acc.ApiClient.SourceService.CreateFile(domain, projectName, &buddy.SourceFileOps{
		Content: &content,
		Path:    &path,
		Message: &message,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func testAccPipelineConfigRemoteToLocal(domain string, projectName string, remoteProjectName string, remoteProjectName2 string, name string, branch string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj2" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "remote" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	name = "%s"
	definition_source = "LOCAL"
  refs = ["refs/heads/%s"]
}
`, domain, projectName, remoteProjectName, remoteProjectName2, name, branch)
}

func testAccPipelineConfigRemoteMain(domain string, projectName string, name string, gitConfigRef string, gitConfig string, remoteProjectName string, remoteProjectName2 string, selectedRemoteProjectName string, remoteBranch string, remotePath string, cmd string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj2" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "remote" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	depends_on = [buddy_project.remote_proj, buddy_project.remote_proj2]
	name = "%s"
  git_config_ref = "%s"
  %s
	definition_source = "REMOTE"
	remote_project_name = "%s"
	remote_ref = "%s"
	remote_path = "%s"
	remote_parameter {
		key = "cmd"
		value = "%s"
	}
}
`, domain, projectName, remoteProjectName, remoteProjectName2, name, gitConfigRef, gitConfig, selectedRemoteProjectName, remoteBranch, remotePath, cmd)
}

func testAccPipelineConfigRemoteBranchMain(domain string, projectName string, name string, gitConfigRef string, gitConfig string, remoteProjectName string, remoteProjectName2 string, selectedRemoteProjectName string, remoteBranch string, remotePath string, cmd string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj2" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_pipeline" "remote" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	depends_on = [buddy_project.remote_proj, buddy_project.remote_proj2]
	name = "%s"
  git_config_ref = "%s"
  %s
	definition_source = "REMOTE"
	remote_project_name = "%s"
	remote_branch = "%s"
	remote_path = "%s"
	remote_parameter {
		key = "cmd"
		value = "%s"
	}
}
`, domain, projectName, remoteProjectName, remoteProjectName2, name, gitConfigRef, gitConfig, selectedRemoteProjectName, remoteBranch, remotePath, cmd)
}

func testAccPipelineConfigRemoteInit(domain string, projectName string, remoteProjectName string, remoteProjectName2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "remote_proj2" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}
`, domain, projectName, remoteProjectName, remoteProjectName2)
}

func testAccPipelineConfigSchedule(domain string, projectName string, name string, startDate string, delay int, paused bool, priority string, pausedFailures int, descRequired bool, concurrentRuns bool, gitChangesetBase string, filesystemChangesetBase string) string {
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
    event {
      type = "SCHEDULE"
      start_date = "%s"
	    delay = %d
    }	
	  paused = %t
	  priority = "%s"
    pause_on_repeated_failures = %d
    description_required = %t
    concurrent_pipeline_runs = %t
    git_changeset_base = "%s"
    filesystem_changeset_base = "%s"
	  fail_on_prepare_env_warning = true
	  fetch_all_refs = true
}
`, domain, projectName, name, startDate, delay, paused, priority, pausedFailures, descRequired, concurrentRuns, gitChangesetBase, filesystemChangesetBase)
}

func TestAccPipeline_event_email(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	eventType := buddy.PipelineEventTypeEmail
	prefix := util.RandString(10)
	newPrefix := util.RandString(10)
	whitelist := util.RandEmail()
	newWhitelist := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigEventEmail(domain, projectName, name, prefix, whitelist),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type:   eventType,
							Prefix: prefix,
							Whitelist: []string{
								whitelist,
							},
						},
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigEventEmail(domain, projectName, newName, newPrefix, newWhitelist),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type:   eventType,
							Prefix: newPrefix,
							Whitelist: []string{
								newWhitelist,
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
				ImportStateVerifyIgnore: []string{"event"},
			},
		},
	})
}

func TestAccPipeline_event_webhook(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	eventType := buddy.PipelineEventTypeWebhook
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigEventWebhook(domain, projectName, name, true),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type: eventType,
							Totp: true,
						},
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigEventWebhook(domain, projectName, newName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type: eventType,
							Totp: false,
						},
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"event"},
			},
		},
	})
}

func TestAccPipeline_event_pull_request(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	branch := util.RandString(10)
	newBranch := util.RandString(10)
	prEvent := buddy.PipelinePullRequestEventClosed
	newPrEvent := buddy.PipelinePullRequestEventEdited
	eventType := buddy.PipelineEventTypePullRequest
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigEventPullRequest(domain, projectName, name, eventType, branch, prEvent),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type:     eventType,
							Branches: []string{branch},
							Events:   []string{prEvent},
						},
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigEventPullRequest(domain, projectName, newName, eventType, newBranch, newPrEvent),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type:     eventType,
							Branches: []string{newBranch},
							Events:   []string{newPrEvent},
						},
					}),
				),
			},
			// import
			{
				ResourceName:            "buddy_pipeline.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"event"},
			},
		},
	})
}

func TestAccPipeline_event(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	identifier := util.UniqueString()
	newIdentifier := util.UniqueString()
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
	eventType := buddy.PipelineEventTypePush
	newEventType := buddy.PipelineEventTypeCreateRef
	user := util.RandEmail()
	newUser := util.RandEmail()
	group := util.RandString(10)
	newGroup := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigEvent(domain, projectName, name, identifier, eventType, ref, tcChangePath, tcVarKey, tcVarValue, tcHours, tcDays, tcZoneId, user, group),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    name,
						Identifier:              identifier,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type: eventType,
							Refs: []string{ref},
						},
						TriggerConditions: []*buddy.PipelineTriggerCondition{
							{
								TriggerCondition: buddy.PipelineTriggerConditionOnChange,
							},
							{
								TriggerCondition:      buddy.PipelineTriggerConditionOnChangeAtPath,
								TriggerConditionPaths: []string{tcChangePath},
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarIs,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarIsNot,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarContains,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarNotContains,
								TriggerVariableKey:   tcVarKey,
								TriggerVariableValue: tcVarValue,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionDateTime,
								TriggerHours:     []int{tcHours},
								TriggerDays:      []int{tcDays},
								Timezone:         tcZoneId,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIs,
								TriggerUser:      user,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIsNot,
								TriggerUser:      user,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIsInGroup,
								TriggerGroup:     group,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIsNotInGroup,
								TriggerGroup:     group,
							},
						},
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigEvent(domain, projectName, newName, newIdentifier, newEventType, newRef, newTcChangePath, newTcVarKey, newTcVarValue, newTcHours, newTcDays, newTcZoneId, newUser, newGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                    newName,
						Identifier:              newIdentifier,
						Project:                 &project,
						Creator:                 &profile,
						FailOnPrepareEnvWarning: false,
						FetchAllRefs:            false,
						Priority:                buddy.PipelinePriorityNormal,
						Event: &buddy.PipelineEvent{
							Type: newEventType,
							Refs: []string{newRef},
						},
						TriggerConditions: []*buddy.PipelineTriggerCondition{
							{
								TriggerCondition: buddy.PipelineTriggerConditionOnChange,
							},
							{
								TriggerCondition:      buddy.PipelineTriggerConditionOnChangeAtPath,
								TriggerConditionPaths: []string{newTcChangePath},
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarIs,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarIsNot,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarContains,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition:     buddy.PipelineTriggerConditionVarNotContains,
								TriggerVariableKey:   newTcVarKey,
								TriggerVariableValue: newTcVarValue,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionDateTime,
								TriggerHours:     []int{newTcHours},
								TriggerDays:      []int{newTcDays},
								Timezone:         newTcZoneId,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIs,
								TriggerUser:      newUser,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIsNot,
								TriggerUser:      newUser,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIsInGroup,
								TriggerGroup:     newGroup,
							},
							{
								TriggerCondition: buddy.PipelineTriggerConditionTriggeringUserIsNotInGroup,
								TriggerGroup:     newGroup,
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
				ImportStateVerifyIgnore: []string{"trigger_condition", "event"},
			},
		},
	})
}

func TestAccPipeline_click(t *testing.T) {
	var pipeline buddy.Pipeline
	var project buddy.Project
	var profile buddy.Profile
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	ref := util.RandString(10)
	newRef := util.RandString(10)
	msgTemplate := util.RandString(10)
	newMsgTemplate := util.RandString(10)
	cpu := buddy.PipelineCpuArm
	newCpu := buddy.PipelineCpuX64
	targetUrl := "https://127.0.0.1"
	newTargetUrl := "https://192.168.1.1"
	cloneDepth := 1
	newCloneDepth := 5
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPipelineCheckDestroy,
		Steps: []resource.TestStep{
			// create pipeline
			{
				Config: testAccPipelineConfigClick(domain, projectName, name, true, false, true, false, true, false, true, msgTemplate, targetUrl, ref, cloneDepth, cpu, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      name,
						Cpu:                       cpu,
						AlwaysFromScratch:         true,
						AutoClearCache:            false,
						NoSkipToMostRecent:        true,
						DoNotCreateCommitStatus:   false,
						IgnoreFailOnProjectStatus: true,
						FetchAllRefs:              true,
						FailOnPrepareEnvWarning:   false,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             targetUrl,
						ExecutionMessageTemplate:  msgTemplate,
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       ref,
						CloneDepth:                cloneDepth,
						ManagePermissionsByYaml:   true,
						ManageVariablesByYaml:     false,
					}),
				),
			},
			// update pipeline
			{
				Config: testAccPipelineConfigClick(domain, projectName, newName, false, true, false, true, false, true, false, newMsgTemplate, newTargetUrl, newRef, newCloneDepth, newCpu, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccPipelineGet("buddy_pipeline.bar", &pipeline),
					testAccProjectGet("buddy_project.proj", &project),
					testAccProfileGet(&profile),
					testAccPipelineAttributes("buddy_pipeline.bar", &pipeline, &testAccPipelineExpectedAttributes{
						Name:                      newName,
						Cpu:                       newCpu,
						AlwaysFromScratch:         false,
						AutoClearCache:            true,
						NoSkipToMostRecent:        false,
						DoNotCreateCommitStatus:   true,
						IgnoreFailOnProjectStatus: false,
						FetchAllRefs:              false,
						FailOnPrepareEnvWarning:   true,
						Priority:                  buddy.PipelinePriorityNormal,
						TargetSiteUrl:             newTargetUrl,
						ExecutionMessageTemplate:  newMsgTemplate,
						Project:                   &project,
						Creator:                   &profile,
						Ref:                       newRef,
						CloneDepth:                newCloneDepth,
						ManagePermissionsByYaml:   false,
						ManageVariablesByYaml:     true,
					}),
				),
			},
			// import pipeline
			{
				ResourceName:      "buddy_pipeline.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPipelineGitConfig(pipeline *buddy.Pipeline, attrs map[string]string, gitConfigRef string, gitConfig *buddy.PipelineGitConfig) error {
	if gitConfigRef != "" {
		if err := util.CheckFieldEqualAndSet("GitConfigRef", pipeline.GitConfigRef, gitConfigRef); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("git_config_ref", attrs["git_config_ref"], gitConfigRef); err != nil {
			return err
		}
		attrsGitConfigProject, attrsGitConfigProjectExists := attrs["git_config.project"]
		attrsGitConfigBranch, attrsGitConfigBranchExists := attrs["git_config.branch"]
		attrsGitConfigPath, attrsGitConfigPathExists := attrs["git_config.path"]
		if gitConfig == nil {
			if pipeline.GitConfig != nil {
				return util.ErrorFieldSet("GitConfig")
			}
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
			if pipeline.GitConfig == nil {
				return util.ErrorFieldEmpty("GitConfig")
			}
			if err := util.CheckFieldEqualAndSet("GitConfig.Project", pipeline.GitConfig.Project, gitConfig.Project); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("GitConfig.Branch", pipeline.GitConfig.Branch, gitConfig.Branch); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("GitConfig.Path", pipeline.GitConfig.Path, gitConfig.Path); err != nil {
				return err
			}
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

func testAccPipelineRemoteAttributes(n string, pipeline *buddy.Pipeline, want *testAccPipelineRemoteExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPipelineId, _ := strconv.Atoi(attrs["pipeline_id"])
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
		if err := util.CheckIntFieldSet("Id", pipeline.Id); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("pipeline_id", attrsPipelineId, pipeline.Id); err != nil {
			return err
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
		if err := util.CheckFieldEqualAndSet("definition_source", attrs["definition_source"], want.DefinitionSource); err != nil {
			return err
		}
		if err := testAccPipelineGitConfig(pipeline, attrs, want.GitConfigRef, want.GitConfig); err != nil {
			return err
		}
		if want.DefinitionSource == buddy.PipelineDefinitionSourceRemote {
			if err := util.CheckFieldEqualAndSet("DefinitionSource", pipeline.DefinitionSource, want.DefinitionSource); err != nil {
				return err
			}
		}
		if want.RemoteProjectName != "" {
			if err := util.CheckFieldEqualAndSet("remote_project_name", attrs["remote_project_name"], want.RemoteProjectName); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteProjectName", pipeline.RemoteProjectName, want.RemoteProjectName); err != nil {
				return err
			}
		}
		if want.RemoteRef != "" {
			if err := util.CheckFieldEqualAndSet("remote_ref", attrs["remote_ref"], want.RemoteRef); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("remote_branch", attrs["remote_branch"], want.RemoteRef); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteRef", pipeline.RemoteRef, want.RemoteRef); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteBranch", pipeline.RemoteBranch, want.RemoteRef); err != nil {
				return err
			}
		}
		if want.RemoteBranch != "" {
			if err := util.CheckFieldEqualAndSet("remote_ref", attrs["remote_ref"], want.RemoteBranch); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("remote_branch", attrs["remote_branch"], want.RemoteBranch); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteRef", pipeline.RemoteRef, want.RemoteBranch); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteBranch", pipeline.RemoteBranch, want.RemoteBranch); err != nil {
				return err
			}
		}
		if want.RemotePath != "" {
			if err := util.CheckFieldEqualAndSet("remote_path", attrs["remote_path"], want.RemotePath); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemotePath", pipeline.RemotePath, want.RemotePath); err != nil {
				return err
			}
		}
		if want.RemoteParam != "" {
			if err := util.CheckFieldEqualAndSet("remote_parameter.0.key", attrs["remote_parameter.0.key"], "cmd"); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("remote_parameter.0.value", attrs["remote_parameter.0.value"], want.RemoteParam); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteParameters[0].Key", pipeline.RemoteParameters[0].Key, "cmd"); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("RemoteParameters[0].Value", pipeline.RemoteParameters[0].Value, want.RemoteParam); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccPipelineAttributes(n string, pipeline *buddy.Pipeline, want *testAccPipelineExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPipelineId, _ := strconv.Atoi(attrs["pipeline_id"])
		attrsConcurrentPipelineRuns, _ := strconv.ParseBool(attrs["concurrent_pipeline_runs"])
		attrsDescriptionRequired, _ := strconv.ParseBool(attrs["description_required"])
		attrsAlwaysFromScratch, _ := strconv.ParseBool(attrs["always_from_scratch"])
		attrsFailOnPrepareEnvWarning, _ := strconv.ParseBool(attrs["fail_on_prepare_env_warning"])
		attrsFetchAllRefs, _ := strconv.ParseBool(attrs["fetch_all_refs"])
		attrsAutoClearCache, _ := strconv.ParseBool(attrs["auto_clear_cache"])
		attrsNoSkipToMostRecent, _ := strconv.ParseBool(attrs["no_skip_to_most_recent"])
		attrsDoNotCreateCommitStatus, _ := strconv.ParseBool(attrs["do_not_create_commit_status"])
		attrsIgnoreFailOnProjectStatus, _ := strconv.ParseBool(attrs["ignore_fail_on_project_status"])
		attrsPausedFailures, _ := strconv.Atoi(attrs["pause_on_repeated_failures"])
		attrsCloneDepth, _ := strconv.Atoi(attrs["clone_depth"])
		attrsPaused, _ := strconv.ParseBool(attrs["paused"])
		attrsCreatorMemberId, _ := strconv.Atoi(attrs["creator.0.member_id"])
		attrsDisabled, _ := strconv.ParseBool(attrs["disabled"])
		attrsManagePermissionsByYaml, _ := strconv.ParseBool(attrs["manage_permissions_by_yaml"])
		attrsManageVariablesByYaml, _ := strconv.ParseBool(attrs["manage_variables_by_yaml"])
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
		if err := util.CheckBoolFieldEqual("DescriptionRequired", pipeline.DescriptionRequired, want.DescriptionRequired); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("ManagePermissionsByYaml", pipeline.ManagePermissionsByYaml, want.ManagePermissionsByYaml); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("manage_permissions_by_yaml", attrsManagePermissionsByYaml, want.ManagePermissionsByYaml); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("ManageVariablesByYaml", pipeline.ManageVariablesByYaml, want.ManageVariablesByYaml); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("manage_variables_by_yaml", attrsManageVariablesByYaml, want.ManageVariablesByYaml); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("ConcurrentPipelineRuns", pipeline.ConcurrentPipelineRuns, want.ConcurrentPipelineRuns); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AlwaysFromScratch", pipeline.AlwaysFromScratch, want.AlwaysFromScratch); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("always_from_scratch", attrsAlwaysFromScratch, want.AlwaysFromScratch); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("description_required", attrsDescriptionRequired, want.DescriptionRequired); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("concurrent_pipeline_runs", attrsConcurrentPipelineRuns, want.ConcurrentPipelineRuns); err != nil {
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
		if want.Cpu != "" {
			if err := util.CheckFieldEqualAndSet("cpu", attrs["cpu"], want.Cpu); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Cpu", pipeline.Cpu, want.Cpu); err != nil {
				return err
			}
		}
		if want.PauseOnRepeatedFailures != 0 {
			if err := util.CheckIntFieldEqualAndSet("pause_on_repeated_failures", attrsPausedFailures, want.PauseOnRepeatedFailures); err != nil {
				return err
			}
			if err := util.CheckIntFieldEqualAndSet("PauseOnRepeatedFailures", pipeline.PauseOnRepeatedFailures, want.PauseOnRepeatedFailures); err != nil {
				return err
			}
		}
		if want.CloneDepth != 0 {
			if err := util.CheckIntFieldEqualAndSet("clone_depth", attrsCloneDepth, want.CloneDepth); err != nil {
				return err
			}
			if err := util.CheckIntFieldEqualAndSet("CloneDepth", pipeline.CloneDepth, want.CloneDepth); err != nil {
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
		if want.Identifier != "" {
			if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], want.Identifier); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Identifier", pipeline.Identifier, want.Identifier); err != nil {
				return err
			}
		}
		if want.ExecutionMessageTemplate != "" {
			if err := util.CheckFieldEqualAndSet("execution_message_template", attrs["execution_message_template"], want.ExecutionMessageTemplate); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("ExecutionMessageTemplate", pipeline.ExecutionMessageTemplate, want.ExecutionMessageTemplate); err != nil {
				return err
			}
		}
		if want.GitChangesetBase != "" {
			if err := util.CheckFieldEqualAndSet("git_changeset_base", attrs["git_changeset_base"], want.GitChangesetBase); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("GitChangesetBase", pipeline.GitChangesetBase, want.GitChangesetBase); err != nil {
				return err
			}
		}
		if want.FilesystemChangesetBase != "" {
			if err := util.CheckFieldEqualAndSet("filesystem_changeset_base", attrs["filesystem_changeset_base"], want.FilesystemChangesetBase); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("FilesystemChangesetBase", pipeline.FilesystemChangesetBase, want.FilesystemChangesetBase); err != nil {
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
		if err := util.CheckBoolFieldEqual("disabled", attrsDisabled, want.Disabled); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Disabled", pipeline.Disabled, want.Disabled); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("disabling_reason", attrs["disabling_reason"], want.DisablingReason); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("DisabledReason", pipeline.DisabledReason, want.DisablingReason); err != nil {
			return err
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
		if want.Permissions != nil {
			if err := util.CheckFieldEqualAndSet("Permissions.Others", pipeline.Permissions.Others, want.Permissions.Others); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("permissions.0.others", attrs["permissions.0.others"], want.Permissions.Others); err != nil {
				return err
			}
			if err := util.CheckIntFieldEqual("Permissions.Users", len(pipeline.Permissions.Users), len(want.Permissions.Users)); err != nil {
				return err
			}
			if err := util.CheckIntFieldEqual("Permissions.Groups", len(pipeline.Permissions.Groups), len(want.Permissions.Groups)); err != nil {
				return err
			}
			if len(want.Permissions.Users) > 0 {
				attrsPermUserId, _ := strconv.Atoi(attrs["permissions.0.user.0.id"])
				wantPermUserId := want.PermissionUser.Id
				if err := util.CheckFieldEqualAndSet("Permissions.Users[0].AccessLevel", pipeline.Permissions.Users[0].AccessLevel, want.Permissions.Users[0].AccessLevel); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("permissions.0.user.0.access_level", attrs["permissions.0.user.0.access_level"], want.Permissions.Users[0].AccessLevel); err != nil {
					return err
				}
				if err := util.CheckIntFieldEqual("Permissions.Users[0].Id", pipeline.Permissions.Users[0].Id, wantPermUserId); err != nil {
					return err
				}
				if err := util.CheckIntFieldEqual("permissions.0.user.0.id", attrsPermUserId, wantPermUserId); err != nil {
					return err
				}
			}
			if len(want.Permissions.Groups) > 0 {
				attrsPermGroupId, _ := strconv.Atoi(attrs["permissions.0.group.0.id"])
				wantPermGroupId := want.PermissionGroup.Id
				if err := util.CheckIntFieldEqual("Permissions.Groups[0].Id", pipeline.Permissions.Groups[0].Id, wantPermGroupId); err != nil {
					return err
				}
				if err := util.CheckIntFieldEqual("permissions.0.group.0.id", attrsPermGroupId, wantPermGroupId); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("Permissions.Groups[0].AccessLevel", pipeline.Permissions.Groups[0].AccessLevel, want.Permissions.Groups[0].AccessLevel); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("permissions.0.group.0.access_level", attrs["permissions.0.group.0.access_level"], want.Permissions.Groups[0].AccessLevel); err != nil {
					return err
				}
			}
		}
		if want.Event != nil {
			if err := util.CheckFieldEqualAndSet("Events[0].Type", pipeline.Events[0].Type, want.Event.Type); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("event.0.type", attrs["event.0.type"], want.Event.Type); err != nil {
				return err
			}
			switch want.Event.Type {
			case buddy.PipelineEventTypePullRequest:
				if err := util.CheckFieldEqualAndSet("Events[0].Branches[0]", pipeline.Events[0].Branches[0], want.Event.Branches[0]); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("event.0.branches.0", attrs["event.0.branches.0"], want.Event.Branches[0]); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("Events[0].Events[0]", pipeline.Events[0].Events[0], want.Event.Events[0]); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("event.0.events.0", attrs["event.0.events.0"], want.Event.Events[0]); err != nil {
					return err
				}
			case buddy.PipelineEventTypeWebhook:
				if err := util.CheckBoolFieldEqual("Events[0].Totp", pipeline.Events[0].Totp, want.Event.Totp); err != nil {
					return err
				}
			case buddy.PipelineEventTypeEmail:
				if err := util.CheckFieldEqualAndSet("Events[0].Prefix", pipeline.Events[0].Prefix, want.Event.Prefix); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("Events[0].Whitelist[0]", pipeline.Events[0].Whitelist[0], want.Event.Whitelist[0]); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("Events[0].Prefix", attrs["event.0.prefix"], want.Event.Prefix); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("Events[0].Whitelist[0]", attrs["event.0.whitelist.0"], want.Event.Whitelist[0]); err != nil {
					return err
				}
			case buddy.PipelineEventTypeSchedule:
				if want.Event.StartDate != "" {
					if err := util.CheckDateFieldEqual("event.0.start_date", attrs["event.0.start_date"], want.Event.StartDate); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet("Events[0].StartDate", pipeline.Events[0].StartDate, want.Event.StartDate); err != nil {
						return err
					}
				}
				if want.Event.Cron != "" {
					if err := util.CheckFieldEqualAndSet("event.0.cron", attrs["event.0.cron"], want.Event.Cron); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet("Events[0].Cron", pipeline.Events[0].Cron, want.Event.Cron); err != nil {
						return err
					}
				}
				if want.Event.Timezone != "" {
					if err := util.CheckFieldEqualAndSet("event.0.timezone", attrs["event.0.timezone"], want.Event.Timezone); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet("Events[0].Timezone", pipeline.Events[0].Timezone, want.Event.Timezone); err != nil {
						return err
					}
				}
				if want.Event.Delay > 0 {
					attrsDelay, _ := strconv.Atoi(attrs["event.0.delay"])
					if err := util.CheckIntFieldEqual("event.0.delay", attrsDelay, want.Event.Delay); err != nil {
						return err
					}
					if err := util.CheckIntFieldEqual("Events[0].Delay", pipeline.Events[0].Delay, want.Event.Delay); err != nil {
						return err
					}
				}
			default:
				if err := util.CheckFieldEqualAndSet("Events[0].Refs[0]", pipeline.Events[0].Refs[0], want.Event.Refs[0]); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("event.0.refs.0", attrs["event.0.refs.0"], want.Event.Refs[0]); err != nil {
					return err
				}
			}
		}
		if want.Loop != "" {
			if err := util.CheckIntFieldEqual("len(pipeline.Loop)", len(pipeline.Loop), 1); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("loop.#", attrs["loop.#"], "1"); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Loop[0]", pipeline.Loop[0], want.Loop); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("loop.0", attrs["loop.0"], want.Loop); err != nil {
				return err
			}
		} else {
			if err := util.CheckIntFieldEqual("len(pipeline.Loop)", len(pipeline.Loop), 0); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("loop.#", attrs["loop.#"] == "" || attrs["loop.#"] == "0", true); err != nil {
				return err
			}
		}
		if len(want.TriggerConditions) > 0 {
			for _, triggerCondition := range want.TriggerConditions {
				i := 0
				for {
					if i >= len(want.TriggerConditions) {
						return fmt.Errorf("trigger condition not found: %s", triggerCondition.TriggerCondition)
					}
					tc := attrs[fmt.Sprintf("trigger_condition.%d.condition", i)]
					if tc == triggerCondition.TriggerCondition {
						break
					}
					i += 1
				}
				if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerCondition", i), pipeline.TriggerConditions[i].TriggerCondition, triggerCondition.TriggerCondition); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.condition", i), attrs[fmt.Sprintf("trigger_condition.%d.condition", i)], triggerCondition.TriggerCondition); err != nil {
					return err
				}
				if triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionOnChangeAtPath {
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerConditionPaths[0]", i), pipeline.TriggerConditions[i].TriggerConditionPaths[0], triggerCondition.TriggerConditionPaths[0]); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.paths.0", i), attrs[fmt.Sprintf("trigger_condition.%d.paths.0", i)], triggerCondition.TriggerConditionPaths[0]); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionTriggeringUserIs || triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionTriggeringUserIsNot {
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerUser", i), pipeline.TriggerConditions[i].TriggerUser, triggerCondition.TriggerUser); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionTriggeringUserIsInGroup || triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionTriggeringUserIsNotInGroup {
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].TriggerGroup", i), pipeline.TriggerConditions[i].TriggerGroup, triggerCondition.TriggerGroup); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionVarIs ||
					triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionVarIsNot ||
					triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionVarContains ||
					triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionVarNotContains {
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
				if triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionDateTime {
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
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("TriggerConditions[%d].Timezone", i), pipeline.TriggerConditions[i].Timezone, triggerCondition.Timezone); err != nil {
						return err
					}
					if err := util.CheckFieldEqualAndSet(fmt.Sprintf("trigger_condition.%d.timezone", i), attrs[fmt.Sprintf("trigger_condition.%d.timezone", i)], triggerCondition.Timezone); err != nil {
						return err
					}
				}
				if triggerCondition.TriggerCondition == buddy.PipelineTriggerConditionSuccessPipeline {
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

func testAccPipelineGet(n string, pipeline *buddy.Pipeline) resource.TestCheckFunc {
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

func testAccPipelineConfigEventEmail(domain string, projectName string, name string, prefix string, whitelist string) string {
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
    event {
        type = "EMAIL"
				prefix = "%s"
				whitelist = ["%s"]
    }
}
`, domain, projectName, name, prefix, whitelist)
}

func testAccPipelineConfigEventWebhook(domain string, projectName string, name string, totp bool) string {
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
    event {
        type = "WEBHOOK"
				totp = %t
    }
}
`, domain, projectName, name, totp)
}

func testAccPipelineConfigEventPullRequest(domain string, projectName string, name string, eventType string, branch string, prEvent string) string {
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
    event {
        type = "%s"
        branches = ["%s"]
				events = ["%s"]
    }
}
`, domain, projectName, name, eventType, branch, prEvent)
}

func testAccPipelineConfigEvent(domain string, projectName string, name string, identifier string, eventType string, ref string, tcChangePath string, tcVarKey string, tcVarValue string, tcHours int, tcDays int, tcZoneId string, user string, group string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_member" "user" {
		domain = "${buddy_workspace.foo.domain}"
  	email  = "%s"
}

resource "buddy_group" "group" {
		domain = "${buddy_workspace.foo.domain}"
  	name = "%s"
}

resource "buddy_permission" "perm" {
  domain = "${buddy_workspace.foo.domain}"
  name                    = "perm"
  pipeline_access_level   = "RUN_ONLY"
  repository_access_level = "READ_ONLY"
  sandbox_access_level    = "DENIED"
}

resource "buddy_project_group" "group_in_project" {
  domain = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  group_id      = "${buddy_group.group.group_id}"
  permission_id = "${buddy_permission.perm.permission_id}"
}

resource "buddy_project_member" "user_in_project" {
  domain = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  member_id     = "${buddy_member.user.member_id}"
  permission_id = "${buddy_permission.perm.permission_id}"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
		identifier = "%s"
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
        timezone = "%s"
    }
		trigger_condition {
				condition = "TRIGGERING_USER_IS"
				trigger_user = "${buddy_member.user.email}"
		}
		trigger_condition {
				condition = "TRIGGERING_USER_IS_NOT"
				trigger_user = "${buddy_member.user.email}"
		}
		trigger_condition {
				condition = "TRIGGERING_USER_IS_IN_GROUP"
				trigger_group = "${buddy_group.group.name}"
		}
		trigger_condition {
				condition = "TRIGGERING_USER_IS_NOT_IN_GROUP"
				trigger_group = "${buddy_group.group.name}"
		}
}
`, domain, projectName, user, group, name, identifier, eventType, ref, tcChangePath, tcVarKey, tcVarValue, tcVarKey, tcVarValue, tcVarKey, tcVarValue, tcVarKey, tcVarValue, tcHours, tcDays, tcZoneId)
}

func testAccPipelinePermissionsEmpty(domain string, projectName string, name string, ref string) string {
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
	refs = ["%s"]
}
`, domain, projectName, name, ref)
}

func testAccPipelinePermissionsUser(domain string, projectName string, name string, ref string, loop string, email string, othersPerm string, userPerm string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
		refs = ["%s"]
		loop = ["%s"]
		permissions {
			others = "%s"
      user {
				id = "${buddy_project_member.bar.member_id}"
        access_level = "%s"
			}
		}
}
`, domain, email, projectName, name, ref, loop, othersPerm, userPerm)
}

func testAccPipelinePermissionsUserGroup(domain string, projectName string, name string, ref string, loop string, email string, groupName string, othersPerm string, userPerm string, groupPerm string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
		refs = ["%s"]
		loop = ["%s"]	
		permissions {
			others = "%s"
      	user {
				id = "${buddy_project_member.bar.member_id}"
        access_level = "%s"
			}
			group {
				id = "${buddy_project_group.bar.group_id}"
				access_level = "%s"
			}
		}
}
`, domain, email, groupName, projectName, name, ref, loop, othersPerm, userPerm, groupPerm)
}

func testAccPipelinePermissionsGroup(domain string, projectName string, name string, ref string, email string, groupName string, othersPerm string, groupPerm string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
		refs = ["%s"]
		loop = []
		permissions {
			others = "%s"
			group {
				id = "${buddy_project_group.bar.group_id}"
				access_level = "%s"
			}
		}
}
`, domain, email, groupName, projectName, name, ref, othersPerm, groupPerm)
}

func testAccPipelinePermissionsBackToEmpty(domain string, projectName string, name string, ref string, email string, groupName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_pipeline" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
	refs = ["%s"]
	permissions {}
}
`, domain, email, groupName, projectName, name, ref)
}

func testAccPipelineConfigClick(domain string, projectName string, name string, alwaysFromScratch bool, failOnPrepareEnvWarning bool, fetchAllRefs bool, autoClearCache bool, noSkipToMostRecent bool, doNotCreateCommitStatus bool, ignoreFailOnProjectStatus bool, executionMessageTemplate string, targetSiteUrl string, ref string, cloneDepth int, cpu string, managePermissionsByYaml bool, manageVariablesByYaml bool) string {
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
    clone_depth = %d
		cpu = "%s"
		manage_variables_by_yaml = %t
		manage_permissions_by_yaml = %t
}
`, domain, projectName, name, alwaysFromScratch, failOnPrepareEnvWarning, fetchAllRefs, autoClearCache, noSkipToMostRecent, doNotCreateCommitStatus, ignoreFailOnProjectStatus, executionMessageTemplate, targetSiteUrl, ref, cloneDepth, cpu, manageVariablesByYaml, managePermissionsByYaml)
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
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
