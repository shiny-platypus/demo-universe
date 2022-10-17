package main

import (
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//// Create a team
		//team, err := github.NewTeam(ctx, "universe", &github.TeamArgs{
		//	Description: pulumi.String("Welcome to the universe"),
		//	Name:        pulumi.String("universe"),
		//	Privacy:     pulumi.String("closed"),
		//})
		//
		//if err != nil {
		//	return err
		//}
		//ctx.Export("team name", team.Name)
		//
		////Add members to Universe team
		//
		//teamMember, err := github.NewTeamMembership(ctx, "joe-universe", &github.TeamMembershipArgs{
		//	TeamId:   team.ID(),
		//	Username: pulumi.String("joeduffy"),
		//}, pulumi.Parent(team))
		//if err != nil {
		//	return err
		//}
		//ctx.Export("team member", teamMember.Username)
		//
		//// Create and initialize a repo
		//
		//demoRepo, err := github.NewRepository(ctx, "andromeda-galaxy", &github.RepositoryArgs{
		//	Name:     pulumi.String("hello-universe-repo"),
		//	AutoInit: pulumi.Bool(true),
		//})
		//
		//if err != nil {
		//	return err
		//}
		//// Set default branch to main
		//_, err = github.NewBranchDefault(ctx, "main", &github.BranchDefaultArgs{
		//	Branch:     pulumi.String("main"),
		//	Repository: demoRepo.Name,
		//})
		//
		//if err != nil {
		//	return err
		//}
		//ctx.Export("reponame", demoRepo.Name)
		//
		//// Add repository permissions
		//
		//perm := "admin"
		//teamRepoResourceName := "andromeda-galaxy" + "-universe-" + perm
		//teamRepo, err := github.NewTeamRepository(ctx, teamRepoResourceName, &github.TeamRepositoryArgs{
		//	Permission: pulumi.String(perm),
		//	Repository: demoRepo.Name,
		//	TeamId:     team.ID(),
		//})
		//
		//if err != nil {
		//	return err
		//}
		//ctx.Export("team repo permission", teamRepo.TeamId)

		//Open PRs across a few existing repos

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
		Title:          pulumi.String("Add Contributor's Guide"),
		Body:           pulumi.String("This pull request serves to create a contributor guide on this repo."),
	}, pulumi.Parent(contributing))

	if err != nil {
		return err
	}
	return nil
}
