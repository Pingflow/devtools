package config

import (
	"encoding/json"
	"io/ioutil"
)

var Vault = VaultFormat{
	Listener: []VaultFormatConfigListener{
		{
			TCP: VaultFormatConfigListenerTCP{
				Address:        "0.0.0.0:8200",
				ClusterAddress: "0.0.0.0:8201",
				TLSDisable:     true,
			},
		},
	},
	Storage: []VaultFormatConfigStorage{
		{
			Consul: VaultFormatConfigStorageConsul{
				Address: "consul:8500",
				Path:    "vault",
			},
		},
	},
	ClusterName: "vault",
	UI:          true,
}

type (
	VaultFormat struct {
		Listener    []VaultFormatConfigListener `json:"listener"`
		Storage     []VaultFormatConfigStorage  `json:"storage"`
		ClusterName string                      `json:"cluster_name"`
		UI          bool                        `json:"ui"`
	}

	VaultFormatConfigListener struct {
		TCP VaultFormatConfigListenerTCP `json:"tcp"`
	}

	VaultFormatConfigListenerTCP struct {
		Address        string `json:"address"`
		ClusterAddress string `json:"cluster_address"`
		TLSDisable     bool   `json:"tls_disable"`
	}

	VaultFormatConfigStorage struct {
		Consul VaultFormatConfigStorageConsul `json:"consul"`
	}

	VaultFormatConfigStorageConsul struct {
		Address string `json:"address"`
		Path    string `json:"path"`
	}
)

func (vf VaultFormat) ToJSON() ([]byte, error) {
	return json.Marshal(vf)
}

func (vf VaultFormat) ToJSONFile(filename string) error {
	b, e := vf.ToJSON()
	if e != nil {
		return e
	}

	return ioutil.WriteFile(filename, b, 0644)
}
