package adnebulae

import (
	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/Pursuit92/chef"
)

type Server struct {
	Nova *nova.Server
	Chef *chef.Node
}
