package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func edit(name string, setup func() ([]byte, error), finish func([]byte) error) error {
	tmpFile := filepath.Join(os.TempDir(), name+".json")
	var file *os.File
	var cont []byte
	var modtime time.Time
	info, err := os.Stat(tmpFile)
	if err != nil {
		fmt.Println("Creating new file", tmpFile)
		file, err = os.Create(tmpFile)
		if err != nil {
			return err
		}

		cont, err := setup()
		if err != nil {
			return err
		}

		_, err = file.Write(cont)
		if err != nil {
			return err
		}
		file.Close()
		info, _ := os.Stat(tmpFile)
		modtime = info.ModTime()
	} else {
		modtime = info.ModTime()
		file, err = os.Open(tmpFile)
		if err != nil {
			return err
		}

	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	fullpath, err := exec.LookPath(editor)
	if err == nil {
		editor = fullpath
	}
	proc, err := os.StartProcess(editor, []string{editor, tmpFile}, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err != nil {
		return err
	}
	state, err := proc.Wait()
	if err != nil || !state.Success() {
		return fmt.Errorf("Unsuccessful editor")
	}

	info, _ = os.Stat(tmpFile)
	if !modtime.Equal(info.ModTime()) {

		file, err = os.Open(tmpFile)
		if err != nil {
			return err
		}

		cont, err = ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		err = finish(cont)
		if err != nil {
			return err
		}
		return os.Remove(tmpFile)
	}
	fmt.Printf("File unmodified.\n")
	return nil

}
