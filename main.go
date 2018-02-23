/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 03-15-2017

* Last Modified : Thu 08 Jun 2017 07:18:51 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"github.com/NYTimes/logrotate"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/bypass"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/hmac"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/querylog"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/ratelimit"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/s3"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/tkmd5"
	_ "gitlab.ccnanext.com/ccna/ngacl/lib/waf"
	"log"
	"net/http"
	"os"
	"runtime"
)

func init() {
	logfile, err := logrotate.NewFile("/var/log/ngacl.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
}

func main() {
	if *core.Version {
		fmt.Printf("%s.%s\n", VERSION, buildtime)
		os.Exit(0)
	}
	if *core.ConfigTest {
		os.Exit(0)
	}
	for _, a := range core.Acls() {
		err := a.LoadConfig()
		if err != nil {
			log.Println(err.Error())
		}
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.ListenAndServe("127.0.0.1:8104", LogHandler(http.HandlerFunc(AclHandler)))
}
