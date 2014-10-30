package main

import (
	//"flag"
	//"fmt"
	"fmt"
	"strings"

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
		cli.New("boot", "Create a new VM", bootVM).
			AddOpts(
			cli.StringOpt("net", "", "Network for the new VM"),
			cli.StringOpt("flavor", "", "New VM size"),
			cli.StringOpt("runlist", "", "Chef run-list"),
			cli.StringOpt("key-name", "", "Keypair to use for the new instance"),
			cli.BoolOpt("chef", true, "Enroll VM in chef"),
			cli.StringOpt("image", "", "Image to boot")),
		//cli.New("update", "Update a VM", nil),
		cli.New("delete", "Delete a VM", deleteVMs)}
)

func listVMs(cmd *cli.Command) error {
	start()
	srvs, err := conn.Servers()
	if err != nil {
		return err
	}

	pt.PrintTable(serverTable(srvs))

	return nil
}

func listImages(cmd *cli.Command) error {
	start()
	imgs, err := conn.Nova.Images()
	if err != nil {
		return err
	}

	pt.PrintTable(imageTable(imgs))

	return nil
}

func listNets(cmd *cli.Command) error {
	start()
	nets, err := conn.Nova.Networks()
	if err != nil {
		return err
	}

	pt.PrintTable(netTable(nets))

	return nil
}

func listFlavors(cmd *cli.Command) error {
	start()
	flavs, err := conn.Nova.FlavorsDetail()
	if err != nil {
		return err
	}

	pt.PrintTable(flavorTable(flavs))

	return nil
}

func listRoles(cmd *cli.Command) error {
	start()
	roles, err := conn.Chef.Roles.List()
	if err != nil {
		return err
	}

	rolesTable := make([][]string, 1)
	rolesTable[0] = []string{"Role Name"}
	for k, _ := range *roles {
		rolesTable = append(rolesTable, []string{k})
	}

	pt.PrintTable(pt.Table(rolesTable))

	return nil
}

func listCookbooks(cmd *cli.Command) error {
	start()
	cbooks, err := conn.Chef.Cookbooks.List()
	if err != nil {
		return err
	}

	pt.PrintTable(cookbookTable(cbooks))

	return nil
}

func deleteVMs(cmd *cli.Command) error {
	start()
	for _, v := range cmd.Args {
		err := conn.Delete(v)
		if err != nil {
			fmt.Println("Delete Error:", err)
		}

	}
	return nil
}

func bootVM(cmd *cli.Command) error {
	start()
	flav, _ := cmd.StringOpt("flavor")
	img, _ := cmd.StringOpt("image")
	net, _ := cmd.StringOpt("net")
	key, _ := cmd.StringOpt("key-name")
	enroll, _ := cmd.BoolOpt("chef")
	runList, _ := cmd.StringOpt("runlist")
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Must supply instance name")
	}
	_, err := conn.Create(cmd.Args[0], img, flav, key, net, enroll, strings.Split(runList, ","))
	return err
}
