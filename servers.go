package adnebulae

import (
	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/go-chef/chef"
)

type Server struct {
	Nova *nova.Server
	Chef *chef.Node
}
