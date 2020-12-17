package main

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getReadme(client *github.Client, owner string) string {
	ctx := context.Background()
	readme, _, err := client.Repositories.GetReadme(ctx, owner, owner, nil)
	if err != nil {
		log.Fatal("Problem in getting readme information", err)
	}

	content, err := readme.GetContent()
	if err != nil {
		log.Fatal("Problem in getting readme content", err)
	}

	return content
}

func updateReadme(client *github.Client, owner string, repository string, newReadme string, commitMessage string) {
	fileContent := []byte(newReadme)

	commit, _, _, err := client.Repositories.GetContents(context.Background(), owner, repository, "README.md",
		&github.RepositoryContentGetOptions{},
	)
	if err != err {
		log.Fatal(err)
	}

	options := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: fileContent,
		SHA:     commit.SHA,
	}

	_, _, err = client.Repositories.UpdateFile(context.Background(), owner, repository, "README.md", options)
	if err != nil {
		log.Fatal(err)
	}
}
