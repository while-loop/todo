package config

type PublisherConfig struct {
	Mock *MockConfig `json:"mock" yaml:"mock"`
}

type MockConfig struct {
	Output string `json:"output" yaml:"output"`
}
