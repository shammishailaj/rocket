package script

import (
	"io"
	"os"
	"os/exec"

	"github.com/astrocorp42/rocket/config"
)

// Deploy deploy the script part of the configuration
func Deploy(conf config.Config) error {
	if conf.Script == nil || *conf.Script == nil {
		return nil
	}

	for _, script := range *conf.Script {
		var err error
		cmd := exec.Command("sh", "-c", script)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		go func() {
			io.Copy(os.Stdout, stdout)
		}()
		go func() {
			io.Copy(os.Stderr, stderr)
		}()

		err = cmd.Start()
		if err != nil {
			return err
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
