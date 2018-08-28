package ghreleases

// TODO:
// create a draft release
// upload assets
// edit release: set draft as false
import (
	"context"

	"github.com/astrocorp42/rocket/config"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
}

func Deploy(conf config.Config) error {
	return nil
}

func NewClient(token string) (GitHubClient, error) {
	var err error

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(context.Background(), ts))

	return GitHubClient{client}, err
}

func CreateDraft() {

}

func UploadAssets() {

}

func publishRelease() {

}
