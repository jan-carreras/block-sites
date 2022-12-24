package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	start = "#BLOCKME\n"
	end   = "#/BLOCKME\n"
)

// Hosts reads and writes the /etc/hosts file adding or removing
// the blocking rules
type Hosts struct {
	path string
}

func NewHosts(path string) *Hosts {
	return &Hosts{path: path}
}

func (h *Hosts) Block(websites []Website) error {
	if err := h.lazyInitialize(); err != nil {
		return err
	}

	f, err := h.read()
	if err != nil {
		return err
	}

	startIdx := bytes.Index(f, []byte(start))
	endIdx := bytes.Index(f, []byte(end))

	buf := bytes.Buffer{}
	buf.Write(f[:startIdx+len(start)])
	for _, website := range websites {
		buf.WriteString(fmt.Sprintf("0.0.0.0 %s\n", website.URL))
	}
	buf.Write(f[endIdx:])

	return h.write(buf.Bytes())
}

func (h *Hosts) Unblock() error {
	if err := h.lazyInitialize(); err != nil {
		return err
	}

	f, err := h.read()
	if err != nil {
		return err
	}

	startIdx := bytes.Index(f, []byte(start))
	endIdx := bytes.Index(f, []byte(end))

	buf := bytes.Buffer{}
	buf.Write(f[:startIdx+len(start)])
	buf.Write(f[:endIdx+len(end)])

	return h.write(buf.Bytes())
}

func (h *Hosts) lazyInitialize() error {
	f, err := h.read()
	if err != nil {
		return err
	}

	startIdx := bytes.Index(f, []byte(start))
	endIdx := bytes.Index(f, []byte(end))
	if startIdx == -1 && endIdx == -1 {
		buf := bytes.NewBuffer(f)
		buf.WriteString("\n" + start + end)
		return h.write(buf.Bytes())
	}

	if startIdx >= endIdx {
		return fmt.Errorf("%q file format messed up (startIdx=%d, endIdx=%d). Fix it manually before I make it worst", h.path, startIdx, endIdx)
	}

	return nil
}

func (h *Hosts) read() ([]byte, error) {
	content, err := os.ReadFile(h.path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (h *Hosts) write(content []byte) error {
	return os.WriteFile(h.path, content, 0644)
}
