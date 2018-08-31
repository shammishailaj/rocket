package ghreleases

import (
	"context"
	"errors"
	"fmt"
	"os"
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
	} else {
		v := config.ExpandEnv(*conf.Name)
		conf.Name = &v
	}

	if conf.Body == nil {
		v := ""
		conf.Body = &v
	} else {
		v := config.ExpandEnv(*conf.Body)
		conf.Body = &v
	}

	if conf.Prerelease == nil {
		v := false
		conf.Prerelease = &v
	}

	if conf.Repo == nil {
		v := os.Getenv("ROCKET_GIT_REPO")
		conf.Repo = &v
	} else {
		v := config.ExpandEnv(*conf.Repo)
		conf.Repo = &v
	}

	if conf.APIKey == nil {
		v := os.Getenv("GITHUB_API_KEY")
		conf.APIKey = &v
	} else {
		v := config.ExpandEnv(*conf.APIKey)
		conf.APIKey = &v
	}

	if conf.Tag == nil {
		v := os.Getenv("ROCKET_LAST_TAG")
		conf.Tag = &v
	} else {
		v := config.ExpandEnv(*conf.Tag)
		conf.Tag = &v
	}

	if conf.Assets == nil {
		conf.Assets = []string{}
	}

	if conf.BaseURL == nil {
		v := os.Getenv("GITHUB_BASE_URL")
		conf.BaseURL = &v
	} else {
		v := config.ExpandEnv(*conf.BaseURL)
		conf.BaseURL = &v
	}

	if conf.UploadURL == nil {
		v := os.Getenv("GITHUB_UPLOAD_URL")
		conf.UploadURL = &v
	} else {
		conf.UploadURL = conf.BaseURL
	}

	if *conf.UploadURL != "" && *conf.BaseURL == "" {
		return errors.New("github: base_url should not be empty when upload_url is set")
	}

	repo, _ := parseRepo(*conf.Repo)
	client, err := NewClient(*conf.APIKey, *conf.BaseURL, *conf.UploadURL)
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

	log.With("files", files).Debug("github: uploading assets")
	err = client.UploadAssets(repo, releaseID, files)
	if err != nil {
		return err
	}

	log.Debug("github: publishing release")
	release, err := client.PublishRelease(repo, releaseID)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("github: release published %s", release.GetHTMLURL()))
	return nil
}

// NewClient create a GitHubClient instance with the given authentication information
func NewClient(token, baseURL, uploadURL string) (GitHubClient, error) {
	var err error
	var client *github.Client

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(context.Background(), ts)

	if baseURL == "" {
		client = github.NewClient(oauthClient)
	} else {
		client, err = github.NewEnterpriseClient(baseURL, uploadURL, oauthClient)
	}

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
		log.Info(fmt.Sprintf("github: deleting existing release %d", release.GetID()))
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

	log.Info(fmt.Sprintf("github: draft release created at %s", release.GetHTMLURL()))
	return release.GetID(), nil
}

// UploadAssets upload the given assets to the given release
func (c *GitHubClient) UploadAssets(repo GitHubRepo, releaseID int64, files []string) error {
	for _, file := range files {
		fileName := filepath.Base(file)
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
		log.With().Info(fmt.Sprintf("github: asset %s uploaded", file))
	}

	return nil
}

// PublishRelease publish the given release (set draft as false)
func (c *GitHubClient) PublishRelease(repo GitHubRepo, releaseID int64) (*github.RepositoryRelease, error) {
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
	return release, err
}

// parseRepo take as input a string in the forme "owner/repo" et return a GitHubRepo struct
func parseRepo(repo string) (GitHubRepo, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return GitHubRepo{}, errors.New("malformed GitHub repo")
	}

	return GitHubRepo{Owner: parts[0], Name: parts[1]}, nil
}
