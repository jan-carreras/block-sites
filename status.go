package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	StatusRunning byte = 1
	StatusPaused       = 2
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

	currentStatus := StatusRunning
	if status.PauseUntil.After(time.Now()) {
		currentStatus = StatusPaused
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