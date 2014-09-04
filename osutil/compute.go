package osutil

import (
	nova "github.com/Pursuit92/openstack/compute/v2"
	"errors"
)

var (
	ErrServerNotFound error = errors.New("Server not found.")
	ErrImageNotFound error = errors.New("Image not found.")
	ErrFlavorNotFound error = errors.New("Flavor not found.")
)

type NovaClient struct {
	*nova.ComputeClient
}

func NewNovaClient(authUrl string) (NovaClient,error) {
	cl,err := nova.NewClient(authUrl)
	return NovaClient{cl},err
}

func (nc NovaClient) ServerByName(name string) (*nova.Server, error) {
	srvs,err := nc.Servers()
	if err != nil {
		return nil,err
	}
	for _,v := range srvs {
		if v.Name == name {
			srv,err := nc.Details(v)
			if err == nil {
				return srv,nil
			}
			return nil,err
		}
	}
	return nil,ErrServerNotFound
}


func (nc NovaClient) ImageByName(name string) (*nova.Image,error) {
	imgs,err := nc.Images()
	if err != nil {
		return nil,err
	}
	for _,v := range imgs {
		if v.Name == name {
			img,err := nc.ImageDetails(v)
			if err == nil {
				return img,nil
			}
			return nil,err
		}
	}
	return nil,ErrImageNotFound
}

func (nc NovaClient) FlavorByName(name string) (*nova.Flavor,error) {
	flavs,err := nc.Flavors()
	if err != nil {
		return nil,err
	}
	for _,v := range flavs {
		if v.Name == name {
			flav,err := nc.FlavorDetails(v)
			if err == nil {
				return flav,nil
			}
			return nil,err
		}
	}
	return nil,ErrFlavorNotFound
}

