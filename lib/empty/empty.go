/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : empty.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Fri 17 Mar 2017 07:32:09 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package empty

import (
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
)

func init() {
	core.AddAcl(new(AclEmpty))
}

type AclEmpty struct{}

func (AclEmpty) Name() string {
	return "empty"
}

func (AclEmpty) LoadConfig() error {
	return nil
}

func (AclEmpty) Pass(r *http.Request) (bool, *http.Request, error) {
	return true, r, nil
}
