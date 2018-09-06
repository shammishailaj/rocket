package zeitnow

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/astrocorp42/rocket/config"
	"github.com/astrocorp42/rocket/version"
	"github.com/astroflow/astroflow-go/log"
	"github.com/z0mbie42/fswalk"
)

// Client is an wrapper to perform various task against the zeit API
type Client struct {
	Token     string
	HTTP      *http.Client
	UserAgent string
	Config    config.ZeitNowConfig
}

type File struct {
	File string `json:"file"`
	SHA  string `json:"sha"`
	Size uint64 `json:"size"`
}

type DeploymentRequest struct {
	Env               map[string]string `json:"env,omitempty"`
	Public            bool              `json:"public"`
	ForceNew          *bool             `json:"forceNew,omitempty"`
	Name              string            `json:"name"`
	DeploymentType    string            `json:"deploymentType"`
	RegistryAuthToken *string           `json:"registryAuthToken,omitempty"`
	Files             []File            `json:"files"`
	Engines           map[string]string `json:"engines,omitempty"`
	SessionAffinity   *string           `json:"sessionAffinity,omitempty"`
	Config            map[string]string `json:"config,omitempty"`
}

type CreateDeploymentResponse struct {
	// TotalSize    int64         `json:"totalSize"`
	DeploymentID string        `json:"deploymentId"`
	URL          string        `json:"url"`
	Warnings     []interface{} `json:"warnings"` // TODO: find the correct schema
	//ReadyState string `json:"readyState"`
}

func Deploy(conf config.ZeitNowConfig) error {
	if conf.Token == nil {
		v := os.Getenv("ZEIT_TOKEN")
		conf.Token = &v
	} else {
		v := config.ExpandEnv(*conf.Token)
		conf.Token = &v
	}

	if conf.Directory == nil {
		v := "."
		conf.Directory = &v
	} else {
		v := config.ExpandEnv(*conf.Directory)
		conf.Directory = &v
	}

	if conf.Env == nil {
		v := map[string]string{}
		conf.Env = v
	}

	if conf.Public == nil {
		v := false
		conf.Public = &v
	}

	if conf.DeploymentType == nil {
		v := "NPM"
		conf.DeploymentType = &v
	} else {
		v := config.ExpandEnv(*conf.DeploymentType)
		conf.DeploymentType = &v
	}

	if conf.Name == nil {
		v := os.Getenv("ZEIT_NOW_NAME")
		conf.Name = &v
	} else {
		v := config.ExpandEnv(*conf.Name)
		conf.Name = &v
	}

	if conf.ForceNew == nil {
		v := true
		conf.ForceNew = &v
	}

	if conf.Engines == nil {
		v := map[string]string{}
		conf.Engines = v
	}

	if conf.SessionAffinity == nil {
		v := "ip"
		conf.SessionAffinity = &v
	} else {
		v := config.ExpandEnv(*conf.SessionAffinity)
		conf.SessionAffinity = &v
	}

	client := NewClient(conf, *conf.Token)
	filesToDeploy := []File{}

	walker, _ := fswalk.NewWalker()
	filesc, _ := walker.Walk(*conf.Directory)
	for file := range filesc {
		if file.Path == "." || file.IsDir || file.IsSymLink {
			continue
		}
		log.With("file", file.Path).Debug("zeit_now: file to upload")
		f, err := client.UploadFile(file.Path)
		if err != nil {
			log.With("file", file.Path).Error(fmt.Sprintf("zeit_now: error uploading a file: %s", err.Error()))
		} else {
			log.Info(fmt.Sprintf("zeit_now: file successfully uploaded %s", file.Path))
			filesToDeploy = append(filesToDeploy, f)
		}
	}

	log.With("files", filesToDeploy).Debug("zeit_now: creating deployment")
	depRes, err := client.CreateDeployment(filesToDeploy)
	if err != nil {
		log.Error(fmt.Sprintf("zeit_now: error creating deployment  %v", err))
	} else {
		log.Info(fmt.Sprintf("zeit_now: deployment successfully created %s", depRes.URL))
	}
	return err
}

func NewClient(conf config.ZeitNowConfig, token string) Client {
	return Client{token, &http.Client{}, fmt.Sprintf("rocket/%s", version.Version), conf}
}

func cleanFilePath(filePath, base string) string {
	if strings.Index(filePath, base) == 0 {
		ret := strings.Replace(filePath, base, "", 1)
		if string(ret[0]) == string(os.PathSeparator) {
			return ret[1:]
		} else {
			return ret
		}
	}
	return filePath
}

func (c *Client) UploadFile(file string) (File, error) {
	var ret File

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ret, err
	}
	reader := bytes.NewReader(data)

	h := sha1.New()
	h.Write(data)
	sum := h.Sum(nil)

	ret.SHA = fmt.Sprintf("%x", sum)
	ret.Size = uint64(len(data))
	ret.File = cleanFilePath(file, *c.Config.Directory)
	sizeStr := fmt.Sprintf("%d", ret.Size)

	req, err := http.NewRequest("POST", "https://api.zeit.co/v2/now/files", reader)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Length", sizeStr)
	req.Header.Set("x-now-digest", ret.SHA)
	req.Header.Set("x-now-size", sizeStr)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != 200 {
		return ret, errors.New(string(body))
	}

	return ret, nil
}

func (c *Client) CreateDeployment(files []File) (CreateDeploymentResponse, error) {
	var ret CreateDeploymentResponse

	request := DeploymentRequest{
		Env:             c.Config.Env,
		Public:          *c.Config.Public,
		DeploymentType:  *c.Config.DeploymentType,
		ForceNew:        c.Config.ForceNew,
		Files:           files,
		Name:            *c.Config.Name,
		Engines:         c.Config.Engines,
		SessionAffinity: c.Config.SessionAffinity,
	}

	jsonToPost, err := json.Marshal(request)
	if err != nil {
		return ret, err
	}
	log.With("request", string(jsonToPost)).Debug("zeit_now: create deployment request")

	req, err := http.NewRequest("POST", "https://api.zeit.co/v3/now/deployments", bytes.NewReader(jsonToPost))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
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

	log.With("response", string(body)).Debug("zeit_now: create deployment response received")

	if resp.StatusCode != 200 {
		return ret, errors.New(string(body))
	}

	err = json.Unmarshal(body, &ret)
	return ret, err
}
