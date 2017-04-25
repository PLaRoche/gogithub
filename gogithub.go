package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
)

func main() {
	var (
		owner      = flag.String("o", "owner", "owner name")
		repository = flag.String("r", "repository", "repositoryName")
	)
	flag.Parse()

	token, err := getToken()
	if err != nil {
		panic(err)
	}
	// commitids, err := getAllCommitID(*owner, *repository, token)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, commit := range commitids {
	// 	fmt.Println(commit)
	// }

	url, err := makeIssue(*owner, *repository, token)
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

}

func getToken() (string, error) {

	val, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return "", errors.New("Token not found")
	}
	return val, nil
}

func getAllCommitID(owner, repo string, token string) ([]string, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	opt := &github.IssueListOptions{}
	ctx := context.Background()

	issueBody := []string{}
	for {
		issues, resp, err := client.Issues.List(ctx, true, opt)
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			issueBody = append(issueBody, issue.GetBody())
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return issueBody, nil
}

func makeIssue(owner, repo string, token string) (string, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	ctx := context.Background()

	title := "Release v1.4.0"
	message := "Release v1.4.0-rc.0"
	body := &github.IssueRequest{
		Title:    &title,
		Body:     &message,
		Assignee: &owner,
		Labels:   &[]string{"ops"},
	}

	issue, _, err := client.Issues.Create(ctx, owner, repo, body)
	if err != nil {
		return "Issues.Create returned error", err
	}

	return issue.GetHTMLURL(), nil

}
