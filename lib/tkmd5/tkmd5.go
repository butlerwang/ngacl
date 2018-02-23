/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : tkmd5.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Thu 29 Jun 2017 09:39:49 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package tkmd5

import (
	"errors"
	"fmt"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	SPLITS = []string{"&h=", "?h="}
)

func init() {
	core.AddAcl(new(AclTkmd5))
}

type AclTkmd5 struct{}

func (AclTkmd5) Name() string {
	return "tkmd5"
}

func (AclTkmd5) LoadConfig() error {
	return nil
}

type Options struct {
	IgnoreExpire bool
}

func parseOptions(s string) *Options {
	o := new(Options)
	kvs := strings.Split(s, ",")
	for _, kv := range kvs {
		kv = strings.TrimSpace(kv)
		p := strings.Split(kv, "=")
		if len(p) == 2 {
			switch p[0] {
			case "ignore_expire":
				if p[1] == "1" {
					o.IgnoreExpire = true
				}
			}
		}
	}
	return o
}

func (AclTkmd5) Pass(r *http.Request) (bool, *http.Request, error) {
	if val, ok := r.Context().Value("bypass").(bool); ok && val {
		return true, r, nil
	}
	e := r.URL.Query().Get("e")
	h := r.URL.Query().Get("h")

	var base string
	for _, sp := range SPLITS {
		p := strings.Split(r.URL.String(), sp)
		if len(p) > 1 {
			base = p[0]
			break
		}
	}
	secret := r.Header.Get("ccna-acl-tkmd5-secret")

	o := r.Header.Get("ccna-acl-tkmd5-options")
	opt := parseOptions(o)

	if !opt.IgnoreExpire || len(e) > 0 {

		t, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return false, r, errors.New("e is not int")
		}
		timestamp := time.Unix(t, 0)
		if time.Now().After(timestamp) {
			return false, r, errors.New("timestamp expired " + fmt.Sprint(time.Now().Sub(timestamp)))
		}

	}

	payload := secret + base
	hash := md5sum(payload)
	if hash == h {
		return true, r, nil
	}

	return false, r, errors.New("auth failed " + hash + " " + h)
}
