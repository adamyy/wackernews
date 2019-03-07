package main

import (
	"log"
	"os"
	"os/user"
	"path"
)

func baseDir() (string, error) {
	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, ".termhn")
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if mkdirErr := os.Mkdir(dir, 0700); mkdirErr != nil {
				return "", mkdirErr
			}
		} else {
			return "", err
		}
	}
	return dir, nil
}

func main() {
	dir, err := baseDir()
	if err != nil {
		log.Fatalf("Could not setup directory: %v", err)
	}
	logfile := path.Join(dir, "hackernews.log")
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	app, err := NewApp()
	if err != nil {
		log.Panicln(err)
	}
	defer app.Close()

	if err := app.Init(); err != nil {
		log.Println(err)
	}
}
