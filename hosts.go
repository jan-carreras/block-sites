package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type AppStatusManager interface {
	InStatus(byte) (bool, error)
	Pause(duration time.Duration) error
	Resume() error
}

type HostReadWriter interface {
	Read() ([]byte, error)
	Write([]byte) error
}

type HostBlocker interface {
	Block([]byte) ([]byte, error)
}

type hostFile struct {
	hostPath string
}

func NewHostFile(hostPath string) HostReadWriter {
	return &hostFile{hostPath: hostPath}
}

func (h *hostFile) Read() ([]byte, error) {
	return ioutil.ReadFile(h.hostPath)
}

func (h *hostFile) Write(data []byte) error {
	return ioutil.WriteFile(h.hostPath, data, 0644)
}

type focusBlocker struct {
}

func NewFocusBlocker() HostBlocker {
	return &focusBlocker{}
}

func (f *focusBlocker) Block(data []byte) ([]byte, error) {
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

type Cmd struct {
	Pause  bool
	Resume bool
}

type App struct {
	readWriter       HostReadWriter
	hostBlocker      HostBlocker
	appStatusManager AppStatusManager
}

func NewApp(readWriter HostReadWriter, hostBlocker HostBlocker, appStatusManager AppStatusManager) *App {
	return &App{readWriter: readWriter, hostBlocker: hostBlocker, appStatusManager: appStatusManager}
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

	if isPaused, err := app.appStatusManager.InStatus(STATUS_PAUSED); err != nil {
		fmt.Println("we return an error?")
		return err
	} else {
		if isPaused {
			fmt.Println("Application in paused state. Doing nothing")
			return nil // Do nothing
		}
	}

	content, err := app.readWriter.Read()
	if err != nil {
		return err
	}

	fb := NewFocusBlocker()
	content, err = fb.Block(content)
	if err != nil {
		return err
	}

	if err := app.readWriter.Write(content); err != nil {
		return err
	}
	return nil
}

const (
	STATUS_RUNNING byte = 1
	STATUS_PAUSED       = 2
)

type appStatus struct {
	PauseUntil time.Time `json:"paused_until,omitempty"`
}

type FileStatusManager struct {
	path string
	f    io.ReadWriter
}

func NewFileStatusManager(path string) AppStatusManager {
	return &FileStatusManager{path: path}
}

func (f *FileStatusManager) InStatus(s byte) (bool, error) {
	currentStatus, err := f.currentStatus()
	if err != nil {
		return false, err
	}

	return currentStatus == s, nil
}

func (f *FileStatusManager) currentStatus() (byte, error) {
	status, err := f.read()
	if err != nil {
		return 0, err
	}

	currentStatus := STATUS_RUNNING
	if status.PauseUntil.After(time.Now()) {
		currentStatus = STATUS_PAUSED
	}
	return currentStatus, nil
}

func (f *FileStatusManager) read() (appStatus, error) {
	fd, err := os.Stat(f.path)
	if os.IsNotExist(err) || fd.Size() == 0 {
		return appStatus{}, nil
	}

	payload, err := ioutil.ReadFile(f.path)
	if err != nil {
		return appStatus{}, err
	}
	s := &appStatus{}
	if err := json.Unmarshal(payload, s); err != nil {
		return appStatus{}, err
	}

	return *s, nil
}

func (f *FileStatusManager) Pause(duration time.Duration) error {
	t := time.Now().Add(duration)
	s := appStatus{
		PauseUntil: t,
	}

	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.path, payload, 0644)
}

func (f *FileStatusManager) Resume() error {
	s, err := f.read()
	if err != nil {
		return err
	}

	s.PauseUntil = time.Time{}

	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.path, payload, 0644)
}