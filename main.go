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
		NewHostFile("/etc/hosts"),
		NewFileStatusManager("/tmp/hostsstatus"),
	)

	if err := app.Handle(getCmd()); err != nil {
		log.Fatalf("Error! %v\n", err)
	}
}

func getCmd() Cmd {
	var pause = flag.Bool("pause", false, "pause the execution for 1 hour")
	var resume = flag.Bool("resume", false, "resume the execution of the script")
	flag.Parse()

	return Cmd{
		Pause:  *pause,
		Resume: *resume,
	}
}
