package zeitnow

import (
	"fmt"
	"net/http"
	"os"

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

	_ = NewClient(*conf.Token)

	walker, _ := fswalk.NewWalker()
	filesc, _ := walker.Walk(*conf.Directory)
	for file := range filesc {
		if file.Path == "." || file.IsDir || file.IsSymLink {
			continue
		}
		log.With("file", file.Path).Debug("zeit_now: file to upload")
	}

	return nil
}

func NewClient(token string) Client {
	return Client{token, &http.Client{}, fmt.Sprintf("rocket/%s", version.Version)}
}
