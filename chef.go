package adnebulae

import (
	"strings"

	"github.com/Pursuit92/openstack-compute/v2"
	"github.com/Pursuit92/chef"
)

func (an *AdNebulae) ChefData(srv *nova.Server) (*chef.Node, error) {
	var chefData chef.Node
	var err error
	chefID, ok := srv.Metadata["chef-id"]
	if ok {
		chefData, err = an.Chef.Nodes.Get(chefID)
	} else {
		for _, w := range []string{srv.Id, srv.Name} {
			chefData, err = an.Chef.Nodes.Get(strings.ToLower(w))
			if err == nil {
				break
			}
		}
	}

	return &chefData, err
}
