package ghreleases

import (
	"context"
	"errors"
	"os"
	// "os/exec"
	"path/filepath"
	"strings"

	"github.com/astrocorp42/rocket/config"
	"github.com/astroflow/astroflow-go/log"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubClient represent an authenticated GitHub clien to perform the API operations
type GitHubClient struct {
	client *github.Client
}

type GitHubRepo struct {
	Owner string
	Name  string
}

// Deploy perform the github release with the following steps:
// Create the release as draft
// upload assets
// publish the release (draft = false)
func Deploy(conf config.GitHubReleasesConfig) error {
	if conf.Name == nil {
		v := os.Getenv("ROCKET_LAST_TAG")
		conf.Name = &v
	}
	if conf.Body == nil {
		v := ""
		conf.Body = &v
	}
	if conf.Prerelease == nil {
		v := false
		conf.Prerelease = &v
	}
	if conf.Repo == nil {
		v := os.Getenv("ROCKET_GIT_REPO")
		conf.Repo = &v
	}
	if conf.APIKey == nil {
		v := os.Getenv("GITHUB_API_KEY")
		conf.APIKey = &v
	}
	if conf.Tag == nil {
		v := os.Getenv("ROCKET_LAST_TAG")
		conf.Tag = &v
	}
	repo, _ := parseRepo(*conf.Repo)
	client, err := NewClient(*conf.APIKey)
	if err != nil {
		return err
	}
	files := []string{}

	for _, pattern := range conf.Assets {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		files = append(files, matches...)
	}

	releaseID, err := client.CreateDraftRelease(
		repo,
		*conf.Name,
		strings.TrimSpace(*conf.Tag),
		*conf.Body,
		*conf.Prerelease,
	)
	if err != nil {
		return err
	}

	log.With("files", files).Debug("uploading assets")
	err = client.UploadAssets(repo, releaseID, files)
	if err != nil {
		return err
	}

	log.Debug("publishing release")
	err = client.PublishRelease(repo, releaseID)
	return err
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
func (c *GitHubClient) CreateDraftRelease(repo GitHubRepo, name, tag, body string, prerelease bool) (int64, error) {
	var release *github.RepositoryRelease
	ctx := context.Background()
	var err error
	var data = &github.RepositoryRelease{
		Name:       github.String(name),
		TagName:    github.String(tag),
		Body:       github.String(body),
		Draft:      github.Bool(true),
		Prerelease: github.Bool(prerelease),
	}

	release, _, err = c.client.Repositories.GetReleaseByTag(
		ctx,
		repo.Owner,
		repo.Name,
		tag,
	)
	if err == nil {
		log.With("tag", tag, "release_id", release.GetID()).Info("deleting existing release")
		_, err = c.client.Repositories.DeleteRelease(
			ctx,
			repo.Owner,
			repo.Name,
			release.GetID(),
		)
		if err != nil {
			return 0, err
		}
	}

	release, _, err = c.client.Repositories.CreateRelease(
		ctx,
		repo.Owner,
		repo.Name,
		data,
	)
	if err != nil {
		return 0, err
	}

	log.With("url", release.GetHTMLURL()).Info("draft release created")
	return release.GetID(), nil
}

// UploadAssets upload the given assets to the given release
func (c *GitHubClient) UploadAssets(repo GitHubRepo, releaseID int64, files []string) error {
	for _, file := range files {
		fileName := filepath.Base(file)
		log.With("file", fileName).Info("uploading asset")
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		_, _, err = c.client.Repositories.UploadReleaseAsset(
			context.Background(),
			repo.Owner,
			repo.Name,
			releaseID,
			&github.UploadOptions{
				Name: fileName,
			},
			f,
		)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

// PublishRelease publish the given release (set draft as false)
func (c *GitHubClient) PublishRelease(repo GitHubRepo, releaseID int64) error {
	var data = &github.RepositoryRelease{
		Draft: github.Bool(false),
	}

	release, _, err := c.client.Repositories.EditRelease(
		context.Background(),
		repo.Owner,
		repo.Name,
		releaseID,
		data,
	)
	if err != nil {
		return err
	}

	log.With("url", release.GetHTMLURL()).Info("release published")
	return nil
}

// parseRepo take as input a string in the forme "owner/repo" et return a GitHubRepo struct
func parseRepo(repo string) (GitHubRepo, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return GitHubRepo{}, errors.New("malformed GitHub repo")
	}

	return GitHubRepo{Owner: parts[0], Name: parts[1]}, nil
}
