package docker

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/astrocorp42/rocket/config"
)

func exe(script string) error {
	script = config.ExpandEnv(script)
	cmd := exec.Command("sh", "-c", script)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
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
	return nil

}

func Deploy(conf config.DockerConfig) error {
	var err error

	if conf.Username == nil {
		v := os.Getenv("DOCKER_USERNAME")
		conf.Username = &v
	} else {
		v := config.ExpandEnv(*conf.Username)
		conf.Username = &v
	}

	if conf.Password == nil {
		v := os.Getenv("DOCKER_PASSWORD")
		conf.Password = &v
	} else {
		v := config.ExpandEnv(*conf.Password)
		conf.Password = &v
	}

	if conf.Login == nil {
		v := true
		conf.Login = &v
	}

	if conf.Images == nil {
		conf.Images = []string{}
	}

	// actually deploy
	if *conf.Login == true {
		if err = exe(fmt.Sprintf("docker login -u %s -p %s", *conf.Username, *conf.Password)); err != nil {
			return err
		}
	}

	for _, image := range conf.Images {
		if err = exe(fmt.Sprintf("docker push %s", config.ExpandEnv(image))); err != nil {
			return err
		}
	}

	return nil
}
