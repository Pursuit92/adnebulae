package main

import (
	"log"
	"io"
	"bytes"
	"github.com/gorilla/mux"
	nova "github.com/Pursuit92/openstack-compute/v2"
	"net/http"
	"encoding/json"
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

	ar.
		PathPrefix("/actions").
		Methods("POST").
		HandlerFunc(s.requireAuthFunc(s.handleActions))

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

type actionReq struct {
	Type string `json:"type,omitempty"`
	Action string `json:"action,omitempty"`
	Ids []string `json:"ids,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func alive(id string, srvs []*nova.Server) bool {
	for _,v := range srvs {
		if id == v.Id {
			return true
		}
	}
	return false
}

func anyAlive(ids []string, srvs []*nova.Server) bool {
	for _,v := range ids {
		if alive(v,srvs) {
			return true
		}
	}
	return false
}

func (s *Server) handleActions(w http.ResponseWriter, r *http.Request) {
	log.Print("Got action request!")
	req := &actionReq{}
	buf := &bytes.Buffer{}
	input := io.TeeReader(r.Body,buf)
	err := json.NewDecoder(input).Decode(req)
	if err != nil {
		http.Error(w,err.Error(),400)
	}
	sess := s.getSession(r,w)
	switch req.Type {
	case "instance":
		switch req.Action {
		case "delete":
			for _,i := range req.Ids {
				sess.OSC.Delete(i)
			}
		}
	}
}
