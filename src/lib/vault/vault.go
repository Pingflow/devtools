package vault

import (
	"encoding/json"
	"fmt"
	"github.com/Pingflow/devtools/src/lib/colors"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Pingflow/devtools/src/lib/docker"
	"github.com/Pingflow/devtools/src/lib/wait"
	"github.com/hashicorp/vault/api"
)

type IClients interface {
	IsInitialized() (bool, error)
	IsSealed() (bool, error)
	Wait() error
	Init() (*api.InitResponse, error)
	UnSeal() error
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

	sort.Strings(c)

	return New(credentialPath, c...), nil
}

func New(credentialPath string, clients ...string) IClients {
	return client{
		credentialPath: credentialPath,
		clients:        clients,
	}
}

func (c client) Wait() error {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return e
	}

	for _, vault := range c.clients {

		addr := strings.TrimPrefix(vault, "http://")

		fmt.Printf("Wait vault %v\t\t...", addr)

		if e := wait.TCP([]string{addr}, time.Second*10); e != nil {
			return e
		}

		if e := vc.SetAddress(vault); e != nil {
			return e
		}

		var ready = false
		for !ready {
			_, e := vc.Sys().Health()
			if e == nil {
				ready = true
				fmt.Printf("\rWait vault %v\t\t... %sdone%s\n", addr, colors.Green.ToString(), colors.Reset.ToString())
			} else {
				fmt.Printf("\rWait vault %v\t\t... %serror%s: %v\n", addr, colors.Red.ToString(), colors.Reset.ToString(), e)
			}
		}
	}

	return nil
}

func (c client) Init() (*api.InitResponse, error) {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
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
		addr := strings.TrimPrefix(c.clients[0], "http://")
		fmt.Printf("Init vault %v\t\t...", addr)

		r, e := vc.Sys().Init(&api.InitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		})
		if e != nil {
			fmt.Printf("\rInit vault %v\t\t... %serror%s: %v\n", addr, colors.Red.ToString(), colors.Reset.ToString(), e)
			return nil, e
		}

		b, e := json.Marshal(r)
		if e != nil {
			return nil, e
		}

		if e := ioutil.WriteFile(c.credentialPath, b, 0644); e != nil {
			return nil, e
		}

		fmt.Printf("\rInit vault %v\t\t... %sdone%s\n", addr, colors.Green.ToString(), colors.Reset.ToString())
		return r, nil
	}

	return nil, nil
}

func (c client) UnSeal() error {

	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return e
	}

	if _, e := os.Stat(c.credentialPath); os.IsNotExist(e) {
		return os.ErrNotExist
	}

	var credential api.InitResponse
	b, e := ioutil.ReadFile(c.credentialPath)
	if e != nil {
		return e
	}

	if e := json.Unmarshal(b, &credential); e != nil {
		return e
	}

	for _, v := range c.clients {

		if e := vc.SetAddress(v); e != nil {
			return e
		}

		seal, e := vc.Sys().SealStatus()
		if e != nil {
			return e
		}

		if seal.Sealed {

			addr := strings.TrimPrefix(v, "http://")
			fmt.Printf("Unseal vault %v\t\t...", addr)

			for _, k := range credential.KeysB64 {
				_, e := vc.Sys().Unseal(k)
				if e != nil {
					return e
				}
			}

			seal, e := vc.Sys().SealStatus()
			if e != nil {
				return e
			}

			if seal.Sealed {
				fmt.Printf("\rUnseal vault %v\t\t... %serror%s\n", addr, colors.Red.ToString(), colors.Reset.ToString())
			} else {
				fmt.Printf("\rUnseal vault %v\t\t... %sdone%s\n", addr, colors.Green.ToString(), colors.Reset.ToString())
			}
		}
	}

	return c.EnableUserpass(credential.RootToken)
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

func (c client) EnableUserpass(rootToken string) error {
	vc, e := api.NewClient(api.DefaultConfig())
	if e != nil {
		return e
	}

	vc.SetToken(rootToken)

	if e := vc.SetAddress(c.clients[0]); e != nil {
		return e
	}

	if e := vc.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{Type: "userpass"}); e != nil {
		fmt.Printf("\rEnable Vault Userpass auth\t\t... %serror%s\n", colors.Green.ToString(), colors.Reset.ToString())
		return e
	}

	fmt.Printf("\rEnable Vault Userpass\t\t... %sdone%s\n", colors.Green.ToString(), colors.Reset.ToString())
	return nil
}
