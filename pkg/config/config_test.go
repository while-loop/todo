package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYmlParse(t *testing.T) {
	a := require.New(t)
	conf, err := ParseFile(filepath.Join(cwd(t), "../../", "config.yml"))
	a.NoError(err)

	a.Equal("my_access_token", conf.VcsConfig.Github.AccessToken)
}

func TestYamlJsonSame(t *testing.T) {
	a := require.New(t)
	ymlConf, err := ParseFile(filepath.Join(cwd(t), "../../", "config.yml"))
	a.NoError(err)

	jsonConf, err := ParseFile(filepath.Join(cwd(t), "../../", "config.json"))
	a.NoError(err)

	a.Equal("my_access_token", ymlConf.TrackerConfig.Github.AccessToken)
	a.Equal("my_access_token", jsonConf.TrackerConfig.Github.AccessToken)

	a.Equal(jsonConf, ymlConf)
}

func TestFileNotFoundErr(t *testing.T) {
	a := require.New(t)
	c, err := ParseFile("fakefile.json")
	a.Nil(c)
	a.Contains(err.Error(), "failed to load")
}

func TestFileExts(t *testing.T) {
	a := require.New(t)
	tcs := []struct {
		name        string
		err         bool
		errContains string
	}{
		{"test1.toml", true, "unsupported config extension"},
		{"test1.exe", true, "unsupported config extension"},
		{"test1", true, "unsupported config extension"},
		{"test1.JSON", false, ""},
		{"test1.json", false, ""},
		{"test1.YML", false, ""},
		{"test1.yml", false, ""},
		{"test1.yaml", false, ""},
		{"test1.YAML", false, ""},
	}

	for idx, tc := range tcs {
		t.Run(fmt.Sprintf("%d-%s", idx, tc.name), func(t *testing.T) {
			withTempFile(t, tc.name, func(f *os.File) {
				require.NoError(t, json.NewEncoder(f).Encode(tmpConfig))
				c, err := ParseFile(f.Name())
				if tc.err {
					a.Nil(c)
					a.Contains(err.Error(), tc.errContains)
				} else {
					a.Nil(err)
					a.Equal(tmpConfig, c)
				}
			})
		})
	}
}

func withTempFile(t *testing.T, name string, fileFunc func(f *os.File)) {
	tmpDir, err := ioutil.TempDir("", "todo")
	require.NoError(t, err)
	f, err := os.Create(filepath.Join(tmpDir, name))
	require.NoError(t, err)
	require.NotNil(t, f)
	fileFunc(f)
	f.Close()
	require.NoError(t, os.RemoveAll(tmpDir))
}

func cwd(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	return dir
}

var tmpConfig = &Config{}
