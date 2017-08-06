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
var version = [3]byte{0, 0, 4}

const source = "https://github.com/aerth/helper"

type Config struct {
	Name      string // what helper is to be named
	Owner     string // who helper knows as its owner
	Resources map[string]string
	Bookmarks map[string]string
	commands  map[string]func()
	Browser   string       // path to browser executable
	OpenLinks bool         // open links in browser
	user      *osuser.User // user info
	configDir string       // where config file is located
	file      string       // config file name
}

var DefaultResources = map[string]string{
	"g":   "https://encrypted.google.com/search?q=%s",
	"ddg": "https://duckduckgo.com/lite?q=%s",
}

var DefaultBookmarks = map[string]string{
	"gonew":   "https://github.com/search?o=desc&q=language%3Ago+stars%3A1&repo=&s=updated&start_value=1&type=Repositories",
	"gotrend": "https://github.com/trending/go",
	"news":    "https://news.ycombinator.com/",
	"r":       "https://reddit.com/r/golang/",
}

var DefaultCommands = map[string]func(){
	"help": func() {
		fmt.Println("Available Resources:")
		for k, v := range DefaultResources {
			fmt.Println("\t", k, v)
		}
		fmt.Println("\nAvailable Bookmarks:")
		for k, v := range DefaultBookmarks {
			fmt.Println("\t", k, v)
		}
	},
}

var DefaultConfig = &Config{
	Name:      "Cere",
	Owner:     "Master",
	Resources: DefaultResources,
	commands:  DefaultCommands,
	Bookmarks: DefaultBookmarks,
	Browser:   "firefox", // sensible-browser ?
	OpenLinks: true,
}

func init() {
	// remove log flags
	log.SetFlags(0)
	log.SetPrefix("")
	log.Println(versionString())
	log.Println(source)
	flag.Usage = func() {
		fmt.Println()
		DefaultCommands["help"]()
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
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
	if err := config.getcommand(); err != nil {
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
	config.Bookmarks = DefaultBookmarks
	config.commands = DefaultCommands
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
func (c *Config) getcommand() error {
	// no args, open google
	if flag.NArg() == 0 {
		return c.OpenLink(fmt.Sprintf(c.Resources["ddg"], "!g+\"aerth\"+helper+github"))
	}
	args := flag.Args()

	// check commands
	if command, ok := c.commands[args[0]]; ok {
		command()
		os.Exit(0)
	}

	// check bookmarks
	if link, ok := c.Bookmarks[args[0]]; ok {
		return c.OpenLink(link)
	}

	// check resources
	if len(args) > 1 {
		if query, ok := c.Resources[args[0]]; ok {
			return c.OpenLink(fmt.Sprintf(query, url.QueryEscape(strings.Join(args[1:], " "))))
		}
	}
	// default: duck duck go
	cmd := url.QueryEscape(strings.Join(args, " ")) // escape concatenated arguments
	link := fmt.Sprintf(c.Resources["ddg"], cmd)
	return c.OpenLink(link)

}

// OpenLink in browser
func (c *Config) OpenLink(cmd string) error {
	log.Println("\n" + cmd + "\n")
	if !c.OpenLinks {
		return nil
	}

	// open in browser
	go func() {
		if err := browser.OpenURL(cmd); err != nil {
			log.Fatal(err)

		}
	}()
	<-time.After(time.Millisecond * 500)

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
