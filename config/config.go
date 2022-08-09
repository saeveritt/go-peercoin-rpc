package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Testnet  bool
	Username string
	Password string
	Host     string
	Port     int
}

var (
	ErrUsernameNotFound  = errors.New("username not found")
	ErrPasswordNotFound  = errors.New("password not found")
	ErrHostNotFound      = errors.New("host not found")
	ErrPortNotFound      = errors.New("port not found")
	ErrDirectoryNotFound = errors.New("directory not found")
)

// load config from json file into Config struct
func LoadConfig() (*Config, error) {
	// create Config struct
	config := &Config{}
	// get parent directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// get config from json file
	err = config.load(wd + "/config.json")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) load(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// unmarshal json into Config struct
	err = json.Unmarshal(content, c)
	log.Printf("Loaded config: %v", c)
	log.Printf("Loaded config from: %s", path)
	if err != nil {
		return err
	}
	// set config username and password
	c.setUsernamePassword()
	return nil
}

// set username and password
func (c *Config) setUsernamePassword() {
	username, err := GetUserName()
	if err != nil {
		log.Printf("Username not found in environment variable GRPPC_USERNAME")
		return
	}
	c.Username = username

	password, err := GetPassword()
	if err != nil {
		log.Printf("Password not found in environment variable GRPPC_PASSWORD")
		return
	}
	c.Password = password
}

// get username from environment variable
func GetUserName() (string, error) {
	username := os.Getenv("GRPPC_USERNAME")
	if username == "" {
		return username, ErrUsernameNotFound
	}
	return username, nil
}

// get password from environment variable
func GetPassword() (string, error) {
	password := os.Getenv("GRPPC_PASSWORD")
	if password == "" {
		return password, ErrPasswordNotFound
	}
	return password, nil
}
