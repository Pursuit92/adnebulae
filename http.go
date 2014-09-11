package adnebulae

import (
	"time"
	"log"
	"path/filepath"
	"net/http"
	"errors"
	nova "github.com/Pursuit92/openstack-compute/v2"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidCreds error = errors.New("Invalid username / password")
)


type Session struct {
	Id string
	OSC *nova.ComputeClient
	Authenticated bool
	AuthFailed bool
	Last string
	Current string
	Error	error
	Renew chan bool
}

func (sess Session) Page() string {
	return filepath.Base(sess.Current)
}

func (s *Server) setupHttp() {
	r := mux.NewRouter()

	r.
		PathPrefix("/login").
		Methods("GET").
		HandlerFunc(s.loginPage)

	r.
		PathPrefix("/dashboard/{page}").
		Methods("GET").
		HandlerFunc(s.requireAuthFunc(s.handleDash))
	r.
		PathPrefix("/dashboard").
		Methods("GET").
		HandlerFunc(s.requireAuthFunc(s.handleDash))

	s.setupApi(r.PathPrefix("/api").Subrouter())

	staticFiles := filepath.Join(s.Config.Main.Files,"static")

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",http.FileServer(http.Dir(staticFiles))))
	r.PathPrefix("/assets/").HandlerFunc(
		func(w http.ResponseWriter,r *http.Request){})

	r.PathPrefix("/").HandlerFunc(s.requireAuthFunc(s.handleHome))

	s.Handler = logger(r)
}

func (s *Server) handleDash(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	vars := mux.Vars(r)
	if _,ok := vars["page"]; ok {
		sess.Current = r.URL.Path
	} else {
		http.Redirect(w,r,"/dashboard/instances",http.StatusFound)
	}
	err := s.renderTemplate(w,"dashboard",sess)
	if err != nil {
		log.Print(err)
	}

}

func (s *Server) loginPage(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	if sess.Authenticated {
		http.Redirect(w,r,sess.Current,http.StatusFound)
	} else {
		err := s.renderTemplate(w,"login",sess)
		if err != nil {
			log.Print(err)
		}
		sess.AuthFailed = false
		sess.Error = nil

	}
}

func (s *Server) handleHome(w http.ResponseWriter,r *http.Request) {
	http.Redirect(w,r,"/dashboard",http.StatusFound)
}

func (s *Server) getSession(r *http.Request,w http.ResponseWriter) (*Session) {
	s.m.Lock()
	defer s.m.Unlock()
	sess, _ := s.store.Get(r,"adnebulae")
	var id string
	if sess.IsNew  {
		id = uuidgen()
		sess.Values["id"] = id
		err := sess.Save(r,w)
		if err != nil {
			log.Print("Error: ",err)
		}
	} else {
		id = sess.Values["id"].(string)
	}
	ret,ok := s.Sessions[id]
	if ! ok {
		anSess := s.newSession()
		anSess.Id = id
		s.Sessions[id] = anSess
		ret = anSess
	}
	ret.Renew <- true
	return ret

}

func (s *Server) newSession() *Session {
	sess := &Session{}
	sess.Renew = make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(time.Duration(s.SessionTimeout) * time.Second):
				s.m.Lock()
				delete(s.Sessions,sess.Id)
				s.m.Unlock()
				return
			case <-sess.Renew:
				continue
			}
		}
	}()
	return sess
}

func (s *Server) renderTemplate(w http.ResponseWriter, tmpl string, sess *Session) error {
	return s.templates.ExecuteTemplate(w,tmpl,sess)
}

func (s *Server) requireAuthFunc(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := s.getSession(r,w)
		sess.Last = r.URL.Path
		if sess.Authenticated {
			f(w,r)
		} else {
			http.Redirect(w,r,"/login",http.StatusFound)
		}
	}
}

