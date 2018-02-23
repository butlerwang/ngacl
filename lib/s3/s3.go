/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : s3.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 18 Apr 2017 05:33:21 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package s3

import (
	"context"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
	"time"
)

// in

// out
// s3header http.Header

func init() {
	core.AddAcl(new(AclS3))
}

type AclS3 struct{}

func (AclS3) Name() string {
	return "s3"
}

func (AclS3) LoadConfig() error {
	return nil
}

func (AclS3) Pass(r *http.Request) (bool, *http.Request, error) {
	bucket := r.Header.Get("ccna-s3-bucket")
	accessKey := r.Header.Get("ccna-s3-accesskey")
	secretKey := r.Header.Get("ccna-s3-secretkey")
	region := r.Header.Get("ccna-s3-region")

	a := &Auth{accessKey, secretKey}
	c := a.NewV4Context(region, bucket, r.Host, r.URL.Path, "", "", time.Now().UTC())

	ctx := context.WithValue(r.Context(), "s3header", c.GetHeader())

	r = r.WithContext(ctx)

	return true, r, nil
}
