package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

type App struct {
	hostFilePath string
	storage      *Storage
}

func NewApp(
	hostFilePath string,
	storage *Storage,
) *App {
	return &App{
		hostFilePath: hostFilePath,
		storage:      storage,
	}
}

type Cmd struct {
	Pause  bool
	Resume bool
}

func (app *App) Handle(cmd Cmd) error {
	if cmd.Pause {
		err := app.storage.Pause(time.Hour)
		if err != nil {
			return err
		}
	}

	if cmd.Resume {
		if err := app.storage.Resume(); err != nil {
			return err
		}
	}

	if isPaused, err := app.storage.IsStatus(StatusPaused); err != nil {
		return err
	} else if isPaused {
		fmt.Println("Application in paused state. Doing nothing")
		return nil // Do nothing
	}

	content, err := os.ReadFile(app.hostFilePath)
	if err != nil {
		return err
	}

	content, err = block(content)
	if err != nil {
		return err
	}

	if err := os.WriteFile(app.hostFilePath, content, 0644); err != nil {
		return err
	}

	return nil
}

func block(data []byte) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))

	block := false

	for index, line := range lines {
		if bytes.Equal(line, []byte("#BLOCKME")) {
			block = true
			continue
		}
		if bytes.Equal(line, []byte("#/BLOCKME")) {
			block = false
			continue
		}
		if !block { // We don't want to block, ignoring
			continue
		}

		if len(line) == 0 { // Empty line, ignoring
			continue
		}

		if line[0] != byte('#') { // Already blocked, ignoring
			continue
		}

		lines[index] = line[1:]
	}
	return bytes.Join(lines, []byte("\n")), nil
}
