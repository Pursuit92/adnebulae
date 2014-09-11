package adnebulae

import (
	"github.com/gorilla/mux"
	"log"
	nova "github.com/Pursuit92/openstack-compute/v2"
	"net/http"
	"encoding/json"
	"fmt"
)

func (s *Server) setupApi(ar *mux.Router) {
	ar.
		PathPrefix("/auth/logout").
		Methods("GET").
		HandlerFunc(s.handleLogout)
	ar.
		PathPrefix("/auth").
		Methods("POST","GET").
		HandlerFunc(s.handleAuth)

	/*
	for _,v := range []string{"/instances/{id}","/instances"} {
		ar.
			PathPrefix(v).
			Methods("GET","POST").
			HandlerFunc(s.requireAuthFunc(s.handleServers))
	}
	*/

	for _,v := range []string{"/images/{id}","/images"} {
		ar.
			PathPrefix(v).
			Methods("GET","POST").
			HandlerFunc(s.handleImages)
	}

}

func (s *Server) handleAuth(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	var err error
	if ! sess.Authenticated {
		user := r.FormValue("user")
		pass := r.FormValue("password")
		sess.OSC,_ = nova.NewClient(s.Config.Openstack.AuthUrl)
		osc := sess.OSC
		osc.PasswordAuth(user,pass)
		err = osc.Authenticate()
	} else {
		tenant := r.FormValue("tenant")
		if tenant != "" {
			osc := sess.OSC
			osc.TokenAuth(osc.Access.Token.Id)
			osc.TenantName(tenant)
			err = osc.Authenticate()
		}
	}
	if err != nil {
		sess.AuthFailed = true
		sess.Authenticated = false
		if err.Error()[:3] == "401" {
			sess.Error = ErrInvalidCreds
		} else {
			sess.Error = err
		}
		http.Redirect(w,r,"/login",http.StatusFound)
	} else {
		sess.Authenticated = true
		http.Redirect(w,r,sess.Last,http.StatusFound)
	}

}

func (s *Server) handleLogout(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	sess.Authenticated = false
	http.Redirect(w,r,"/login",http.StatusFound)
}

func (s *Server) handleServers(w http.ResponseWriter,r *http.Request) {
	log.Print("Got servers request")
	sess := s.getSession(r,w)
	vars := mux.Vars(r)
	var osresp interface{}
	var err error
	id,ok := vars["id"]
	if ok {
		osresp, err = sess.OSC.Details(id)
	} else {
		osresp, err = sess.OSC.ServersDetail()
	}

	errOut := func(err error) {
			switch err.Error()[:3] {
			case "404":
				w.WriteHeader(404)
			default:
				w.WriteHeader(500)
			}
			fmt.Fprintf(w,"Error getting data from openstack: %v\n",err)
			return
	}

	if err != nil {
		errOut(err)
		return
	}

	moreInfo := func(srv *nova.Server) error {
		flav,err := sess.OSC.FlavorDetails(srv.Flavor.Id)
		if err != nil {
			return err
		}
		img,err := sess.OSC.ImageDetails(srv.Image.Id)
		if err != nil {
			return err
		}
		srv.Image = img
		srv.Flavor = flav
		return nil
	}

	if ok {
		realResp := osresp.(*nova.Server)
		err := moreInfo(realResp)
		if err != nil {
			errOut(err)
			return
		}
	} else {
		realResp := osresp.([]*nova.Server)
		for _,v := range realResp {
			err := moreInfo(v)
			if err != nil {
				errOut(err)
				return
			}
		}
	}


	resp,err := json.MarshalIndent(osresp,"","  ")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w,"Error marshalling json: %v\n",err)
		return
	}

	w.Write(resp)
}

func (s *Server) handleImages(w http.ResponseWriter,r *http.Request) {
	log.Print("Got servers request")
	sess := s.getSession(r,w)
	vars := mux.Vars(r)
	var osresp interface{}
	var err error
	if id,ok := vars["id"]; ok {
		osresp, err = sess.OSC.ImageDetails(id)
	} else {
		osresp, err = sess.OSC.ImagesDetail()
	}

	if err != nil {
		switch err.Error()[:3] {
		case "404":
			w.WriteHeader(404)
		default:
			w.WriteHeader(500)
		}
		fmt.Fprintf(w,"Error getting data from openstack: %v\n",err)
		return
	}
	resp,err := json.MarshalIndent(osresp,"","  ")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w,"Error marshalling json: %v\n",err)
		return
	}

	w.Write(resp)
}
