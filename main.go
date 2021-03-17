package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var (
	githubToken = flag.String("token", "", "GitHub API token.")
)

func main() {
	flag.Parse()

	if *githubToken == "" {
		log.Fatal("Missing GitHub API token")
	}

	ctx := context.Background()
	client := NewGitHubClient(ctx, *githubToken)

	repos, err := GetAllRepos(ctx, client)
	if err != nil {
		log.Fatal("List user repositories:", err)
	}

	var alfredItems []*AlfredItem
	for _, r := range repos {
		item := NewAlfredItem(r)
		alfredItems = append(alfredItems, item)
	}

	items := AlfredItems{Items: alfredItems}
	bytes, err := json.Marshal(items)
	if err != nil {
		log.Fatal("Serialization:", err)
	}

	_, err = fmt.Fprintf(os.Stdout, "%s", string(bytes))
	if err != nil {
		log.Fatal("Output:", err)
	}
}

func NewGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetAllRepos(ctx context.Context, client *github.Client) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

type AlfredItems struct {
	Items []*AlfredItem `json:"items"`
}

type AlfredItem struct {
	Uid      string `json:"uid"`
	Arg      string `json:"arg"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

func NewAlfredItem(r *github.Repository) *AlfredItem {
	var description string
	if r.Description == nil {
		description = ""
	} else {
		description = *r.Description
	}

	return &AlfredItem{
		Uid:      *r.NodeID,
		Arg:      *r.HTMLURL,
		Title:    *r.FullName,
		Subtitle: description,
	}
}
