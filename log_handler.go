/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : log_handler.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 21 Mar 2017 01:46:10 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		id := random(16)
		ctx := context.Background()
		ctx = context.WithValue(ctx, "id", id)

		remoteAddr := r.Header.Get("Ccna-Acl-Remote-Addr")
		if len(remoteAddr) == 0 {
			remoteAddr = "-"
		}
		r.RemoteAddr = remoteAddr

		method := r.Header.Get("Ccna-Acl-Request-Method")
		if len(method) == 0 {
			method = "-"
		}
		r.Method = method

		var ipchan []string
		xff := r.Header.Get("X-Forwarded-For")
		if len(xff) > 0 {
			for _, i := range strings.Split(xff, ",") {
				if p := net.ParseIP(i); p != nil {
					ipchan = append(ipchan, i)
				}
			}
		}

		ctx = context.WithValue(ctx, "xff", ipchan)
		if len(ipchan) > 0 {
			ctx = context.WithValue(ctx, "xff0", ipchan[0])
		}

		w.Header().Set("X-Cc-Auth-Id", id)
		writer := statusWriter{w, 0, 0}

		r = r.WithContext(ctx)
		next.ServeHTTP(&writer, r)

		log.Println(id, r.RemoteAddr, writer.status, writer.length, r.Method, r.Host+r.RequestURI, time.Since(t1))

	})
}
