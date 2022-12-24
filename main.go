package main

import (
	"flag"
	"log"
)

// Be able to --pause (for 1h) and --resume the script functions
// In Paused state it will prevent the script from doing anything
// "Resume" with return the script into Running state

func main() {
	app := NewApp(
		NewHosts("/etc/hosts"),
		NewStorage("/tmp/hostsstatus"),
	)

	if err := app.Handle(getCmd()); err != nil {
		log.Fatalf("Error! %v\n", err)
	}
}

func getCmd() (cmd Cmd) {
	flag.BoolVar(&cmd.Unblock, "unblock", false, "allow access to websites for a limited amount of time")
	flag.BoolVar(&cmd.Block, "block", false, "block all sites")
	flag.StringVar(&cmd.Website, "web", "", "website to ban from now on")
	flag.Parse()

	return cmd
}
