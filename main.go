package main

import (
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a team
		team, err := github.NewTeam(ctx, "animals", &github.TeamArgs{
			Description: pulumi.String("Welcome to the Zoo"),
			Name:        pulumi.String("animals"),
			Privacy:     pulumi.String("closed"),
		})

		if err != nil {
			return err
		}
		ctx.Export("team name", team.Name)

		// Add members to animals team

		//teamMember, err := github.NewTeamMembership(ctx, "animals", &github.TeamMembershipArgs{
		//	TeamId:   team.ID(),
		//	Username: pulumi.String("guineveresaenger"),
		//}, pulumi.Parent(team))
		//if err != nil {
		//	return err
		//}
		//ctx.Export("team member", teamMember.Username)

		// Create and initialize a repo

		//testRepo, err := github.NewRepository(ctx, "test-repo", &github.RepositoryArgs{
		//	Name:     pulumi.String("test-repo"),
		//	AutoInit: pulumi.Bool(true),
		//})
		//
		//if err != nil {
		//	return err
		//}
		//// Set default branch to main
		//_, err = github.NewBranchDefault(ctx, "main", &github.BranchDefaultArgs{
		//	Branch:     pulumi.String("main"),
		//	Repository: testRepo.Name,
		//})
		//
		//if err != nil {
		//	return err
		//}
		//ctx.Export("reponame", testRepo.Name)

		// Add repository permissions

		//perm := "admin"
		//teamRepoResourceName := "test-repo" + "-animals-" + perm
		//teamRepo, err := github.NewTeamRepository(ctx, teamRepoResourceName, &github.TeamRepositoryArgs{
		//	Permission: pulumi.String(perm),
		//	Repository: testRepo.Name,
		//	TeamId:     team.ID(),
		//})
		//
		//if err != nil {
		//	return err
		//}
		//ctx.Export("team repo permission", teamRepo.TeamId)

		//Open PRs across a few existing repos

		//
		//managedRepos := []string{"development", "staging", "production"}
		//for _, repo := range managedRepos {
		//	err = createPR(ctx, repo)
		//}

		return nil
	})

}

func createPR(ctx *pulumi.Context, repoName string) error {
	branchResourceName := repoName + "-feature-branch"
	featureBranch, err := github.NewBranch(ctx, branchResourceName, &github.BranchArgs{
		Branch:     pulumi.String("new-feature"),
		Repository: pulumi.String(repoName),
	})
	if err != nil {
		return err
	}

	fileResourceName := repoName + "-contributing"
	contributing, err := github.NewRepositoryFile(ctx, fileResourceName, &github.RepositoryFileArgs{
		Content:    pulumi.String("Remember - every repository should have a contributor's guide."),
		File:       pulumi.String("CONTRIBUTING.md"),
		Repository: pulumi.String(repoName),
		Branch:     featureBranch.Branch,
	}, pulumi.Parent(featureBranch))

	if err != nil {
		return err
	}

	//let's open a PR!
	prResourceName := repoName + "-feature-pr"
	_, err = github.NewRepositoryPullRequest(ctx, prResourceName, &github.RepositoryPullRequestArgs{
		BaseRef:        pulumi.String("main"),
		BaseRepository: pulumi.String(repoName),
		HeadRef:        featureBranch.Branch,
		Title:          pulumi.String("Do you wanna make a PR?"),
		Body:           pulumi.String("This pull request serves to create a contributor guide on this repo."),
	}, pulumi.Parent(contributing))

	if err != nil {
		return err
	}
	return nil
}
