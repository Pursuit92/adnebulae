package adnebulae

import (
	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/marpaia/chef-golang"
)

type Server struct {
	Nova *nova.Server
	Chef *chef.Node
}
