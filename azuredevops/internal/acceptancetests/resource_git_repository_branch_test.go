//go:build (all || core || resource_git_repository_branch) && !exclude_resource_git_repository_branch
// +build all core resource_git_repository_branch
// +build !exclude_resource_git_repository_branch

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRepoBranch_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()
	branchNameChanged := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepoBranches(projectName, gitRepoName, branchName),
				Check: resource.ComposeTestCheckFunc(
					// test-branch
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "name", fmt.Sprintf("testbranch-%s", branchName)),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "ref_branch", "master"),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_master", "last_commit_id"),
					// test-branch2
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_commit_id", "name", fmt.Sprintf("testbranch2-%s", branchName)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "ref_commit_id"),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "last_commit_id"),
				),
			},
			// Test replace/update branch when name changes
			{
				Config: hclGitRepoBranches(projectName, gitRepoName, branchNameChanged),
				Check: resource.ComposeTestCheckFunc(
					// test-branch
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "name", fmt.Sprintf("testbranch-%s", branchNameChanged)),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "ref_branch", "master"),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_master", "last_commit_id"),
					// test-branch2
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_commit_id", "name", fmt.Sprintf("testbranch2-%s", branchNameChanged)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "ref_commit_id"),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "last_commit_id"),
				),
			},
		},
	},
	)
}

func TestAccGitRepoBranch_InvalidRef(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config:      hclGitRepoBranchInvalidRef(projectName, gitRepoName, branchName),
				ExpectError: regexp.MustCompile(`No refs found that match ref "refs/tags/0.0.0"`),
			},
		},
	},
	)
}

func hclGitRepoBranches(projectName, gitRepoName, branchName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository_branch" "from_master" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch-%[3]s"
  ref_branch    = "master"
}

resource "azuredevops_git_repository_branch" "from_commit_id" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch2-%[3]s"
  ref_commit_id = azuredevops_git_repository_branch.from_master.last_commit_id
}
  `, projectName, gitRepoName, branchName)
}

func hclGitRepoBranchInvalidRef(projectName, gitRepoName, branchName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository_branch" "from_master" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch-%[3]s"
  ref_branch    = "master"
}
resource "azuredevops_git_repository_branch" "from_commit_id" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch2-%[3]s"
  ref_commit_id = azuredevops_git_repository_branch.from_master.last_commit_id
}

resource "azuredevops_git_repository_branch" "from_nonexistent_tag" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch-non-existent-tag"
  ref_tag       = "0.0.0"
}


  `, projectName, gitRepoName, branchName)
}
