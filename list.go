package adnebulae

import (
	"fmt"
	"sync"

	nova "github.com/Pursuit92/openstack-compute/v2"
)

func (an *AdNebulae) Servers() ([]*Server, error) {
	insts, err := an.Nova.ServersDetail()
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	srvs := make([]*Server, len(insts))

	wg := &sync.WaitGroup{}
	for i, v := range insts {
		wg.Add(1)

		go func(i int, v *nova.Server) {
			chefData, _ := an.ChefData(v)

			srvs[i] = &Server{v, chefData}
			wg.Done()
		}(i, v)
	}

	wg.Wait()
	return srvs, nil
}

func (an *AdNebulae) Server(name string) (*Server, error) {
	srv, err := an.Nova.ServerByName(name)
	if err != nil {
		return nil, err
	}

	chefData, _ := an.ChefData(srv)
	ret := &Server{srv, chefData}
	return ret, nil

}
