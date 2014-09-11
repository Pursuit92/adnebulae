package adnebulae

import (
	"log"
	"os"
	"path/filepath"
	"html/template"
	"github.com/howeyc/fsnotify"
)

func readTemplates(dir string) (*template.Template,error) {
	t := template.New("top")
	//t.Delims("=%","%=")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil &&
		! info.IsDir() &&
		filepath.Base(path)[0] != '.' &&
		err == nil {
			name := filepath.Base(path)
			if name[len(name)-5:] == ".tmpl" {
				_,err := t.ParseFiles(path)
				return err
			}
			return nil
		}
		return nil
	})
	return t,err
}

func (s *Server) WatchChanges() {
	watcher,err := fsnotify.NewWatcher()
	if err != nil {
		log.Print("Error: ",err)
		return
	}
	templates := filepath.Join(s.Config.Main.Files,"template")

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() && ev.Name[len(ev.Name)-5:] == ".tmpl" && ev.Name[0] != '.' {
					t,err := readTemplates(templates)
					if err == nil {
						log.Print("Reloaded templates.")
						s.templates = t
					}
				}
			case <-watcher.Error:
				break
			}
		}
	}()

	err = watcher.Watch(templates)
	if err != nil {
		log.Print("Error: ",err)
	}
}

