package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
)

// A message processes parameter and returns the result on responseChan.
// ctx is places in a struct, but this is ok to do.
type message struct {
	responseChan chan<- int
	parameter    string
	ctx          context.Context
}

func main() {
	var (
		owner      = flag.String("o", "owner", "owner name")
		repository = flag.String("r", "repository", "repositoryName")
	)
	flag.Parse()
	commitids, err := getAllCommitID(*owner, *repository)
	if err != nil {
		panic(err)
	}
	for _, commit := range commitids {
		fmt.Println(commit)
	}

}

func getAllCommitID(owner, repo string) ([]string, error) {

	var token = ""

	val, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Printf("GITHUB TOKEN not set")
	} else {
		token = val
	}

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
