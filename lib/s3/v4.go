/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : v4.go

* Purpose :

* Creation Date : 03-23-2017

* Last Modified : Tue 18 Apr 2017 05:32:25 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	// 	"log"
	"net/http"
	"time"
)

var (
	timeFormat        = "20060102T150405Z"
	shortTimeFormat   = "20060102"
	XamzContentSHA256 = "x-amz-content-sha256"
	XamzDate          = "x-amz-date"
	Service           = "s3"
)

func (ctx *V4Context) signedHeaders() string {
	if len(ctx.Range) > 0 {
		return fmt.Sprintf("%s;%s;%s;%s", "host", "range", XamzContentSHA256, XamzDate)
	}
	return fmt.Sprintf("%s;%s;%s", "host", XamzContentSHA256, XamzDate)
}

type Auth struct {
	accessKey string
	secretKey string
}

type V4Context struct {
	Region  string
	Bucket  string
	Host    string
	Uri     string
	Query   string
	Amzdate string
	Date    string
	Payload string
	Range   string
	time    time.Time
	auth    *Auth
}

func (a *Auth) NewV4Context(region, bucket, host, uri, query, payload string, t time.Time) *V4Context {
	return &V4Context{
		Region:  region,
		Bucket:  bucket,
		Host:    host,
		Uri:     uri,
		Query:   query,
		Amzdate: t.Format(timeFormat),
		Date:    t.Format(shortTimeFormat),
		Payload: payload,
		time:    t,
		auth:    a,
	}
}

/*
func NewV4Context(region, bucket, host, uri, query, payload, accessKey, secretKey string, t time.Time) *V4Context {
	return &V4Context{
		Region:    region,
		Bucket:    bucket,
		Host:      host,
		Uri:       uri,
		Query:     query,
		Amzdate:   t.Format(timeFormat),
		Date:      t.Format(shortTimeFormat),
		Payload:   payload,
		accessKey: accessKey,
		secretKey: secretKey,
		time:      t,
	}
}
*/

func (V4Context) hashedPayload(payload string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
}

func (ctx *V4Context) canonicalHeaders() string {
	if len(ctx.Range) > 0 {
		return fmt.Sprintf(
			`host:%s
%s:%s
%s:%s
%s:%s
`, ctx.Host,
			"range", ctx.Range,
			XamzContentSHA256, ctx.hashedPayload(ctx.Payload),
			XamzDate, ctx.Amzdate)

	}
	return fmt.Sprintf(
		`host:%s
%s:%s
%s:%s
`, ctx.Host,
		XamzContentSHA256, ctx.hashedPayload(ctx.Payload),
		XamzDate, ctx.Amzdate)
}

func (ctx *V4Context) canonicalRequest() string {
	res := fmt.Sprintf(
		`%s
%s
%s
%s
%s
%s`,
		"GET",
		ctx.Uri,
		ctx.Query,
		ctx.canonicalHeaders(),
		ctx.signedHeaders(),
		ctx.hashedPayload(ctx.Payload),
	)
	// 	log.Println(res)
	return res
}

func (ctx *V4Context) hexCanonicalRequest() string {
	return ctx.hashedPayload(ctx.canonicalRequest())
}

func (ctx *V4Context) stringToSign() string {
	res := fmt.Sprintf(
		`AWS4-HMAC-SHA256
%s
%s/%s/%s/aws4_request
%s`,
		ctx.Amzdate,
		ctx.Date, ctx.Region, Service,
		ctx.hexCanonicalRequest())
	// 	log.Println(res)
	return res
}

func makeHmac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func (ctx *V4Context) signature() string {
	secret := ctx.auth.secretKey
	date := makeHmac([]byte("AWS4"+secret), []byte(ctx.Date))
	region := makeHmac(date, []byte(ctx.Region))
	service := makeHmac(region, []byte(Service))
	credentials := makeHmac(service, []byte("aws4_request"))
	signature := makeHmac(credentials, []byte(ctx.stringToSign()))
	return hex.EncodeToString(signature)
}

func (ctx *V4Context) authorization() string {
	// 	log.Println("sign", ctx.signature())
	return fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/%s/aws4_request, SignedHeaders=%s, Signature=%s", ctx.auth.accessKey, ctx.Date, ctx.Region, Service, ctx.signedHeaders(), ctx.signature())
}

func (ctx *V4Context) GetHeader() http.Header {
	h := make(http.Header)
	h.Set(XamzContentSHA256, ctx.hashedPayload(""))
	h.Set(XamzDate, ctx.Amzdate)
	h.Set("Authorization", ctx.authorization())
	return h
}
