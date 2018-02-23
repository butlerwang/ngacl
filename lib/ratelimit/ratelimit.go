/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : ratelimit.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 16 May 2017 10:00:27 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package ratelimit

import (
	"errors"
	"fmt"
	_ratelimit "github.com/juju/ratelimit"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"net/http"
	"sync"
	"time"
)

/*

INPUT: request.Context()

- "xff0" string
- "bypass" bool

OUTPUT:

*/

var (
	tokenBuckets map[string]*_ratelimit.Bucket
	lock         = new(sync.RWMutex)
)

func init() {
	tokenBuckets = make(map[string]*_ratelimit.Bucket)
	core.AddAcl(new(AclRatelimit))
	data = &Data{
		Meter: make(map[Host]map[Ip]int64),
		start: time.Now(),
		lock:  new(sync.RWMutex),
	}
	ch = make(chan *http.Request)
	go cron()
	go writeData(ch)
	go server()
}

func cron() {
	done := make(chan bool)
	ticker := time.NewTicker(time.Minute * 60)
	go func() {
		for range ticker.C {
			cleanMap()
			cleanMapBlackList()
		}
	}()
	<-done
}

func cleanMap() {
	lock.Lock()
	newBuckets := make(map[string]*_ratelimit.Bucket)
	tokenBuckets = newBuckets
	lock.Unlock()
}

type AclRatelimit struct {
}

func (AclRatelimit) Name() string {
	return "ratelimit"
}

func (AclRatelimit) Pass(r *http.Request) (bool, *http.Request, error) {
	// if other config is filling ctx to bypass acl, then return
	if val, ok := r.Context().Value("bypass").(bool); ok && val {
		return true, r, nil
	}
	// if Host config exist, then do ratelimit, otherwise do not waste time
	c, ok := ratelimitConfigs[r.Host]
	if !ok {
		return true, r, nil
	}

	ip := r.RemoteAddr
	if val, ok := r.Context().Value("xff0").(string); ok {
		ip = val
	}

	// if ip already blocked, then do not waste time
	data.lock.RLock()
	val, ok := data.Meter[Host(r.Host)]
	if ok {
		_, ok := val[Ip(ip)]
		data.lock.RUnlock()
		if ok {
			ch <- r
			return false, r, errors.New("ip blocked " + ip)
		}
	} else {
		data.lock.RUnlock()
	}

	for re, bucketConf := range c.BucketConf {
		if re.MatchString(r.URL.RequestURI()) {
			name := fmt.Sprintf("%s:%s:%s", ip, r.Host, re.String())
			lock.Lock()
			tokenBucket, ok := tokenBuckets[name]
			if !ok {
				tokenBucket = _ratelimit.NewBucket(bucketConf.Dur, bucketConf.Max)
				tokenBuckets[name] = tokenBucket
			}
			lock.Unlock()

			available := tokenBucket.TakeAvailable(1)
			if available <= 0 {
				ch <- r
				return false, r, errors.New("rate limit " + name)
			}
		}
	}
	return true, r, nil
}
