package adnebulae

import (
	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/marpaia/chef-golang"
	"os"
)

type AdNebulae struct {
	*nova.ComputeClient
	*chef.Chef
}

var (
	cfgVars = []string{"OS_AUTH_URL", "OS_TENANT_NAME", "OS_USERNAME", "OS_PASSWORD", "CHEF_CFG"}
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
	var chefClient *chef.Chef
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

	chefClient, err = chef.Connect(cfg["CHEF_CFG"])
	if err != nil {
		return nil, err
	}
	chefClient.SSLNoVerify = true

	return &AdNebulae{osClient, chefClient}, err

}
