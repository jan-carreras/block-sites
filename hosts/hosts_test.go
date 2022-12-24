package hosts

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestApp_Handle(t *testing.T) {
	tests := []struct {
		name    string
		having  string
		expects string
	}{
		{
			name: "Base Case",
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
			name: "Empty line should not be commented",
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

			err = ioutil.WriteFile(file.Name(), []byte(test.having), 0666)
			assert.NoError(t, err)

			app := NewApp(
				NewHostFile(file.Name()),
				NewFocusBlocker(),
			)

			err = app.Handle()
			assert.NoError(t, err)

			data, err := ioutil.ReadFile(file.Name())
			assert.NoError(t, err)
			assert.Equal(t, test.expects, string(data))
		})
	}
}
