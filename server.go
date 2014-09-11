package adnebulae

import (
	fp "path/filepath"
	"net/url"
	"sync"
	"net/http"
	"log"
	"html/template"
	"github.com/gorilla/sessions"
	sc "github.com/gorilla/securecookie"
)

type Server struct {
	Config   *Config
	Handler http.Handler
	Sessions map[string]*Session
	SessionTimeout int
	store *sessions.CookieStore
	m *sync.Mutex
	templates *template.Template
}

func NewServer(configPath string) (*Server,error) {
	an := &Server{}
	var err error
	log.Printf("Reading config file: %s",configPath)
	an.Config,err = readConfig(configPath)
	if err != nil {
		return nil,err
	}
	an.store = sessions.NewCookieStore(sc.GenerateRandomKey(32),sc.GenerateRandomKey(32))
	an.Sessions = make(map[string]*Session)
	an.m = &sync.Mutex{}
	_,err = url.Parse(an.Config.Openstack.AuthUrl)
	if err != nil {
		return nil,err
	}
	an.templates,err = readTemplates(fp.Join(an.Config.Main.Files,"template"))
	if err != nil {
		return nil,err
	}
	an.SessionTimeout = 300

	return an,nil
}

func (s *Server) Run() error {
	log.Print("Starting server")
	s.setupHttp()
	if s.Config.Main.SSL {
		key := fp.Join(s.Config.Main.ConfBase, s.Config.Main.SSLCert)
		crt := fp.Join(s.Config.Main.ConfBase, s.Config.Main.SSLKey)
		log.Printf("SSL Listening on %s",s.Config.Main.Listen)
		log.Printf("SSL key: %s",key)
		log.Printf("SSL cert: %s",crt)
		return http.ListenAndServeTLS(s.Config.Main.Listen,
			key,
			crt,
			s.Handler)
	} else {
		log.Printf("SSL Listening on %s",s.Config.Main.Listen)
		return http.ListenAndServe(s.Config.Main.Listen,
			s.Handler)
	}
}
