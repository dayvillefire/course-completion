package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mail struct {
		ServerName string `yaml:"server-name"`
		ServerPort int    `yaml:"server-port"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Subject    string `yaml:"subject"`
		Template   string `yaml:"template"`
		FromName   string `yaml:"from-name"`
		FromEmail  string `yaml:"from-email"`
	} `yaml:"mail"`
	Template struct {
		PdfFile            string                 `yaml:"pdf-file"`
		GlobalReplacements map[string]string      `yaml:"global-replacements"`
		Replacements       map[string]Replacement `yaml:"replacements"`
	} `yaml:"template"`
	Data struct {
		NameField  string `yaml:"name-field"`
		EmailField string `yaml:"email-field"`
	} `yaml:"data"`
}

type Replacement struct {
	Default    string  `yaml:"default"`
	FontFamily string  `yaml:"font-family"`
	FontJson   string  `yaml:"font-json"`
	FontSize   float64 `yaml:"font-size"`
	NameY      float64 `yaml:"name-y"`
}

func loadConfig(configPath string) (*Config, error) {
	c := &Config{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), c)

	return c, err
}
