package docker

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/Pingflow/devtools/config"
	"github.com/Pingflow/devtools/src/lib/exec"
	"github.com/Pingflow/devtools/src/lib/slice"
	"gopkg.in/yaml.v2"
)

const bin = "docker-compose"

type ComposeFile string

func Compose(composeFile string) ComposeFile {
	return ComposeFile(composeFile)
}

func (c ComposeFile) exec(args ...string) error {
	return exec.Run(bin, append([]string{"-f", c.String(), "-p", "pf4"}, args...)...)
}

func (c ComposeFile) String() string {
	return string(c)
}

func (c ComposeFile) Up() error {
	return c.exec("up", "-d")
}

func (c ComposeFile) Start() error {
	return c.exec("start")
}

func (c ComposeFile) Stop() error {
	return c.exec("stop")
}

func (c ComposeFile) Down() error {
	return c.exec("down")
}

func (c ComposeFile) Ps() error {
	return c.exec("ps")
}

func (c ComposeFile) Logs(services ...string) error {
	return c.exec(append([]string{"logs", "-f",}, slice.RemoveEmpty(services)...)...)
}

func (c ComposeFile) TTY(service string, command ...string) error {
	return c.exec(append([]string{"exec", service}, command...)...)
}

func (c ComposeFile) Services() (config.DockerComposeFormatServices, error) {

	dc := config.DockerComposeFormat{}

	data, e := ioutil.ReadFile(c.String())
	if e != nil {
		return nil, e
	}

	if e := yaml.Unmarshal(data, &dc); e != nil {
		return nil, e
	}

	return dc.Services, nil
}

func (c ComposeFile) ServicesStartWith(prefix string) (config.DockerComposeFormatServices, error) {

	var s = config.DockerComposeFormatServices{}

	l, e := c.Services()
	if e != nil {
		return nil, e
	}

	for k, svc := range l {
		if strings.HasPrefix(k, prefix) {
			s[k] = svc
		}
	}

	return s, nil
}

func (c ComposeFile) Service(service string) (config.DockerComposeFormatService, error) {
	l, e := c.Services()
	if e != nil {
		return config.DockerComposeFormatService{}, e
	}

	if v, ok := l[service]; ok {
		return v, nil
	}

	return config.DockerComposeFormatService{}, errors.New("service not found")
}
