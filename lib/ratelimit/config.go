/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : config.go

* Purpose :

* Creation Date : 03-18-2017

* Last Modified : Sun 19 Mar 2017 02:52:44 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package ratelimit

import (
	"github.com/go-ini/ini"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var config *ini.File
var ratelimitConfigs map[string]*RatelimitConfig

type BucketConf struct {
	// every xx, give a token
	Dur time.Duration
	// max have xx token, if no token left, then reqire wait
	Max int64
}

type RatelimitConfig struct {
	RateLimit  []string                       `ini:"ratelimit,,allowshadow"`
	BucketConf map[*regexp.Regexp]*BucketConf `ini:"-"`
}

func (a *AclRatelimit) LoadConfig() error {
	var err error
	p := path.Join(core.Config().IncludeDir, a.Name())
	config, err = ini.Load(p)
	if err != nil {
		return err
	}
	ratelimitConfigs = make(map[string]*RatelimitConfig)
	for _, c := range config.Sections() {
		if c.Name() == "DEFAULT" {
			continue
		}
		// 		log.Println(c.Name())
		a := new(RatelimitConfig)
		c.MapTo(&a)
		a.BucketConf = make(map[*regexp.Regexp]*BucketConf)
		for _, v := range a.RateLimit {
			part := strings.Fields(v)
			if len(part) != 3 {
				continue
			}
			re, err := regexp.Compile(part[0])
			if err != nil {
				log.Println(err.Error())
				continue
			}
			dur, err := time.ParseDuration(part[1])
			if err != nil {
				log.Println(err.Error())
				continue
			}
			max, err := strconv.ParseInt(part[2], 10, 64)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			a.BucketConf[re] = &BucketConf{dur, max}
		}
		// 		log.Println(a)
		ratelimitConfigs[c.Name()] = a
	}
	return nil
}
