/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : hmac_test.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 06 Jun 2017 12:28:36 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package hmac

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHmac(t *testing.T) {
	u, _ := url.Parse("http://a.com/b")

	type T struct {
		p1   string
		p2   string
		p3   string
		p4   string
		acl  string
		pass bool
	}

	m := map[string]string{
		"Key1": "abc",
		"Key2": "def",
	}

	goodTime := fmt.Sprint(time.Now().Add(time.Minute).Unix())
	badTime := fmt.Sprint(time.Now().Add(-time.Minute).Unix())

	testCase := []T{
		{"a", "1", "1", "1", "", false},
		{badTime, "1", "1", "1", "", false},
		{goodTime, "1", "1", "1", "", false},
		{goodTime, "3", "1", "1", "", false},
		{goodTime, "1", "4", "1", "", false},
		{goodTime, "1", "1", generate(m["Key1"], u.Path, goodTime, "1", "1", ""), "", true},
		{goodTime, "1", "2", generate(m["Key1"], u.Path, goodTime, "1", "2", ""), "", true},
		{goodTime, "1", "3", generate(m["Key1"], u.Path, goodTime, "1", "3", ""), "", true},
		{goodTime, "2", "1", generate(m["Key2"], u.Path, goodTime, "2", "1", ""), "", true},
		{goodTime, "2", "2", generate(m["Key2"], u.Path, goodTime, "2", "2", ""), "", true},
		{goodTime, "2", "3", generate(m["Key2"], u.Path, goodTime, "2", "3", ""), "", true},

		{goodTime, "1", "1", generate(m["Key1"], u.Path, goodTime, "1", "1", "/"), "/", true},
		{goodTime, "1", "2", generate(m["Key1"], u.Path, goodTime, "1", "2", "/"), "/", true},
		{goodTime, "1", "3", generate(m["Key1"], u.Path, goodTime, "1", "3", "/"), "/", true},
		{goodTime, "2", "1", generate(m["Key2"], u.Path, goodTime, "2", "1", "/"), "/", true},
		{goodTime, "2", "2", generate(m["Key2"], u.Path, goodTime, "2", "2", "/"), "/", true},
		{goodTime, "2", "3", generate(m["Key2"], u.Path, goodTime, "2", "3", "/"), "/", true},

		{goodTime, "2", "2", generate(m["Key2"], u.Path, goodTime, "2", "2", "/a"), "/", false},
	}

	for _, c := range testCase {
		v := u.Query()
		v.Set("P1", c.p1)
		v.Set("P2", c.p2)
		v.Set("P3", c.p3)
		v.Set("P4", c.p4)
		v.Set("ACL", c.acl)

		u.RawQuery = v.Encode()
		req, _ := http.NewRequest("GET", u.String(), nil)
		req.Header.Set("Ccna-Acl-Hmac-Key1", m["Key1"])
		req.Header.Set("Ccna-Acl-Hmac-Key2", m["Key2"])
		pass, _, err := (&AclHmac{}).Pass(req)
		if err != nil {
			log.Println(err.Error())
		}
		if pass != c.pass {
			t.Fatal(u.String())
		}
	}
}
