package src

import (
	"fmt"
	"github.com/Pingflow/devtools/config"
	"github.com/Pingflow/devtools/src/lib/vault"
	"github.com/c-bata/go-prompt"
	"log"
	"os"
	"os/user"

	"github.com/Pingflow/devtools/src/lib/docker"
	"github.com/gobuffalo/envy"
)

const name = "pf4"

var (
	dockerComposePath = getPath("docker-compose.yml")
	traefikPath       = getPath("traefik.toml")
	vaultPath         = getPath("vault.json")
	credentialsPath   = getPath("credentials.json")

	dc docker.ComposeFile
)

type (
	app struct {
		Version string
		Commit  string
		Date    string
	}

	IApp interface {
		Run() error
	}
)

func App(version, commit, date string) IApp {
	return app{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

func (a app) Run() error {
	fmt.Printf("[DT] GoLang DevTools for GoMicro %v (%v) built at %v\n", a.Version, a.Commit, a.Date)
	if e := start(); e != nil {
		return e
	}
	fmt.Println("\nPlease use `exit` or `Ctrl-D` to exit this program.")
	defer stop()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("dt> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionMaxSuggestion(10),
	)
	p.Run()
	return nil
}

func start() error {

	if _, e := os.Stat(appPath()); os.IsNotExist(e) {
		if e := os.MkdirAll(appPath(), os.ModePerm); e != nil {
			return e
		}
	}

	if e := config.DockerCompose.ToYAMLFile(dockerComposePath); e != nil {
		return e
	}

	if e := config.Vault.ToJSONFile(vaultPath); e != nil {
		return e
	}

	dc = docker.Compose(dockerComposePath)

	if e := dc.Up(); e != nil {
		return e
	}

	fmt.Println()

	vc, e := vault.NewFromDockerCompose(dc, credentialsPath)
	if e != nil {
		return e
	}

	if e := vc.Wait(); e != nil {
		return e
	}

	if _, e := vc.Init(); e != nil {
		return e
	}

	if e := vc.UnSeal(); e != nil {
		return e
	}

	return nil
}

func stop() {
	if e := dc.Stop(); e != nil {
		log.Fatal(e)
	}
	os.Exit(0)
}

func resetPF4() error {
	return os.RemoveAll(appPath())
}

func appPath() string {
	u, _ := user.Current()
	return fmt.Sprintf("%s/.devtools", envy.Get("DEVTOOLS_PATH", u.HomeDir))
}

func getPath(path string) string {
	return fmt.Sprintf("%s/%s", appPath(), path)
}
