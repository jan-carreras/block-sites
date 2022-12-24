package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestApp_Handle(t *testing.T) {
	tests := []struct {
		name      string
		havingCmd Cmd
		having    string
		expects   string
	}{
		{
			name:      "Base Case",
			havingCmd: Cmd{Resume: false},
			having:    testdata(t, "basecase.input.txt"),
			expects:   testdata(t, "basecase.output.txt"),
		},
		{
			name:      "Empty line should not be commented",
			havingCmd: Cmd{Pause: false},
			having:    testdata(t, "emptyline.input.txt"),
			expects:   testdata(t, "emptyline.output.txt"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := ioutil.TempFile("/tmp", "hoststest")
			require.NoError(t, err)
			defer os.Remove(file.Name())

			status, err := ioutil.TempFile("/tmp", "status")
			require.NoError(t, err)
			defer os.Remove(status.Name())

			err = ioutil.WriteFile(file.Name(), []byte(test.having), 0666)
			require.NoError(t, err)

			app := NewApp(
				NewHostFile(file.Name()),
				NewFileStatusManager(status.Name()),
			)

			err = app.Handle(test.havingCmd)
			require.NoError(t, err)

			data, err := ioutil.ReadFile(file.Name())
			require.NoError(t, err)
			assert.Equal(t, test.expects, string(data))
		})
	}
}

func testdata(t *testing.T, name string) string {
	f, err := os.Open(filepath.Join("testdata", name))
	require.NoError(t, err)

	content, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	return string(content)
}
