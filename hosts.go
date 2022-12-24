package main

import (
	"os"
)

type hostFile struct {
	hostPath string
}

func NewHostFile(hostPath string) hostFile {
	return hostFile{hostPath: hostPath}
}

func (h *hostFile) Read() ([]byte, error) {
	return os.ReadFile(h.hostPath)
}

func (h *hostFile) Write(data []byte) error {
	return os.WriteFile(h.hostPath, data, 0644)
}

type Cmd struct {
	Pause  bool
	Resume bool
}
