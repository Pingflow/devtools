package config

import (
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var DockerCompose = DockerComposeFormat{
	Version: "2.1",
	Services: map[string]DockerComposeFormatService{
		"consul-1": {
			Image:         "consul:1.5",
			Hostname:      "consul-1",
			ContainerName: "consul-1",
			Command:       "consul agent -dev -bind=0.0.0.0 -client=0.0.0.0 -ui",
			Ports: []string{
				"8500:8500",
			},
			HealthCheck: DockerComposeFormatServiceHealthCheck{
				Test:     []string{"CMD", "curl", "-f", "http://localhost:8500/v1/status/leader"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  10,
			},
		},
		"consul-2": {
			Image:         "consul:1.5",
			Hostname:      "consul-2",
			ContainerName: "consul-2",
			Command:       "consul agent -dev -bind=0.0.0.0 -client=0.0.0.0 -join=consul-1 -ui",
			Ports: []string{
				"8501:8500",
			},
			HealthCheck: DockerComposeFormatServiceHealthCheck{
				Test:     []string{"CMD", "curl", "-f", "http://localhost:8500/v1/status/leader"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  10,
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-1": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-1",
			},
		},
		"consul-3": {
			Image:         "consul:1.5",
			Hostname:      "consul-3",
			ContainerName: "consul-3",
			Command:       "consul agent -dev -bind=0.0.0.0 -client=0.0.0.0 -join=consul-1 -ui",
			Ports: []string{
				"8502:8500",
			},
			HealthCheck: DockerComposeFormatServiceHealthCheck{
				Test:     []string{"CMD", "curl", "-f", "http://localhost:8500/v1/status/leader"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  10,
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-1": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-1",
			},
		},

		"vault-1": {
			Image:         "vault:1.1.2",
			Hostname:      "vault-1",
			ContainerName: "vault-1",
			Command:       "vault server -config /config/vault.json",
			Volumes: []string{
				"./vault.json:/config/vault.json",
			},
			Ports: []string{
				"8200:8200",
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-1": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-1:consul",
			},
			CapAdd: []string{
				"IPC_LOCK",
			},
		},

		"vault-2": {
			Image:         "vault:1.1.2",
			Hostname:      "vault-2",
			ContainerName: "vault-2",
			Command:       "vault server -config /config/vault.json",
			Volumes: []string{
				"./vault.json:/config/vault.json",
			},
			Ports: []string{
				"8201:8200",
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-2": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-2:consul",
			},
			CapAdd: []string{
				"IPC_LOCK",
			},
		},

		"vault-3": {
			Image:         "vault:1.1.2",
			Hostname:      "vault-3",
			ContainerName: "vault-3",
			Command:       "vault server -config /config/vault.json",
			Volumes: []string{
				"./vault.json:/config/vault.json",
			},
			Ports: []string{
				"8202:8200",
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-3": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-3:consul",
			},
			CapAdd: []string{
				"IPC_LOCK",
			},
		},

		"web": {
			Image:         "microhq/micro",
			Hostname:      "web",
			ContainerName: "web",
			Command:       "--enable_stats --registry=consul --registry_address=\"consul:8500\" --server_name=\"com.pingflow.web\" web --namespace=\"com.pingflow\"",
			Ports: []string{
				"8001:8082",
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-2": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-2:consul",
			},
		},

		"api": {
			Image:         "microhq/micro",
			Hostname:      "api",
			ContainerName: "api",
			Command:       "--enable_stats --registry=consul --registry_address=\"consul:8500\" --server_name=\"com.pingflow.api\" api --handler=web --namespace=\"com.pingflow\"",
			Ports: []string{
				"8000:8080",
			},
			DependsOn: map[string]DockerComposeFormatServiceDependsOn{
				"consul-3": {
					Condition: "service_healthy",
				},
			},
			Links: []string{
				"consul-3:consul",
			},
		},
	},
}

type (
	DockerComposeFormat struct {
		Version  string                      `yaml:"version"`
		Services DockerComposeFormatServices `yaml:"services"`
	}

	DockerComposeFormatServices map[string]DockerComposeFormatService

	DockerComposeFormatService struct {
		Image         string                                         `yaml:"image,omitempty"`
		Hostname      string                                         `yaml:"hostname,omitempty"`
		ContainerName string                                         `yaml:"container_name,omitempty"`
		Command       string                                         `yaml:"command,omitempty"`
		Ports         []string                                       `yaml:"ports,omitempty"`
		DependsOn     map[string]DockerComposeFormatServiceDependsOn `yaml:"depends_on,omitempty"`
		Links         []string                                       `yaml:"links,omitempty"`
		HealthCheck   DockerComposeFormatServiceHealthCheck          `yaml:"healthcheck,omitempty"`
		Volumes       []string                                       `yaml:"volumes,omitempty"`
		CapAdd        []string                                       `yaml:"cap_add,omitempty"`
	}

	DockerComposeFormatServiceDependsOn struct {
		Condition string `yaml:"condition,omitempty"`
	}
	DockerComposeFormatServiceHealthCheck struct {
		Test     []string `yaml:"test,flow,omitempty"`
		Interval string   `yaml:"interval,omitempty"`
		Timeout  string   `yaml:"timeout,omitempty"`
		Retries  int      `yaml:"retries,omitempty"`
	}
)

func (dcf DockerComposeFormat) ToYAML() ([]byte, error) {
	return yaml.Marshal(dcf)
}

func (dcf DockerComposeFormat) ToYAMLFile(filename string) error {
	b, e := dcf.ToYAML()
	if e != nil {
		return e
	}

	return ioutil.WriteFile(filename, b, 0644)
}

type DockerComposeFormatServicePorts struct {
	Host  int
	Local int
}

func (s DockerComposeFormatService) ListPorts() []DockerComposeFormatServicePorts {
	var p []DockerComposeFormatServicePorts

	for _, v := range s.Ports {
		vp := strings.Split(v, ":")
		if len(vp) == 2 {
			host, _ := strconv.Atoi(vp[0])
			local, _ := strconv.Atoi(vp[1])
			p = append(p, DockerComposeFormatServicePorts{Host: host, Local: local})
		} else {
			port, _ := strconv.Atoi(vp[0])
			p = append(p, DockerComposeFormatServicePorts{Host: port, Local: port})

		}
	}

	return p
}
