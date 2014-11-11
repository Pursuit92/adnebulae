package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Pursuit92/adnebulae"
	"github.com/Pursuit92/openstack-compute/v2"
	"github.com/Pursuit92/chef"
)

type cookbooksTable chef.CookbookListResult

func (c cookbooksTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Name", "Versions"}
	for k, v := range c {
		rec := make([]string, 2)
		rec[0] = k
		vers := make([]string, len(v.Versions))
		for i, w := range v.Versions {
			vers[i] = w.Version
		}
		rec[1] = strings.Join(vers, "\n")
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

type vmTable adnebulae.Server

func (vm vmTable) Table() [][]string {
	tab := make([][]string, 1)
	tab[0] = []string{"Property", "Value"}
	addProp := func(p, v string) {
		tab = append(tab, []string{p, v})
	}
	var addrs []string
	for k, a := range vm.Nova.Addresses {
		ips := []string{}
		for _, ip := range a {
			ips = append(ips, ip.Addr)
		}
		addrs = append(addrs, fmt.Sprintf("%s=%s", k, strings.Join(ips, ",")))
	}
	addProp("Name", vm.Nova.Name)
	addProp("Network", strings.Join(addrs, ";"))
	addProp("Flavor", vm.Nova.Flavor.Name)
	addProp("Image", vm.Nova.Image.Name)
	addProp("Created", vm.Nova.Created)
	addProp("Updated", vm.Nova.Updated)
	addProp("Status", vm.Nova.Status)
	addProp("Id", vm.Nova.Id)
	addProp("TenantId", vm.Nova.TenantId)
	if vm.Chef.Name != "" {
		addProp("Environment", vm.Chef.Environment)
		addProp("RunList", strings.Join(vm.Chef.RunList, ",\n"))
		if _, ok := vm.Chef.AutomaticAttributes["platform"]; ok {
			addProp("Platform", vm.Chef.AutomaticAttributes["platform"].(string)+" "+vm.Chef.AutomaticAttributes["platform_version"].(string))
		}
		normal, _ := json.MarshalIndent(vm.Chef.NormalAttributes, "", "  ")
		addProp("Attributes", string(normal))
	}
	return tab
}

type environmentTable chef.Environment

func (rt environmentTable) Table() [][]string {
	tab := make([][]string, 0)
	addProp := func(p, v string) {
		tab = append(tab, []string{p, v})
	}
	addProp("Property", "Value")
	addProp("Name", rt.Name)
	addProp("Desc", rt.Description)
	def, _ := json.MarshalIndent(rt.DefaultAttributes, "", "  ")
	override, _ := json.MarshalIndent(rt.OverrideAttributes, "", "  ")
	addProp("Default", string(def))
	addProp("Override", string(override))
	books, _ := json.MarshalIndent(rt.CookbookVersions, "", "  ")
	addProp("Cookbooks", string(books))
	return tab
}

type roleTable chef.Role

func (rt roleTable) Table() [][]string {
	tab := make([][]string, 0)
	addProp := func(p, v string) {
		tab = append(tab, []string{p, v})
	}
	addProp("Property", "Value")
	addProp("Name", rt.Name)
	addProp("Desc", rt.Description)
	addProp("RunList", strings.Join(rt.RunList, ",\n"))
	def, _ := json.MarshalIndent(rt.DefaultAttributes, "", "  ")
	override, _ := json.MarshalIndent(rt.OverrideAttributes, "", "  ")
	addProp("Default", string(def))
	addProp("Override", string(override))
	return tab
}

type cookbookTable chef.Cookbook

func collectName(is []chef.CookbookItem) string {
	var names string
	for _, v := range is {
		if len(names) > 0 {
			names = fmt.Sprintf("%s\n%s", names, v.Name)
		} else {
			names = v.Name
		}
	}
	return names
}

func (ct cookbookTable) Table() [][]string {
	tab := make([][]string, 0)
	addProp := func(p, v string) {
		tab = append(tab, []string{p, v})
	}
	addProp("Property", "Value")
	addProp("Name", ct.Name)
	addProp("Recipes", collectName(ct.Recipes))
	return tab
}
