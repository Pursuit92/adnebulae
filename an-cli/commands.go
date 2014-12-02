package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Pursuit92/chef"
	"github.com/Pursuit92/cli"
	"github.com/Pursuit92/openstack-compute/v2"
	pt "github.com/Pursuit92/prettytables"
)

var (
	AnCliCommands = []*cli.Command{
		cli.New("list", "List VMs", listVMs).
			Subs(
			cli.New("environments", "List available environments", listEnvironments),
			cli.New("images", "List available Images", listImages),
			cli.New("vms", "List VMs", listVMs),
			cli.New("nets", "List available Networks", listNets),
			cli.New("flavors", "List available Flavors", listFlavors),
			cli.New("roles", "List available Roles", listRoles),
			cli.New("databags", "List databags", listDBs),
			cli.New("cookbooks", "List available Cookbooks", listCookbooks)),
		cli.New("boot", "Create a new VM", bootVM).
			AddOpts(
			cli.StringOpt("net", "", "Network for the new VM"),
			cli.StringOpt("flavor", "", "New VM size"),
			cli.StringOpt("runlist", "", "Chef run-list"),
			cli.StringOpt("key-name", "", "Keypair to use for the new instance"),
			cli.BoolOpt("chef", true, "Enroll VM in chef"),
			cli.StringOpt("image", "", "Image to boot")),
		cli.New("update", "Update a VM", updateVMs).
			AddOpts(
			cli.StringOpt("environment", "", "Chef environment"),
			cli.StringOpt("runlist", "", "Chef run-list")),
		cli.New("show", "Show more info", showVM).
			Subs(
			cli.New("environment", "Show environment info", showEnvironment),
			cli.New("cookbook", "Show cookbook info", showCookbook),
			cli.New("vm", "Show info about VM", showVM),
			cli.New("databag", "Show info about DataBags", showDB),
			cli.New("role", "Show info about Roles", showRole)),
		cli.New("delete", "Delete a VM", deleteVMs),
		cli.New("edit", "Edit Chef data", editNode).
			Subs(
			cli.New("databag", "Edit databag", editDB),
			cli.New("node", "Edit node", editNode),
			cli.New("role", "Edit role", editRole))}
)

func listEnvironments(cmd *cli.Command) error {
	start()
	environments, err := conn.Chef.Environments.List()
	if err != nil {
		return err
	}

	environmentsTable := make([][]string, 1)
	environmentsTable[0] = []string{"Role Name"}
	for k, _ := range *environments {
		environmentsTable = append(environmentsTable, []string{k})
	}

	pt.PrintTable(pt.Table(environmentsTable))
	return nil
}

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
	cbooks, err := conn.Chef.Cookbooks.ListAvailableVersions("4")
	if err != nil {
		return err
	}

	pt.PrintTable(cookbooksTable(cbooks))

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
	log.Print(runList)
	_, err := conn.Create(cmd.Args[0], img, flav, key, net, enroll, strings.Split(runList, ","))
	if err != nil {
		return err
	}

	return showVM(cmd)
}

func updateVMs(cmd *cli.Command) error {
	start()
	runList, _ := cmd.StringOpt("runlist")
	environment, _ := cmd.StringOpt("environment")
	srvs, err := conn.Servers()
	if err != nil {
		return err
	}

	if len(cmd.Args) < 1 {
		cmd.PrintHelp(nil)
		return nil
	}
	for _, v := range cmd.Args {
		arg := strings.ToLower(v)
		for _, s := range srvs {
			if arg == strings.ToLower(s.Nova.Name) ||
				arg == strings.ToLower(s.Nova.Id) {
				node, err := conn.ChefData(s.Nova)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					break
				}
				splitList := strings.Split(runList, ",")
				if len(splitList) > 0 && splitList[0] == "" {
					splitList = []string{}
				}
				//log.Print(splitList)

				if len(splitList) != 0 {
					node.RunList = splitList
				}
				if environment != "" {
					node.Environment = environment
				}
				_, err = conn.Chef.Nodes.Put(*node)
				if err != nil {
					return err
				}
				return showVM(cmd)

			}
		}
	}
	return nil
}

func showVM(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 1 {
		err := fmt.Errorf("Not enough arguments")
		cmd.PrintHelp(err)
		return err
	}
	srv, err := conn.Server(cmd.Args[0])
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		srv.Nova.Image, err = conn.Nova.ImageDetails(srv.Nova.Image.Id)
		if err != nil {
			srv.Nova.Image = &nova.Image{Name: "-"}
		}
		wg.Done()
	}()
	go func() {
		srv.Nova.Flavor, err = conn.Nova.FlavorDetails(srv.Nova.Flavor.Id)
		if err != nil {
			srv.Nova.Flavor = &nova.Flavor{Name: "-"}
		}
		wg.Done()
	}()

	wg.Wait()

	pt.PrintTable(vmTable(*srv))
	return nil
}

func showRole(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Not enough arguments")
	}
	role := strings.ToLower(cmd.Args[0])
	roles, err := conn.Chef.Roles.List()
	if err != nil {
		return err
	}
	for k, _ := range *roles {
		if strings.ToLower(k) == role {
			role, err := conn.Chef.Roles.Get(k)
			if err != nil {
				return err
			}

			pt.PrintTable(roleTable(*role))
		}
	}

	return nil
}

func showCookbook(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 2 {
		cookbook := strings.ToLower(cmd.Args[0])
		cookbooks, err := conn.Chef.Cookbooks.ListAvailableVersions("4")
		if err != nil {
			return err
		}
		for k, v := range cookbooks {
			if strings.ToLower(k) == cookbook {
				pt.PrintTable(cookbooksTable(map[string]chef.CookbookVersions{k: v}))
			}
		}
		return nil

	}
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Not enough arguments")
	}
	cookbook := strings.ToLower(cmd.Args[0])
	version := strings.ToLower(cmd.Args[1])
	cookbooks, err := conn.Chef.Cookbooks.List()
	if err != nil {
		return err
	}
	for k, _ := range cookbooks {
		if strings.ToLower(k) == cookbook {
			cookbook, err := conn.Chef.Cookbooks.GetVersion(k, version)
			if err != nil {
				return err
			}

			pt.PrintTable(cookbookTable(cookbook))
		}
	}

	return nil
}

func showEnvironment(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Not enough arguments")
	}
	environment := strings.ToLower(cmd.Args[0])
	environments, err := conn.Chef.Environments.List()
	if err != nil {
		return err
	}
	for k, _ := range *environments {
		if strings.ToLower(k) == environment {
			environment, err := conn.Chef.Environments.Get(k)
			if err != nil {
				return err
			}

			pt.PrintTable(environmentTable(*environment))
		}
	}

	return nil
}

func editNode(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Not enough arguments")
	}
	node, err := conn.Server(cmd.Args[0])
	if err != err {
		return err
	}
	if node.Chef == nil {
		return fmt.Errorf("Not a chef node")
	}

	err = edit(node.Chef.Name, func() ([]byte, error) {
		cont, err := json.MarshalIndent(node.Chef.NormalAttributes, "", "  ")
		if err != nil {
			return nil, err
		}

		return cont, nil
	}, func(cont []byte) error {
		err = json.Unmarshal(cont, &node.Chef.NormalAttributes)
		if err != nil {
			return err
		}

		_, err = conn.Chef.Nodes.Put(*node.Chef)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Done editing %s!\n", cmd.Args[0])

	return nil
}

func editRole(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Not enough arguments")
	}
	roleName := strings.ToLower(cmd.Args[0])
	roles, err := conn.Chef.Roles.List()
	var role *chef.Role
	if err != nil {
		return err
	}
	for k, _ := range *roles {
		if strings.ToLower(k) == roleName {
			role, err = conn.Chef.Roles.Get(k)
			if err != nil {
				return err
			}

		}
	}
	if role == nil {
		return fmt.Errorf("Role not found!")
	}

	err = edit(role.Name, func() ([]byte, error) {
		cont, err := json.MarshalIndent(role, "", "  ")
		if err != nil {
			return nil, err
		}

		return cont, nil
	}, func(cont []byte) error {
		err = json.Unmarshal(cont, role)
		if err != nil {
			return err
		}

		_, err = conn.Chef.Roles.Put(role)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Done editing %s!\n", cmd.Args[0])

	return nil

}

func showDB(cmd *cli.Command) error {
	start()
	switch len(cmd.Args) {
	case 1:
		dbList, err := conn.Chef.DataBags.ListItems(cmd.Args[0])
		if err != nil {
			return err
		}
		pt.PrintTable(dbListTable(*dbList))

	case 2:
		dbInfo, err := conn.Chef.DataBags.GetItem(cmd.Args[0], cmd.Args[1])
		if err != nil {
			return err
		}
		pt.PrintTable(dbTable((interface{})(dbInfo).(map[string]interface{})))

	default:
		return fmt.Errorf("Incorrect number of arguments")
	}
	return nil
}

func listDBs(cmd *cli.Command) error {
	start()
	switch len(cmd.Args) {
	case 0:
		dbList, err := conn.Chef.DataBags.List()
		if err != nil {
			return err
		}
		pt.PrintTable(dbListTable(*dbList))

	default:
		return fmt.Errorf("Incorrect number of arguments")
	}
	return nil
}

func editDB(cmd *cli.Command) error {
	start()
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough arguments")
	}
	dbItem, err := conn.Chef.DataBags.GetItem(cmd.Args[0], cmd.Args[1])
	if err != nil {
		return err
	}

	err = edit(cmd.Args[0]+"-"+cmd.Args[1], func() ([]byte, error) {
		cont, err := json.MarshalIndent(dbItem, "", "  ")
		if err != nil {
			return nil, err
		}

		return cont, nil
	}, func(cont []byte) error {
		err = json.Unmarshal(cont, &dbItem)
		if err != nil {
			return err
		}

		err = conn.Chef.DataBags.UpdateItem(cmd.Args[0], cmd.Args[1], dbItem)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Done editing %s %s!\n", cmd.Args[0], cmd.Args[1])

	return nil

}
