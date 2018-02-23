/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : core.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Sat 18 Mar 2017 08:15:06 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package core

import (
	"flag"
	"github.com/go-ini/ini"
	// 	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	configPath *string = flag.String("c", "/etc/ngacl.ini", "config path")
	ConfigTest *bool   = flag.Bool("t", false, "config test and exit")
	Version    *bool   = flag.Bool("v", false, "config test and exit")
	config     *CoreConfig
)

type CoreConfig struct {
	File       *ini.File
	IncludeDir string
}

func init() {
	flag.Parse()
	loadConfig()
}

func loadConfig() {
	cfg, err := ini.Load(*configPath)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	// 	enableLog, err := cfg.Section("core").Key("log").Bool()
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds)
	// 	if !enableLog {
	// 		log.SetOutput(ioutil.Discard)
	// 	}

	include := cfg.Section("core").Key("include").String()
	if f, err := os.Stat(include); err != nil || !f.IsDir() {
		log.Println("[core]include", include, "is not dir")
		os.Exit(1)
	}

	config = &CoreConfig{
		File:       cfg,
		IncludeDir: include,
	}
}

type Acl interface {
	Pass(*http.Request) (bool, *http.Request, error)
	Name() string
	LoadConfig() error
}

func Acls() []Acl {
	return aclList
}

var aclList []Acl

func AddAcl(a Acl) {
	aclList = append(aclList, a)
}

func Config() *CoreConfig {
	return config
}
