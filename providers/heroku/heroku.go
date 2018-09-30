package heroku

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bloom42/rocket/config"
	"github.com/bloom42/rocket/version"
	"github.com/bloom42/astroflow-go/log"
	"github.com/z0mbie42/fswalk"
)

// CreateSourceResp is the response to the https://api.heroku.com/apps/{app}/builds API call
type CreateSourceResp struct {
	SourceBlob struct {
		GetURL string `json:"get_url"`
		PutURL string `json:"put_url"`
	} `json:"source_blob"`
}

// CreateBuildResp is the response to the https://api.heroku.com/apps/{app}/sources API call
type CreateBuildResp struct {
	App struct {
		ID string `json:"id"`
	} `json:"app"`
	BuildPacks []struct {
		URL string `json:"url"`
	} `json:"buildpacks"`
	CreatedAt       time.Time `json:"created_at"`
	ID              string    `json:"id"`
	OutputStreamURL string    `json:"output_stream_url"`
	Release         *string   `json:"release"`
	Slug            *struct {
		ID *string `json:"id"`
	} `json:"slug"`
	SourceBlob struct {
		Checksum           *string `json:"checksum"`
		URL                string  `json:"url"`
		Version            string  `json:"version"`
		VersionDescription *string `json:"version_description"`
	} `json:"source_blob"`
	Stack     string    `json:"stack"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	User      struct {
		Email string `json:"email"`
		ID    string `json:"id"`
	} `json:"user"`
}

type CreateBuildReq struct {
	SourceBlob CreateBuildSourceBlob `json:"source_blob"`
}

type CreateBuildSourceBlob struct {
	URL     string `json:"url"`
	Version string `json:"version"`
}

// Client is an wrapper to perform various task against the Heroku API
type Client struct {
	APIKey    string
	App       string
	HTTP      *http.Client
	UserAgent string
}

// Deploy deploy the script part of the configuration
// create an archive then release using the API
// https://devcenter.heroku.com/articles/build-and-release-using-the-api
func Deploy(conf config.HerokuConfig) error {
	if conf.App == nil {
		v := os.Getenv("HEROKU_APP")
		conf.App = &v
	} else {
		v := config.ExpandEnv(*conf.App)
		conf.App = &v
	}

	if conf.APIKey == nil {
		v := os.Getenv("HEROKU_API_KEY")
		conf.APIKey = &v
	} else {
		v := config.ExpandEnv(*conf.APIKey)
		conf.APIKey = &v
	}

	if conf.Directory == nil {
		v := "."
		conf.Directory = &v
	} else {
		v := config.ExpandEnv(*conf.Directory)
		conf.Directory = &v
	}

	if conf.Version == nil {
		v := os.Getenv("ROCKET_COMMIT_HASH")
		conf.Version = &v
	} else {
		v := config.ExpandEnv(*conf.Version)
		conf.Version = &v
	}

	// create the archive
	tmpFile, err := ioutil.TempFile("", "rocket.*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	// set up the gzip writer
	gw := gzip.NewWriter(tmpFile)
	tw := tar.NewWriter(gw)

	walker, _ := fswalk.NewWalker()
	filesc, _ := walker.Walk(*conf.Directory)
	for file := range filesc {
		if file.Path == "." || file.IsDir || file.IsSymLink {
			continue
		}
		log.With("archive", tmpFile.Name(), "file", file.Path).Debug("heroku: adding file to final archive")
		err = addFile(tw, file.Path)
		if err != nil {
			return err
		}
	}

	if err = tw.Close(); err != nil {
		return err
	}
	if err = gw.Close(); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	// upload it
	client := NewClient(*conf.APIKey, *conf.App)
	sourceRep, err := client.CreateSource()
	log.With("response", sourceRep).Debug("heroku: create source response")
	log.Info("heroku: source created")
	if err != nil {
		return err
	}

	err = client.UploadRelease(tmpFile.Name(), sourceRep.SourceBlob.PutURL)
	if err != nil {
		return err
	}
	log.Info("heroku: release uploaded")

	buildResp, err := client.CreateBuild(CreateBuildReq{SourceBlob: CreateBuildSourceBlob{URL: sourceRep.SourceBlob.GetURL, Version: *conf.Version}})
	if err != nil {
		return err
	}
	log.With("response", buildResp).Debug("heroku: create build response")
	log.Info("heroku: build created")

	return nil
}

func NewClient(apiKey, app string) Client {
	return Client{apiKey, app, &http.Client{}, fmt.Sprintf("rocket/%s", version.Version)}
}

func addFile(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// now lets create the header as needed for this file within the tarball
	header := new(tar.Header)
	header.Format = tar.FormatGNU
	header.Name = path
	header.Size = stat.Size()
	header.Mode = int64(stat.Mode())
	header.ModTime = stat.ModTime()
	// write the header to the tarball archive
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	// copy the file data to the tarball
	if _, err := io.Copy(tw, file); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateSource() (CreateSourceResp, error) {
	var ret CreateSourceResp

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.heroku.com/apps/%s/sources", c.App), nil)
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(body, &ret)
	return ret, err
}

func (c *Client) UploadRelease(file, putURL string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", putURL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(body) != 0 {
		return errors.New(string(body))
	}
	return nil
}

func (c *Client) CreateBuild(payload CreateBuildReq) (CreateBuildResp, error) {
	var ret CreateBuildResp
	var err error

	data, err := json.Marshal(&payload)
	if err != nil {
		return ret, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.heroku.com/apps/%s/builds", c.App), bytes.NewBuffer(data))
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("User-Agent", c.UserAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(body, &ret)
	return ret, err
}
