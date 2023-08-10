package main

type Config struct {
	Token        string `yaml:"token"`
	ImagePath    string `yaml:"imagePath"`
	AdminAccount int64  `yaml:"adminAccount"`
	Debug        bool   `yaml:"debug"`
	Info         string `yaml:"info"`
}
