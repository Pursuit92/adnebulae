package adnebulae

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/Pursuit92/openstack-compute/v2"
)

const chefClientTmpl string = `chef_server_url "{{ .Server }}"
validation_client_name "chef-validator"
node_name "{{ .Id }}"
`
const chefBootScript string = `#!/bin/bash
curl -L https://www.opscode.com/chef/install.sh | bash
mkdir /etc/chef
mv /tmp/client.rb /etc/chef
mv /tmp/validation.pem /etc/chef
mv /tmp/first-boot.json /etc/chef
chef-client -j /etc/chef/first-boot.json
`

func (an *AdNebulae) Create(name, img, flav, key, net string, enroll bool, runList []string) (*Server, error) {
	if len(runList) == 1 && runList[0] == "" {
		runList = nil
	}
	srv := nova.NewServer()
	if name == "" {
		return nil, fmt.Errorf("No name specified")
	}
	if flav == "" {
		return nil, fmt.Errorf("No flavor specified")
	}
	if img == "" {
		return nil, fmt.Errorf("No image specified")
	}
	if net == "" {
		return nil, fmt.Errorf("No network specified")
	}
	if key != "" {
		srv.KeyName = key
	}

	srv.Name = name
	srv.Image = &nova.Image{Name: img}
	srv.Flavor = &nova.Flavor{Name: flav}
	srv.NetNames = []string{net}
	if enroll {
		id := uuidgen()
		firstBoot, _ := json.Marshal(map[string][]string{"run_list": runList})
		validator, err := ioutil.ReadFile(an.Validator)
		if err != nil {
			return nil, err
		}
		clientTmpl, err := template.New("client").Parse(chefClientTmpl)
		if err != nil {
			panic(err)
		}

		clientConfig := map[string]string{"Id": id, "Server": an.Chef.BaseURL.String()}
		var clientBuf bytes.Buffer
		err = clientTmpl.Execute(&clientBuf, clientConfig)

		srv.Metadata = map[string]string{"chef-id": id}
		srv.Personality = []nova.Personality{
			nova.Personality{
				Path:     "/tmp/first-boot.json",
				Contents: base64.StdEncoding.EncodeToString(append(firstBoot, byte('\n')))},
			nova.Personality{
				Path:     "/tmp/client.rb",
				Contents: base64.StdEncoding.EncodeToString(clientBuf.Bytes())},
			nova.Personality{
				Path:     "/tmp/validation.pem",
				Contents: base64.StdEncoding.EncodeToString(validator)}}
		srv.UserData = base64.StdEncoding.EncodeToString([]byte(chefBootScript))
	}

	srv, err := an.Nova.Create(srv)
	return &Server{srv, nil}, err
}
