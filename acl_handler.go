/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : acl_handler.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Wed 14 Jun 2017 10:12:54 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func AclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value("id").(string)
	if !ok {
		id = "-"
	}

	pass := true
	methods := strings.Fields(r.Header.Get("ccna-acl-method"))
	if len(methods) == 0 {
		w.WriteHeader(200)
		return
	}
	rawPath := r.URL.String()
	scheme := r.Header.Get("ccna-acl-scheme")
	host := r.Host

	r.URL, _ = url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, rawPath))

Here:
	for _, method := range methods {
		for _, acl := range core.Acls() {
			if method == acl.Name() {
				var p bool
				var err error
				p, r, err = acl.Pass(r)
				if !p && err != nil {
					pass = false
					log.Println(id, r.RemoteAddr, method, r.Host+r.RequestURI, err.Error())
					break Here
				} else if err != nil {
					log.Println(id, r.RemoteAddr, method, r.Host+r.RequestURI, err.Error())
				}
			}
		}
	}
	if pass {
		if h, ok := r.Context().Value("s3header").(http.Header); ok {
			for k, v := range h {
				w.Header().Set(k, v[0])
			}
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}
}
