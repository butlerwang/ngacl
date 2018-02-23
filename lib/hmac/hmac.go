/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : hmac.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 06 Jun 2017 12:35:27 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package hmac

import (
	_hmac "crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*

INPUT: request.Context()

- "xff0" string
- "bypass" bool

OUTPUT:

*/

func init() {
	core.AddAcl(new(AclHmac))
}

type AclHmac struct{}

func (AclHmac) Name() string {
	return "hmac"
}

func (AclHmac) LoadConfig() error {
	return nil
}

func getQuery(r *http.Request, key string) string {
	val := r.URL.Query().Get(strings.ToUpper(key))
	if len(val) > 0 {
		return val
	}

	return r.URL.Query().Get(strings.ToLower(key))
}

func (AclHmac) Pass(r *http.Request) (bool, *http.Request, error) {
	if val, ok := r.Context().Value("bypass").(bool); ok && val {
		return true, r, nil
	}
	p1 := getQuery(r, "P1")
	p2 := getQuery(r, "P2")
	p3 := getQuery(r, "P3")
	p4 := getQuery(r, "P4")

	acl := getQuery(r, "ACL")

	k1 := r.Header.Get("Ccna-Acl-Hmac-Key1")
	k2 := r.Header.Get("Ccna-Acl-Hmac-Key2")
	if len(k1) == 0 {
		return false, r, errors.New("missing header 'Ccna-Acl-Hmac-Key1'")
	}
	if len(k2) == 0 {
		return false, r, errors.New("missing header 'Ccna-Acl-Hmac-Key2'")
	}
	var key string

	t, err := strconv.ParseInt(p1, 10, 64)
	if err != nil {
		return false, r, errors.New("P1 is not int")
	}
	timestamp := time.Unix(t, 0)
	if time.Now().After(timestamp) {
		return false, r, errors.New("timestamp expired " + fmt.Sprint(time.Now().Sub(timestamp)))
	}
	switch p2 {
	case "1":
		key = k1
	case "2":
		key = k2
	default:
		return false, r, errors.New("P2 not supported " + p2)
	}
	var payload []byte
	if len(acl) == 0 {
		payload = []byte(fmt.Sprintf("%s%s", r.URL.Path, p1))
	} else {
		payload = []byte(fmt.Sprintf("%s%s", acl, p1))
	}
	var res string
	switch p3 {
	case "1":
		h := _hmac.New(sha1.New, []byte(key))
		h.Write(payload)
		res = base64.StdEncoding.EncodeToString(h.Sum(nil))
	case "2":
		h := _hmac.New(sha256.New, []byte(key))
		h.Write(payload)
		res = base64.StdEncoding.EncodeToString(h.Sum(nil))
	case "3":
		h := _hmac.New(md5.New, []byte(key))
		h.Write(payload)
		res = base64.StdEncoding.EncodeToString(h.Sum(nil))
	default:
		return false, r, errors.New("P3 not supported " + p3)
	}

	res = strings.ToLower(res)
	dec, _ := url.QueryUnescape(p4)
	dec = strings.ToLower(dec)
	p4 = strings.ToLower(p4)

	if res == p4 || res == dec {
		return true, r, nil
	}
	return false, r, errors.New("auth failed " + res + " " + p4)
}
