/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : bypass.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 21 Mar 2017 01:48:53 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package bypass

import (
	"context"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
)

/*

INPUT: request.Context()

- "xff0" string

OUTPUT:

- "bypass" bool

*/

var ips []string = []string{"106.48.13.103", "24.43.12.42", "127.0.0.1"}

func init() {
	core.AddAcl(new(AclBypass))
}

type AclBypass struct{}

func (AclBypass) Name() string {
	return "bypass"
}

func (AclBypass) LoadConfig() error {
	return nil
}

func (AclBypass) Pass(r *http.Request) (bool, *http.Request, error) {
	pass := false

	ip := r.RemoteAddr

	if in(ips, ip) {
		pass = true
	}

	xff0, ok := r.Context().Value("xff0").(string)
	if ok {
		if in(ips, xff0) {
			pass = true
		}
	}

	ctx := r.Context()

	ctx = context.WithValue(r.Context(), "bypass", pass)

	r = r.WithContext(ctx)

	return true, r, nil
}

func in(is []string, i string) bool {
	for _, v := range is {
		if v == i {
			return true
		}
	}
	return false
}
