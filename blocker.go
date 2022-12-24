package main

import "bytes"

func Block(data []byte) ([]byte, error) {
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
