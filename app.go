package main

import (
	"fmt"
	"time"
)

type AppStatusManager interface {
	InStatus(byte) (bool, error)
	Pause(duration time.Duration) error
	Resume() error
}

type App struct {
	hostFile         hostFile
	appStatusManager AppStatusManager
}

func NewApp(
	hostFile hostFile,
	appStatusManager AppStatusManager,
) *App {
	return &App{
		hostFile:         hostFile,
		appStatusManager: appStatusManager,
	}
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

	content, err := app.hostFile.Read()
	if err != nil {
		return err
	}

	content, err = Block(content)
	if err != nil {
		return err
	}

	if err := app.hostFile.Write(content); err != nil {
		return err
	}

	return nil
}
