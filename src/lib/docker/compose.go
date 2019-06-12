package docker

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Pingflow/devtools/src/lib"
)

const bin = "docker-compose"

type ComposeFile string

type ComposeConfig struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}

type Services map[string]Service

type Service struct {
	Image         string   `yaml:"image"`
	Hostname      string   `yaml:"hostname"`
	ContainerName string   `yaml:"container_name"`
	Command       string   `yaml:"command"`
	Ports         []string `yaml:"ports"`
	DependsOn     []string `yaml:"depends_on"`
	Links         []string `yaml:"links"`
}

func Compose(composeFile string) ComposeFile {
	return ComposeFile(composeFile)
}

func (c ComposeFile) exec(args ...string) error {
	return lib.Exec(bin, append([]string{"-f", c.String(), "-p", "pf4"}, args...)...)
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
	return c.exec(append([]string{"logs", "-f",}, lib.RemoveEmptySlice(services)...)...)
}

func (c ComposeFile) TTY(service string, command ...string) error {
	return c.exec(append([]string{"exec", service}, command...)...)
}

func (c ComposeFile) Services() (Services, error) {

	dc := ComposeConfig{}

	data, e := ioutil.ReadFile(c.String())
	if e != nil {
		return nil, e
	}
	if e := yaml.Unmarshal(data, &dc); e != nil {
		return nil, e
	}

	return dc.Services, nil
}

func (c ComposeFile) ServicesStartWith(prefix string) (Services, error) {

	var s = Services{}

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

func (c ComposeFile) Service(service string) (Service, error) {
	l, e := c.Services()
	if e != nil {
		return Service{}, e
	}

	if v, ok := l[service]; ok {
		return v, nil
	}

	return Service{}, errors.New("service not found")
}

type PortsService struct {
	Host  int
	Local int
}

func (s Service) ListPorts() []PortsService {
	var p []PortsService

	for _, v := range s.Ports {
		var np PortsService
		vp := strings.Split(v, ":")
		if len(vp) == 2 {
			host, _ := strconv.Atoi(vp[0])
			np.Host = host
			local, _ := strconv.Atoi(vp[1])
			np.Local = local
		} else {
			port, _ := strconv.Atoi(vp[0])
			np.Host = port
			np.Local = port
		}

		p = append(p, np)
	}

	return p
}
