package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	StatusRunning = "running"
	StatusPaused  = "paused"
)

const databaseFile = "db.json"

type Website struct {
	URL string
}

type appStatus struct {
	PauseUntil time.Time `json:"paused_until,omitempty"`
	Websites   []Website `json:"websites"`
}

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{path: filepath.Join(path, databaseFile)}
}

func (s *Storage) Websites() ([]Website, error) {
	status, err := s.read()
	if err != nil {
		return nil, err
	}

	return status.Websites, nil
}

func (s *Storage) IsStatus(status string) (bool, error) {
	currentStatus, err := s.currentStatus()
	if err != nil {
		return false, err
	}

	return currentStatus == status, nil
}

func (s *Storage) Pause(duration time.Duration) error {
	d, err := s.read()
	if err != nil {
		return err
	}

	d.PauseUntil = time.Now().Add(duration)

	return s.write(d)
}

func (s *Storage) Resume() error {
	d, err := s.read()
	if err != nil {
		return err
	}

	d.PauseUntil = time.Time{}

	return s.write(d)
}

func (s *Storage) BanWebsite(website string) error {
	data, err := s.read()
	if err != nil {
		return err
	}

	for _, w := range data.Websites {
		if w.URL == website {
			return nil
		}
	}

	data.Websites = append(data.Websites, Website{URL: website})

	return s.write(data)
}

func (s *Storage) currentStatus() (string, error) {
	status, err := s.read()
	if err != nil {
		return "", err
	}

	currentStatus := StatusRunning
	if status.PauseUntil.After(time.Now()) {
		currentStatus = StatusPaused
	}

	return currentStatus, nil
}

func (s *Storage) write(status appStatus) error {
	s.lazyCreateDirectory()

	payload, err := json.Marshal(status)
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, payload, 0644)
}

func (s *Storage) read() (appStatus, error) {
	s.lazyCreateDirectory()

	fd, err := os.Stat(s.path)
	if errors.Is(err, fs.ErrNotExist) || (fd != nil && fd.Size() == 0) {
		return appStatus{}, nil
	}

	payload, err := os.ReadFile(s.path)
	if err != nil {
		return appStatus{}, err
	}

	d := &appStatus{}
	if err := json.Unmarshal(payload, d); err != nil {
		return appStatus{}, err
	}

	return *d, nil
}

func (s *Storage) lazyCreateDirectory() {
	_ = os.MkdirAll(filepath.Dir(s.path), 0700)
}
