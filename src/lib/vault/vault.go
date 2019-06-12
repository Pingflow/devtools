package vault

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Pingflow/devtools/src/lib"
	"github.com/Pingflow/devtools/src/lib/docker"
	"github.com/hashicorp/vault/api"
)

type IClients interface {
	IsInitialized() (bool, error)
	IsSealed() (bool, error)
	Init() (*api.InitResponse, error)
	UnSeal() (*api.SealStatusResponse, error)
	Leader() (*api.LeaderResponse, error)
}

type client struct {
	credentialPath string
	clients        []string
}

func NewFromDockerCompose(compose docker.ComposeFile, credentialPath string) (IClients, error) {

	var c []string

	s, e := compose.ServicesStartWith("vault")
	if e != nil {
		return nil, e
	}

	for _, v := range s {
		for _, p := range v.ListPorts() {
			if p.Local == 8200 {
				c = append(c, fmt.Sprintf("http://127.0.0.1:%d", p.Host))
			}
		}
	}

	return New(credentialPath, c...), nil
}

func New(credentialPath string, clients ...string) IClients {
	return client{
		credentialPath: credentialPath,
		clients:        clients,
	}
}

func (c client) Init() (*api.InitResponse, error) {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return nil, e
	}

	if e := lib.WaitForServices([]string{strings.TrimPrefix(c.clients[0], "http://")}, time.Second*10); e != nil {
		return nil, e
	}

	if e := vc.SetAddress(c.clients[0]); e != nil {
		return nil, e
	}

	init, e := vc.Sys().InitStatus()
	if e != nil {
		return nil, e
	}

	if !init {
		r, e := vc.Sys().Init(&api.InitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		})
		if e != nil {
			return nil, e
		}

		b, e := json.Marshal(r)
		if e != nil {
			return nil, e
		}

		if e := ioutil.WriteFile(c.credentialPath, b, 0644); e != nil {
			return nil, e
		}

		return r, nil
	}

	return nil, nil
}

func (c client) UnSeal() (*api.SealStatusResponse, error) {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return nil, e
	}

	if _, e := os.Stat(c.credentialPath); os.IsNotExist(e) {
		return nil, os.ErrNotExist
	}

	var credential api.InitResponse
	b, e := ioutil.ReadFile(c.credentialPath)
	if e != nil {
		return nil, e
	}

	if e := json.Unmarshal(b, &credential); e != nil {
		return nil, e
	}

	for _, v := range c.clients {
		if e := vc.SetAddress(v); e != nil {
			return nil, e
		}

		seal, e := vc.Sys().SealStatus()
		if e != nil {
			return nil, e
		}

		if seal.Sealed {
			for _, k := range credential.KeysB64 {
				return vc.Sys().Unseal(k)
			}
		}
	}

	return nil, nil
}

func (c client) Leader() (*api.LeaderResponse, error) {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return nil, e
	}

	if e := vc.SetAddress(c.clients[0]); e != nil {
		return nil, e
	}

	return vc.Sys().Leader()
}

func (c client) IsInitialized() (bool, error) {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return false, e
	}

	if e := vc.SetAddress(c.clients[0]); e != nil {
		return false, e
	}

	return vc.Sys().InitStatus()
}

func (c client) IsSealed() (bool, error) {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return false, e
	}

	if e := vc.SetAddress(c.clients[0]); e != nil {
		return false, e
	}

	r, e := vc.Sys().SealStatus()
	if e != nil {
		return false, e
	}

	return r.Sealed, nil
}
