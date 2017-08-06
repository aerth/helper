// Copyright 2017 aerth. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

// helper command line assistant
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	osuser "os/user"
	"strings"
	"time"

	"github.com/pkg/browser"

	"path/filepath"
)

// helper version
var version = [3]byte{0, 0, 3}

type Config struct {
	Name      string // what helper is to be named
	Owner     string // who helper knows as its owner
	Resources map[string]string
	Browser   string       // path to browser executable
	OpenLinks bool         // open links in browser
	user      *osuser.User // user info
	configDir string       // where config file is located
	file      string       // config file name
}

var DefaultResources = map[string]string{
	"google": "https://encrypted.google.com/search?q=%s",
	"ddg":    "https://duckduckgo.com/lite?q=%s",
}

var DefaultConfig = &Config{
	Name:      "Cere",
	Owner:     "Master",
	Resources: DefaultResources,
	Browser:   "firefox", // sensible-browser ?
	OpenLinks: true,
}

func init() {
	// remove log flags
	log.SetFlags(0)
	log.SetPrefix("")
	log.Println(versionString())
}

// stringer for version
func versionString() string {
	return fmt.Sprintf("helper v%v.%v.%v", version[0], version[1], version[2])
}

// parse flags, run command
func main() {
	noopenlinks := flag.Bool("x", false, "dont open links, just print them")
	flag.Parse()
	config := readconfig()
	if *noopenlinks {
		config.OpenLinks = false
	}
	command := getcommand()
	if err := config.Run(command); err != nil {
		log.Fatal(err)
	}
}

// read config file (~/.config/helper/config.json)
func readconfig() *Config {
	user, err := osuser.Current()
	if err != nil {
		log.Fatal(err)
	}

	// configdir /home/user/.config/helper/
	configdir := filepath.Join(user.HomeDir, ".config", "helper")
	err = os.MkdirAll(configdir, 0700)
	if err != nil {
		log.Fatal(err)
	}

	// configfile /home/user/.config/helper/config.json
	configfile := filepath.Join(configdir, "config.json")
	b, err := ioutil.ReadFile(configfile)
	if err != nil {
		if err.Error() != fmt.Sprintf("open %s: no such file or directory", configfile) {
			log.Fatal(err)
		}

		// write default config
		err = ioutil.WriteFile(configfile, DefaultConfig.Marshal(), 0700)
		if err != nil {
			log.Fatal(err)
		}
		b = DefaultConfig.Marshal() // avoid restart (first boot)

	}
	// unmarshal with defaults
	config := new(Config)
	config.Name = "Cere"
	config.Owner = user.Username
	config.Resources = DefaultResources
	config.OpenLinks = true
	config.user = user
	config.configDir = configdir
	config.file = configfile
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// get command from command line arguments
func getcommand() string {
	if flag.NArg() == 0 {
		return "!g+\"aerth\"+helper+github"
	}
	return url.QueryEscape(strings.Join(flag.Args(), " ")) // escape concatenated arguments

}

// Run a command
func (c *Config) Run(cmd string) error {
	link := fmt.Sprintf(c.Resources["ddg"], cmd)
	log.Println(link)
	if c.OpenLinks {
		// open in browser
		go func() {
			if err := browser.OpenURL(link); err != nil {
				log.Fatal(err)

			}
		}()
		<-time.After(time.Millisecond * 500)
	}
	return nil
}

// Marshal into []byte
func (c *Config) Marshal() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}
