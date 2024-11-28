package main

import (
	"github.com/heltirj/image_previewer/internal/logger"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	LogLevel    logger.LogLevel `yaml:"logLevel"`
	StoragePath string          `yaml:"storagePath"`
	LRUSize     int             `yaml:"lruSize"`
	Host        string          `yaml:"host"`
	Port        int             `yaml:"port"`
}

func NewConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
