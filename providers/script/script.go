package script

import (
	"fmt"
	"os/exec"

	"github.com/astrocorp42/rocket/config"
)

// Deploy deploy the script part of the configuration
func Deploy(conf config.Config) error {
	if conf.Script == nil || *conf.Script == nil {
		return nil
	}

	for _, cmd := range *conf.Script {
		out, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	}

	return nil
}
