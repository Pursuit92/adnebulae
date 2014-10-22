package adnebulae

import (
	"fmt"
	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/marpaia/chef-golang"
	"strings"
	"sync"
)

func (an *AdNebulae) Servers() ([]*Server, error) {
	insts, err := an.ServersDetail()
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	srvs := make([]*Server, len(insts))

	wg := &sync.WaitGroup{}
	for i, v := range insts {
		wg.Add(1)

		go func(i int, v *nova.Server) {
			var chefData *chef.Node
			var chefOK bool
			if chefID, ok := v.Metadata["chef-id"]; ok {
				chefData, chefOK, err = an.GetNode(chefID)
				if err != nil {
					return
				}
			}
			if !chefOK {
				for _, w := range []string{v.Id, v.Name} {
					chefData, chefOK, err = an.GetNode(strings.ToLower(w))
					if err != nil {
						return
					}
					if chefOK {
						break
					}
				}
			}

			srvs[i] = &Server{v, chefData}
			wg.Done()
		}(i, v)
	}

	wg.Wait()
	return srvs, nil
}
