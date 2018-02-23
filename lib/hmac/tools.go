/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : tools.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 06 Jun 2017 12:24:09 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package hmac

import (
	_hmac "crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

func generate(key, u, p1, p2, p3, acl string) string {
	var b string
	var payload []byte
	if len(acl) == 0 {
		payload = []byte(fmt.Sprintf("%s%s", u, p1))
	} else {
		payload = []byte(fmt.Sprintf("%s%s", acl, p1))
	}
	switch p3 {
	case "1":
		h := _hmac.New(sha1.New, []byte(key))
		h.Write(payload)
		b = base64.StdEncoding.EncodeToString(h.Sum(nil))

	case "2":
		h := _hmac.New(sha256.New, []byte(key))
		h.Write(payload)
		b = base64.StdEncoding.EncodeToString(h.Sum(nil))
	case "3":
		h := _hmac.New(md5.New, []byte(key))
		h.Write(payload)
		b = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
	return url.QueryEscape(b)
}
