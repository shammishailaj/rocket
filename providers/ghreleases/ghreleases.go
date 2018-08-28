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

// GitHubClient represent an authenticated GitHub clien to perform the API operations
type GitHubClient struct {
	client *github.Client
}

// Deploy actually perform the github release and upload assets according to the given configuration
func Deploy(conf config.Config) error {
	return nil
}

// NewClient create a GitHubClient instance with the given authentication information
func NewClient(token string) (GitHubClient, error) {
	var err error

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(context.Background(), ts))

	return GitHubClient{client}, err
}

// CreateDraftRelease create a draft release with the given information
func CreateDraftRelease() {

}

// UploadAssets upload the given assets to the given release
func UploadAssets() {

}

// PublishRelease publish the given release (set draft as false)
func PublishRelease() {

}
