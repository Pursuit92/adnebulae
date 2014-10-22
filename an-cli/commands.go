package main

import (
	//"flag"
	//"fmt"
	"github.com/Pursuit92/cli"
	//nova "github.com/Pursuit92/openstack-compute/v2"
	pt "github.com/Pursuit92/prettytables"
	//"github.com/marpaia/chef-golang"
	//"os"
	//"strings"
)

var (
	AnCliCommands = []*cli.Command{
		cli.New("list", "List VMs", listVMs).
			Subs(
			cli.New("images", "List available Images", listImages),
			cli.New("vms", "List VMs", listVMs),
			cli.New("nets", "List available Networks", listNets),
			cli.New("flavors", "List available Flavors", listFlavors),
			cli.New("roles", "List available Roles", listRoles),
			cli.New("cookbooks", "List available Cookbooks", listCookbooks)),
		cli.New("boot", "Create a new VM", nil).
			AddOpts(
			cli.IntOpt("demo", 5, "just a placeholder")),
		cli.New("update", "Update a VM", nil),
		cli.New("delete", "Delete a VM", nil)}
)

func listVMs(cmd *cli.Command) error {
	srvs, err := conn.Servers()
	if err != nil {
		return err
	}

	pt.PrintTable(serverTable(srvs))

	return nil
}

func listImages(cmd *cli.Command) error {
	imgs, err := conn.Images()
	if err != nil {
		return err
	}

	pt.PrintTable(imageTable(imgs))

	return nil
}

func listNets(cmd *cli.Command) error {
	nets, err := conn.Networks()
	if err != nil {
		return err
	}

	pt.PrintTable(netTable(nets))

	return nil
}

func listFlavors(cmd *cli.Command) error {
	flavs, err := conn.FlavorsDetail()
	if err != nil {
		return err
	}

	pt.PrintTable(flavorTable(flavs))

	return nil
}

func listRoles(cmd *cli.Command) error {
	roles, err := conn.GetRoles()
	if err != nil {
		return err
	}

	rolesTable := make([][]string, 1)
	rolesTable[0] = []string{"Role Name"}
	for k, _ := range roles {
		rolesTable = append(rolesTable, []string{k})
	}

	pt.PrintTable(pt.Table(rolesTable))

	return nil
}

func listCookbooks(cmd *cli.Command) error {
	cbooks, err := conn.GetCookbooks()
	if err != nil {
		return err
	}

	pt.PrintTable(cookbookTable(cbooks))

	return nil
}
