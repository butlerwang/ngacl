/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : md5.go

* Purpose :

* Creation Date : 06-08-2017

* Last Modified : Thu 29 Jun 2017 09:39:00 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package tkmd5

import (
	"crypto/md5"
	"fmt"
)

func md5sum(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
