package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error! %v\n", err)
	}
}

func run() error {
	storagePath, err := getStoragePath()
	if err != nil {
		return err
	}

	app := NewApp(NewHosts(getHostsFile()), NewStorage(storagePath))

	if err := app.Handle(getCmd()); err != nil {
		return err
	}

	return nil
}

func getCmd() (cmd Cmd) {
	flag.BoolVar(&cmd.Unblock, "unblock", false, "allow access to websites for a limited amount of time")
	flag.BoolVar(&cmd.Block, "block", false, "block all sites")
	flag.StringVar(&cmd.Website, "web", "", "website to ban from now on")
	flag.Parse()

	return cmd
}

func getStoragePath() (string, error) {
	path := os.Getenv("BLOCK_SITES_PATH")

	if path != "" {
		return path, nil
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path = filepath.Join(dirname, "/.blocksites")

	return path, nil
}

func getHostsFile() string {
	path := os.Getenv("BLOCK_SITES_HOSTS_FILE")
	if path != "" {
		return path
	}

	return "/etc/hosts"
}
