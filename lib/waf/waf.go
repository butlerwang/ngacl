/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : waf.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Thu 11 May 2017 05:20:38 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package waf

import (
	"errors"
	"fmt"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
	"strings"
)

func init() {
	core.AddAcl(new(AclWaf))
}

type AclWaf struct{}

func (AclWaf) Name() string {
	return "waf"
}

func (AclWaf) LoadConfig() error {
	return nil
}

func (AclWaf) Pass(r *http.Request) (bool, *http.Request, error) {
	// CVE-2017-5638
	// http://blog.talosintelligence.com/2017/03/apache-0-day-exploited.html
	contentType := r.Header.Get("content-type")
	for _, v := range []string{"@java", "new java", "flush()"} {
		if strings.Contains(contentType, v) {
			return false, r, errors.New(fmt.Sprintf(ERR_VUN_CVE_2017_5638, contentType))
		}
	}
	return true, r, nil
}
