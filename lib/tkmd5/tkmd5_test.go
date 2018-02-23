/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : tkmd5_test.go

* Purpose :

* Creation Date : 06-08-2017

* Last Modified : Fri 30 Jun 2017 12:39:23 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package tkmd5

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func Test_Md5sum(t *testing.T) {
	if md5sum("rspEQ9JQhakwghttp://cdn.wacom.com/s/developer/CommonDeviceLibrarySDK.zip?e=1496779545") != "bc85007f6d8d5e40d5680fbaa9ac183d" {
		t.Fatal("md5 error")
	}
}

func Test_Pass(t *testing.T) {
	// 	r, _ := http.NewRequest("GET", "http://cdn.wacom.com/s/developer/CommonDeviceLibrarySDK.zip?e=1496779545", nil)
	secret := "rspEQ9JQhakwg"
	u, _ := url.Parse("http://a.com/b")

	type T struct {
		e    string
		h    string
		pass bool
	}

	goodTime := fmt.Sprint(time.Now().Add(time.Minute).Unix())
	badTime := fmt.Sprint(time.Now().Add(-time.Minute).Unix())

	testCase := []T{
		{"a", "", false},
		{goodTime, u.String(), false},
		{goodTime, md5sum(secret + u.String() + "?e=" + badTime), false},
		{goodTime, md5sum(secret + u.String() + "?e=" + goodTime), true},
		{badTime, md5sum(secret + u.String() + "?e=" + badTime), false},
	}
	for k, c := range testCase {
		v := u.Query()
		v.Set("e", c.e)
		v.Set("h", c.h)
		u.RawQuery = v.Encode()
		req, _ := http.NewRequest("GET", u.String(), nil)
		req.Header.Set("ccna-acl-tkmd5-secret", secret)
		pass, _, err := (&AclTkmd5{}).Pass(req)
		if err != nil {
			log.Println(k, err.Error())
		}
		if pass != c.pass {
			t.Fatal(k, u.String())
		}

	}

	u, _ = url.Parse("http://a.com/b")
	v := u.Query()
	v.Set("h", md5sum(secret+u.String()))
	u.RawQuery = v.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("ccna-acl-tkmd5-secret", secret)
	req.Header.Set("ccna-acl-tkmd5-options", "ignore_expire=1")
	pass, _, err := (&AclTkmd5{}).Pass(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !pass {
		t.Fatal(u.String())
	}
}
