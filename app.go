package main

import (
	"fmt"
	"os"
	"time"
)

type AppStatusManager interface {
	InStatus(byte) (bool, error)
	Pause(duration time.Duration) error
	Resume() error
}

type App struct {
	hostFilePath     string
	appStatusManager AppStatusManager
}

func NewApp(
	hostFilePath string,
	appStatusManager AppStatusManager,
) *App {
	return &App{
		hostFilePath:     hostFilePath,
		appStatusManager: appStatusManager,
	}
}

type Cmd struct {
	Pause  bool
	Resume bool
}

func (app *App) Handle(cmd Cmd) error {
	if cmd.Pause {
		err := app.appStatusManager.Pause(time.Hour)
		if err != nil {
			return err
		}
	}

	if cmd.Resume {
		if err := app.appStatusManager.Resume(); err != nil {
			return err
		}
	}

	if isPaused, err := app.appStatusManager.InStatus(StatusPaused); err != nil {
		return err
	} else if isPaused {
		fmt.Println("Application in paused state. Doing nothing")
		return nil // Do nothing
	}

	content, err := os.ReadFile(app.hostFilePath)
	if err != nil {
		return err
	}

	content, err = Block(content)
	if err != nil {
		return err
	}

	if err := os.WriteFile(app.hostFilePath, content, 0644); err != nil {
		return err
	}

	return nil
}
