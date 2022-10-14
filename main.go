package main

import (
	"fmt"
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
)

type Team struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	ParentTeamId int
	Teams        []Team   `yaml:"teams"`
	Slug         string   `yaml:"slug"`
	Members      []Member `yaml:"members"`
}
type Organization struct {
	Org   string `yaml:"org"`
	Teams []Team `yaml:"teams"`
}

type Member struct {
	UserName string `yaml:"username"`
	Role     string `yaml:"role"`
}

// TeamPermissions describes a github team and the levels of permissions available (i.e. admin, maintain)
type TeamPermissions struct {
	TeamName    string       `yaml:"team-name"`
	Permissions []Permission `yaml:"permissions"`
}

// Permission is a single permission and a list of repos with that permission.
type Permission struct {
	Role  string   `yaml:"role"`
	Repos []string `yaml:"repos"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//import accurate team info from yaml
		teamsFilePath, err := filepath.Abs("./teams.yaml")
		if err != nil {
			return err
		}

		yamlFile, err := os.ReadFile(teamsFilePath)
		if err != nil {
			return err
		}

		var org Organization
		err = yaml.Unmarshal(yamlFile, &org)
		if err != nil {
			return err
		}

		for _, orgTeam := range org.Teams {
			err = setupTeams(ctx, &orgTeam)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func setupTeams(ctx *pulumi.Context, rootTeam *Team) error {
	// set up org team, i.e. engineering
	parent, err := github.NewTeam(ctx, rootTeam.Slug, &github.TeamArgs{
		Description: pulumi.String(rootTeam.Description),
		Name:        pulumi.String(rootTeam.Name),
		Privacy:     pulumi.String("closed"),
	}, pulumi.Protect(false))
	if err != nil {
		fmt.Println("encountered error creating new Pulumi github team: ", rootTeam.Name)
		return err
	}

	if len(rootTeam.Members) > 0 {
		addMembers(ctx, rootTeam.Members, parent, rootTeam.Name)
	}
	err = setupTeamRepos(ctx, parent, rootTeam.Slug)
	if err != nil {
		return err
	}

	//set up nested teams
	for _, childTeam := range rootTeam.Teams {
		// set each child team's parent team ID to the current team ID
		team, err := github.NewTeam(ctx, childTeam.Slug, &github.TeamArgs{
			Description: pulumi.String(childTeam.Description),
			Name:        pulumi.String(childTeam.Name),
			Privacy:     pulumi.String("closed"),
			ParentTeamId: parent.ID().ApplyT(func(id interface{}) int {
				x := fmt.Sprintf("%v", id) // we brutally abuse the standard library here
				y, _ := strconv.Atoi(x)
				return y
			}).(pulumi.IntOutput),
		},
			pulumi.Parent(parent), // establish the parental relationship and show in the pulumi cli
			pulumi.Protect(false), // this is explicit
			pulumi.Aliases([]pulumi.Alias{{NoParent: pulumi.Bool(true)}}),
		)
		if err != nil {
			fmt.Println("encountered error creating new Pulumi github team: ", childTeam.Name)
			return err
		}

		err = addMembers(ctx, childTeam.Members, team, childTeam.Name)
		if err != nil {
			return err
		}

		err = setupTeamRepos(ctx, team, childTeam.Slug)
		if err != nil {
			return err
		}
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

func setupTeamRepos(ctx *pulumi.Context, team *github.Team, teamName string) error {

	repoFilePath, err := filepath.Abs("./team-repos/" + teamName + ".yaml")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	yamlFile, err := os.ReadFile(repoFilePath)

	if err != nil {
		return err
	}

	var teamPerms TeamPermissions
	err = yaml.Unmarshal(yamlFile, &teamPerms)
	if err != nil {
		return err
	}
	for _, permission := range teamPerms.Permissions {
		for _, repo := range permission.Repos {
			_, err := github.NewTeamRepository(ctx, fmt.Sprintf("%s-%s-%s", repo, teamName, permission.Role), &github.TeamRepositoryArgs{
				Permission: pulumi.String(permission.Role),
				Repository: pulumi.String(repo),
				TeamId:     team.ID(),
			}, pulumi.Protect(false))

			if err != nil {
				return err
			}
		}
	}
	return nil
}
