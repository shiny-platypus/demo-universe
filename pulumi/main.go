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
		}, pulumi.Protect(false))
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
			pulumi.Protect(false),
			pulumi.Parent(animalsTeam),
		)

		if err != nil {
			return err
		}

		fmt.Println(platypusTeam.Name)
		fmt.Println(animalsTeamMember.Username)
		return nil
	})

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
