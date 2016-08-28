package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List("", nil)

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List("silverstripeltd", nil)

	if err != nil {
		panic(err)
	}

	for _, org := range orgs {
		fmt.Println(org)
	}

	fmt.Println("done")

}
