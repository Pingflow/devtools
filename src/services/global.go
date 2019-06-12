package services

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/Pingflow/devtools/src/lib/docker"
	"github.com/Pingflow/devtools/src/lib/vault"
	"github.com/gobuffalo/envy"
)

var dc  docker.ComposeFile

func Start() error {

	dc = docker.Compose(envy.Get("PF4_DOCKER_COMPOSE", getPath("docker-compose.yml")))

	if e := dc.Up(); e != nil {
		return e
	}

	vc, e := vault.NewFromDockerCompose(dc, getPath("credentials.json"))
	if e != nil {
		return e
	}

	if _, e := vc.Init(); e != nil {
		return e
	}

	if _, e := vc.UnSeal(); e != nil {
		return e
	}

	return nil
}

func Stop() {
	if e := dc.Stop(); e != nil {
		log.Fatal(e)
	}
	os.Exit(0)
}

func ResetPF4() error {
	return os.RemoveAll(appPath())
}

func appPath() string {

	u, _ := user.Current()

	return fmt.Sprintf("%s/.pf4", u.HomeDir)
}

func getPath(path string) string {
	return fmt.Sprintf("%s/%s", appPath(), path)
}
