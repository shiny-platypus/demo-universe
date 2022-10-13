package main

import (
	"fmt"
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Member struct {
	UserName string
	Role     string
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a team
		animalsTeam, err := github.NewTeam(ctx, "animals", &github.TeamArgs{
			Description: pulumi.String("Welcome to the Zoo"),
			Name:        pulumi.String("animals"),
			Privacy:     pulumi.String("closed"),
		})
		if err != nil {
			return err
		}

		// Add members to animals team

		guin := Member{UserName: "guineveresaenger"}

		// unique name for TeamMembership
		utmName := "animals" + "-" + guin.UserName
		animalsTeamMember, err := github.NewTeamMembership(ctx, utmName, &github.TeamMembershipArgs{
			TeamId:   animalsTeam.ID(),
			Username: pulumi.String(guin.UserName),
		})

		//guin := Member{UserName: "guineveresaenger"}
		//
		//// unique name for TeamMembership
		//utmName := "animals" + "-" + guin.UserName
		//animalsTeamMember, err := github.NewTeamMembership(ctx, utmName, &github.TeamMembershipArgs{
		//	TeamId:   animalsTeam.ID(),
		//	Username: pulumi.String(guin.UserName),
		//})
		//
		//guin := Member{UserName: "guineveresaenger"}
		//
		//// unique name for TeamMembership
		//utmName := "animals" + "-" + guin.UserName
		//animalsTeamMember, err := github.NewTeamMembership(ctx, utmName, &github.TeamMembershipArgs{
		//	TeamId:   animalsTeam.ID(),
		//	Username: pulumi.String(guin.UserName),
		//})
		//
		//guin := Member{UserName: "guineveresaenger"}
		//
		//// unique name for TeamMembership
		//utmName := "animals" + "-" + guin.UserName
		//animalsTeamMember, err := github.NewTeamMembership(ctx, utmName, &github.TeamMembershipArgs{
		//	TeamId:   animalsTeam.ID(),
		//	Username: pulumi.String(guin.UserName),
		//})

		// Create a sub-team
		platypusTeam, err := github.NewTeam(ctx, "platypuses", &github.TeamArgs{
			Description: pulumi.String("Duck-billed mammals"),
			Name:        pulumi.String("platypuses"),
			ParentTeamId: animalsTeam.ID().ApplyT(func(id interface{}) int {
				x := fmt.Sprintf("%v", id) // we need the ID to be a real value, not an interface
				y, _ := strconv.Atoi(x)
				return y
			}).(pulumi.IntOutput),
			Privacy: pulumi.String("closed"),
		},
			pulumi.Parent(animalsTeam),
		)

		if err != nil {
			return err
		}

		// Create a repo

		testRepo, err := github.NewRepository(ctx, "test-repo", &github.RepositoryArgs{
			Name:     pulumi.String("test-repo"),
			AutoInit: pulumi.Bool(true),
		}, pulumi.Protect(false))

		_, err = github.NewBranchDefault(ctx, "main", &github.BranchDefaultArgs{
			Branch:     pulumi.String("main"),
			Repository: testRepo.Name,
		})

		// Add repository permissions

		perm := "admin"

		_, err = github.NewTeamRepository(ctx, fmt.Sprintf("%s-%s-%s", testRepo, platypusTeam.Name, perm), &github.TeamRepositoryArgs{
			Permission: pulumi.String(perm),
			Repository: testRepo.Name,
			TeamId:     platypusTeam.ID(),
		}, pulumi.Protect(false))

		if err != nil {
			return err
		}

		//Open PRs across a few pre-made repos, which will trigger an Action that is also pre-set-up.
		//
		managedRepos := []string{"development", "staging", "production"}

		for _, repo := range managedRepos {
			err = createPR(ctx, repo)
		}

		ctx.Export("pTeamName", platypusTeam.Name)
		ctx.Export("aTeamMember", animalsTeamMember.Username)
		//ctx.Export("reponame", testRepo.Name)
		//ctx.Export("rteam repo permission", teamRepo.TeamId)
		//ctx.Export("pull request number", pr.Number)
		return nil
	})

}

func createPR(ctx *pulumi.Context, repoName string) error {
	branchName := repoName + "-feature-branch"
	featureBranch, err := github.NewBranch(ctx, branchName, &github.BranchArgs{
		Branch:     pulumi.String("new-feature"),
		Repository: pulumi.String(repoName),
	})
	if err != nil {
		return err
	}

	fileName := repoName + "-contributing"
	contributing, err := github.NewRepositoryFile(ctx, fileName, &github.RepositoryFileArgs{
		Content:    pulumi.String("Remember - every repository should have a contributor's guide."),
		File:       pulumi.String("CONTRIBUTING.md"),
		Repository: pulumi.String(repoName),
		Branch:     featureBranch.Branch,
	}, pulumi.Parent(featureBranch))

	if err != nil {
		return err
	}

	//let's open a PR!
	prName := repoName + "-feature-pr"
	_, err = github.NewRepositoryPullRequest(ctx, prName, &github.RepositoryPullRequestArgs{
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

func addMembers(ctx *pulumi.Context, members []Member, team *github.Team, teamName string) error {
	for _, member := range members {
		// unique name for TeamMembership
		utmName := teamName + "-" + member.UserName
		// set up team maintainers if that Role is set
		if member.Role != "" {
			_, err := github.NewTeamMembership(ctx, utmName, &github.TeamMembershipArgs{
				TeamId:   team.ID(),
				Username: pulumi.String(member.UserName),
				Role:     pulumi.String(member.Role),
			})
			if err != nil {
				return err
			}
		} else { // default to non-maintainer team membership
			_, err := github.NewTeamMembership(ctx, utmName, &github.TeamMembershipArgs{
				TeamId:   team.ID(),
				Username: pulumi.String(member.UserName),
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}
