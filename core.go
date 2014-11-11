package adnebulae

import (
	"io/ioutil"
	"os"

	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/Pursuit92/chef"
)

type AdNebulae struct {
	Nova      *nova.ComputeClient
	Chef      *chef.Client
	Validator string
}

var (
	cfgVars = []string{"OS_AUTH_URL", "OS_TENANT_NAME", "OS_USERNAME", "OS_PASSWORD", "CHEF_ENDPOINT", "CHEF_USERNAME", "CHEF_KEY_FILE", "CHEF_VALIDATOR"}
)

func cfgFromEnv() map[string]string {
	ret := make(map[string]string)
	for _, v := range cfgVars {
		ret[v] = os.Getenv(v)
	}
	return ret
}

func New() (*AdNebulae, error) {
	var osClient *nova.ComputeClient
	var chefClient *chef.Client
	var err error
	cfg := cfgFromEnv()
	osClient, err = nova.NewClient(cfg["OS_AUTH_URL"])
	if err != nil {
		return nil, err

	}
	osClient.PasswordAuth(cfg["OS_USERNAME"], cfg["OS_PASSWORD"])
	osClient.TenantName(cfg["OS_TENANT_NAME"])
	err = osClient.Authenticate()
	if err != nil {
		return nil, err
	}

	chefKey, err := ioutil.ReadFile(cfg["CHEF_KEY_FILE"])
	if err != nil {
		return nil, err
	}

	chefClient, err = chef.NewClient(&chef.Config{
		Name:    cfg["CHEF_USERNAME"],
		BaseURL: cfg["CHEF_ENDPOINT"],
		Key:     string(chefKey),
		SkipSSL: true})

	if err != nil {
		return nil, err
	}

	return &AdNebulae{osClient, chefClient, cfg["CHEF_VALIDATOR"]}, err

}
