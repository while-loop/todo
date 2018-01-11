package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"path/filepath"
	"strings"

	"fmt"
	"os"

	"github.com/pkg/errors"
	configp "github.com/while-loop/todo/pkg/publisher/config"
	configt "github.com/while-loop/todo/pkg/tracker/config"
	configv "github.com/while-loop/todo/pkg/vcs/config"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TrackerConfig   *configt.TrackerConfig   `json:"trackers" yaml:"trackers"`
	VcsConfig       *configv.VcsConfig       `json:"vcs" yaml:"vcs"`
	PublisherConfig *configp.PublisherConfig `json:"publishers" yaml:"publishers"`
}

func ParseFile(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config file")
	}
	defer f.Close()

	ext := strings.TrimLeft(strings.ToLower(filepath.Ext(filePath)), ".")
	if isJson(ext) {
		return ParseJson(f)
	} else if isYaml(ext) {
		return ParseYaml(f)
	}

	return nil, fmt.Errorf("unsupported config extension %s", ext)
}

func ParseJson(reader io.ReadCloser) (*Config, error) {
	defer reader.Close()

	var c Config
	return &c, json.NewDecoder(reader).Decode(&c)
}

func ParseYaml(reader io.ReadCloser) (*Config, error) {
	defer reader.Close()

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	var c Config
	err = yaml.UnmarshalStrict(bs, &c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	return &c, nil
}

func isYaml(ext string) bool {
	return ext == "yml" || ext == "yaml"
}

func isJson(ext string) bool {
	return ext == "json"
}
