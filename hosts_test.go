package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
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
			having: `# LOCALHOST
127.0.0.1       localhost

#BLOCKME
#0.0.0.0 www.facebook.com
#/BLOCKME

127.0.0.1       localhost
`,
			expects: `# LOCALHOST
127.0.0.1       localhost

#BLOCKME
0.0.0.0 www.facebook.com
#/BLOCKME

127.0.0.1       localhost
`,
		},
		{
			name:      "Empty line should not be commented",
			havingCmd: Cmd{Pause: false},
			having: `# LOCALHOST
#BLOCKME
#0.0.0.0 www.facebook.com

#/BLOCKME
`,
			expects: `# LOCALHOST
#BLOCKME
0.0.0.0 www.facebook.com

#/BLOCKME
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := ioutil.TempFile("/tmp", "hoststest")
			assert.NoError(t, err)
			defer os.Remove(file.Name())

			status, err := ioutil.TempFile("/tmp", "status")
			assert.NoError(t, err)
			defer os.Remove(status.Name())

			err = ioutil.WriteFile(file.Name(), []byte(test.having), 0666)
			assert.NoError(t, err)

			app := NewApp(
				NewHostFile(file.Name()),
				NewFocusBlocker(),
				NewFileStatusManager(status.Name()),
			)

			err = app.Handle(test.havingCmd)
			require.NoError(t, err)

			data, err := ioutil.ReadFile(file.Name())
			assert.NoError(t, err)
			assert.Equal(t, test.expects, string(data))
		})
	}
}
