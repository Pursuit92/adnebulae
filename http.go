package adnebulae

import (
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
	OSC *nova.ComputeClient
	Authenticated bool
	AuthFailed bool
	Page string
	Path string
	Error	error
}

func (s *Server) setupHttp() {
	r := mux.NewRouter()
	r.HandleFunc("/",s.requireAuthFunc(s.handleHome))
	r.HandleFunc("/login",s.loginPage)
	r.HandleFunc("/auth",s.handleAuth)

	dashroute := r.PathPrefix("/dashboard").Subrouter()
	dashroute.HandleFunc("/{page}",s.requireAuthFunc(s.handleDash))
	dashroute.HandleFunc("/",s.requireAuthFunc(s.handleDash))

	staticFiles := filepath.Join(s.Config.Main.Files,"static")

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",http.FileServer(http.Dir(staticFiles))))

	s.Handler = logger(r)
}

func (s *Server) handleAuth(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	user := r.FormValue("user")
	pass := r.FormValue("password")
	osc := sess.OSC
	osc.PasswordAuth(user,pass)
	err := osc.Authenticate()
	if err != nil {
		sess.AuthFailed = true
		if err.Error()[:3] == "401" {
			sess.Error = ErrInvalidCreds
		} else {
			sess.Error = err
		}
		http.Redirect(w,r,"/login",http.StatusFound)
	} else {
		sess.Authenticated = true
		http.Redirect(w,r,"/dashboard/",http.StatusFound)
	}
}

func (s *Server) handleDash(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	vars := mux.Vars(r)
	sess.Path = "/dashboard/"
	if page,ok := vars["page"]; ok {
		sess.Page = page
	} else {
		http.Redirect(w,r,"overview",http.StatusFound)
	}
	err := s.renderTemplate(w,"dashboard",sess)
	if err != nil {
		log.Print(err)
	}

}

func (s *Server) loginPage(w http.ResponseWriter,r *http.Request) {
	sess := s.getSession(r,w)
	if sess.Authenticated {
		http.Redirect(w,r,sess.Path+sess.Page,http.StatusFound)
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
	sess := s.getSession(r,w)
	http.Redirect(w,r,sess.Path+sess.Page,http.StatusFound)
}

func (s *Server) getSession(r *http.Request,w http.ResponseWriter) (*Session) {
	s.m.Lock()
	defer s.m.Unlock()
	sess, _ := s.store.Get(r,"adnebulae")
	var id string
	if sess.IsNew  {
		id = uuidgen()
		sess.Values["id"] = id
		anSess := s.newSession()
		s.Sessions[id] = anSess
		err := sess.Save(r,w)
		if err != nil {
			log.Print("Error: ",err)
		}
	} else {
		id = sess.Values["id"].(string)
	}
	return s.Sessions[id]

}

func (s *Server) newSession() *Session {
	sess := &Session{}
	sess.OSC,_ = nova.NewClient(s.Config.Openstack.AuthUrl)
	return sess
}

func (s *Server) renderTemplate(w http.ResponseWriter, file string, sess *Session) error {
	return s.templates.ExecuteTemplate(w,file,sess)
}

func (s *Server) requireAuthFunc(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := s.getSession(r,w)
		if sess.Authenticated {
			f(w,r)
		} else {
			http.Redirect(w,r,"/login",http.StatusFound)
		}
	}
}

