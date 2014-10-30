package main

import (
	"fmt"
	"strings"

	"github.com/Pursuit92/adnebulae"
	"github.com/Pursuit92/openstack-compute/v2"
	"github.com/go-chef/chef"
)

type cookbookTable chef.CookbookListResult

func (c cookbookTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Name", "Versions"}
	for k, v := range c {
		rec := make([]string, 2)
		rec[0] = k
		vers := make([]string, len(v.Versions))
		for i, w := range v.Versions {
			vers[i] = w.Version
		}
		rec[1] = strings.Join(vers, ", ")
		tab = append(tab, rec)
	}
	return tab
}

type flavorTable []*nova.Flavor

func (i flavorTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Id", "Name", "Memory", "Disk", "CPUs"}
	for _, v := range i {
		rec := make([]string, 5)
		rec[0] = v.Id
		rec[1] = v.Name
		rec[2] = fmt.Sprintf("%v", v.Ram)
		rec[3] = fmt.Sprintf("%v", v.Disk)
		rec[4] = fmt.Sprintf("%v", v.Vcpus)
		tab = append(tab, rec)
	}
	return tab
}

type imageTable []*nova.Image

func (i imageTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Id", "Name"}
	for _, v := range i {
		rec := make([]string, 2)
		rec[0] = v.Id
		rec[1] = v.Name
		tab = append(tab, rec)
	}
	return tab
}

type netTable []*nova.Network

func (i netTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Id", "Label"}
	for _, v := range i {
		rec := make([]string, 2)
		rec[0] = v.Id
		rec[1] = v.Label
		tab = append(tab, rec)
	}
	return tab
}

type serverTable []*adnebulae.Server

func (s serverTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Id", "Name", "Status", "Networks", "Run List"}
	for _, v := range s {
		rec := make([]string, 5)
		rec[0] = v.Nova.Id
		rec[1] = v.Nova.Name
		rec[2] = v.Nova.Status
		addrs := []string{}
		for k, a := range v.Nova.Addresses {
			ips := []string{}
			for _, ip := range a {
				ips = append(ips, ip.Addr)
			}
			addrs = append(addrs, fmt.Sprintf("%s=%s", k, strings.Join(ips, ",")))
		}

		rec[3] = strings.Join(addrs, "; ")
		if v.Chef != nil {
			rec[4] = strings.Join(v.Chef.RunList, ",")
		}
		tab = append(tab, rec)
	}
	return tab
}
