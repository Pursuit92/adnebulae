package adnebulae

func (an *AdNebulae) Delete(id string) error {
	srv, err := an.Nova.Details(id)
	if err != nil {
		srv, err = an.Nova.ServerByName(id)
		if err != nil {
			return err
		}
	}
	an.Nova.Delete(srv.Id)
	chefData, err := an.ChefData(srv)
	if err == nil {
		an.Chef.Nodes.Delete(chefData.Name)
		an.Chef.Clients.Delete(chefData.Name)
	}

	return nil
}
