package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// IssueMessage : a simple struct to hold all the components of an issue
type IssueMessage struct {
	Release string //Message : of issue
}

func main() {
	var (
		owner      = flag.String("o", "owner", "owner name")
		repository = flag.String("r", "repository", "repositoryName")
		version    = flag.String("v", "version", "version tag")
		fileName   = flag.String("f", "tmpl/issue.txt", "template file")
	)
	flag.Parse()

	token, err := getToken()
	if err != nil {
		panic(err)
	}

	url, err := makeIssue(*owner, *repository, *version, *fileName, token)
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

func makeIssue(owner, repo, version, templateFile, token string) (string, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	ctx := context.Background()

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return "Template Parse File Issue: ", err
	}
	ghIssue := IssueMessage{version}

	var buf bytes.Buffer

	err = t.ExecuteTemplate(&buf, "Message", ghIssue)
	if err != nil {
		return "Could not make the Message: ", err
	}
	message := buf.String()
	buf.Reset()

	err = t.ExecuteTemplate(&buf, "Title", ghIssue)
	if err != nil {
		return "Could not make the title: ", err
	}
	title := buf.String()
	buf.Reset()

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
