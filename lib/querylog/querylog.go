/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : querylog.go

* Purpose :

* Creation Date : 03-16-2017

* Last Modified : Tue 11 Apr 2017 08:42:20 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package querylog

import (
	"encoding/json"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"gitlab.ccnanext.com/ccna/ngacl/lib/core"
	"html/template"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Host string

type Uri string

type QueryKey string
type QueryValue string

type Data struct {
	Meter  map[Host]map[Uri]map[QueryKey]map[QueryValue]uint64
	start  time.Time
	Uptime time.Duration
	lock   *sync.RWMutex
	Count  uint64
}

type User struct {
	Username string
	Password string
}

var (
	data   *Data
	ch     chan *http.Request
	admins []User = []User{{"ccna", "{SHA}7biP3mux5ydy0FXAN/R4UMD9XJc="}}
	users  []User = []User{{"airasia", "{SHA}80PBfOeg1mlg5l/wUiHeBAiQKsc="}}
)

func init() {
	core.AddAcl(new(AclQuerylog))
	data = &Data{
		Meter: make(map[Host]map[Uri]map[QueryKey]map[QueryValue]uint64),
		start: time.Now(),
		lock:  new(sync.RWMutex),
	}
	users = append(users, admins...)
	ch = make(chan *http.Request)
	go do(ch)
	go cron()
	go server()
}

func cron() {
	done := make(chan bool)
	ticker := time.NewTicker(time.Hour * 24)
	go func() {
		for range ticker.C {
			cleanMap()
		}
	}()
	<-done
}

func cleanMap() {
	meter := make(map[Host]map[Uri]map[QueryKey]map[QueryValue]uint64)
	data.lock.Lock()
	data.start = time.Now()
	data.Meter = meter
	data.Count = 0
	data.lock.Unlock()
}

type AclQuerylog struct{}

func (AclQuerylog) Name() string {
	return "querylog"
}

func (AclQuerylog) LoadConfig() error {
	return nil
}

func (AclQuerylog) Pass(r *http.Request) (bool, *http.Request, error) {
	ch <- r
	return true, r, nil
}

func do(r chan *http.Request) {
	putone := func(host, uri, queryKey, queryValue string) {
		h := Host(host)
		u := Uri(uri)
		qk := QueryKey(queryKey)
		qv := QueryValue(queryValue)

		_, ok := data.Meter[h]
		if !ok {
			m := make(map[Uri]map[QueryKey]map[QueryValue]uint64)
			data.lock.Lock()
			data.Meter[h] = m
			data.lock.Unlock()
		}

		_, ok = data.Meter[h][u]
		if !ok {
			m := make(map[QueryKey]map[QueryValue]uint64)
			data.lock.Lock()
			data.Meter[h][u] = m
			data.lock.Unlock()
		}

		_, ok = data.Meter[h][u][qk]
		if !ok {
			m := make(map[QueryValue]uint64)
			data.lock.Lock()
			data.Meter[h][u][qk] = m
			data.lock.Unlock()
		}

		data.Meter[h][u][qk][qv] += 1
	}
	for val := range r {
		data.Count++
		for qk, qv := range val.URL.Query() {
			var v string
			if len(qv) == 0 {
				v = " "
			} else {
				v = qv[0]
			}
			putone(val.Host, val.URL.Path, qk, v)
		}
		key1 := val.URL.Query().Get("o1") + "-" + val.URL.Query().Get("d1")
		if key1 != "-" {
			putone(val.Host, val.URL.Path, "o1-d1", key1)
		}
	}
}

var IndexTemplate string = `
<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
<p>
{{ .Uptime }}
<br>
{{ .Count }}
</p>
{{ if .Back }}<div><a href="{{.Back}}">back</a></div><br>{{ end }}
{{ range .Index }}
<li><a href="{{.Link}}">{{.Name}}</li>
{{ end }}
<pre><code>{{.Body}}</code></pre>
{{ if .Back }}<br><div><a href="{{.Back}}">back</a></div>{{ end }}
</body>
</html>

`

type Html struct {
	Title string

	Uptime string
	Count  string

	Back string

	Index []RangeIndex

	Body string
}
type RangeIndex struct {
	Name string
	Link string
}

func NewHtml() *Html {
	return &Html{
		Uptime: fmt.Sprintf("uptime: %v", time.Since(data.start)),
		Count:  fmt.Sprintf("count: %v", data.Count),
	}
}

// /list
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 	user, pass, _ := r.BasicAuth()
	// 	if !isAdmin(user, pass) {
	// 		w.WriteHeader(403)
	// 		return
	// 	}

	html := NewHtml()
	var list []string
	for k := range data.Meter {
		name := string(k)
		list = append(list, name)
	}
	sort.Strings(list)
	for _, name := range list {
		html.Index = append(html.Index, RangeIndex{Name: name, Link: "/list/" + name})
	}
	html.Title = "Index"
	t, _ := template.New("index").Parse(IndexTemplate)
	t.Execute(w, html)
}

func secretUser(user, realm string) string {
	for _, u := range users {
		if user == u.Username {
			return u.Password
		}
	}
	return ""
}
func secretAdmin(user, realm string) string {
	for _, u := range admins {
		if user == u.Username {
			return u.Password
		}
	}
	return ""
}

// /list/host.com
func hostHandler(w http.ResponseWriter, r *http.Request) {
	// 	user, pass, _ := r.BasicAuth()
	// 	if !isUser(user, pass) {
	// 		w.WriteHeader(403)
	// 		return
	// 	}

	val := mux.Vars(r)
	h := val["host"]
	host := Host(h)
	html := NewHtml()

	u := r.URL.Query().Get("uri")

	if len(u) == 0 {
		var list []string
		for k := range data.Meter[host] {
			name := string(k)
			list = append(list, name)
		}
		sort.Strings(list)
		for _, name := range list {
			html.Index = append(html.Index, RangeIndex{Name: name, Link: "/list/" + h + "?uri=" + name})
		}
		html.Title = h
		html.Back = "/list"
		t, _ := template.New("index").Parse(IndexTemplate)
		t.Execute(w, html)
	} else {
		uri := Uri(u)

		html.Title = h + u
		html.Body = toJson(data.Meter[host][uri])

		html.Back = "/list/" + h

		t, _ := template.New("index").Parse(IndexTemplate)
		t.Execute(w, html)
	}
}

func server() {
	authAdmin := auth.NewBasicAuthenticator("", secretAdmin)
	authUser := auth.NewBasicAuthenticator("", secretUser)
	r := mux.NewRouter()
	r.HandleFunc("/list", auth.JustCheck(authAdmin, indexHandler))
	r.HandleFunc("/list/{host}", auth.JustCheck(authUser, hostHandler))
	http.ListenAndServe(":8105", gziphandler.GzipHandler(r))
}

func toJson(i interface{}) string {
	var res string
	data.lock.RLock()
	j, err := json.MarshalIndent(i, "", "  ")
	data.lock.RUnlock()
	if err != nil {
		res += fmt.Sprintln(err.Error())
	} else {
		res += string(j)
	}
	return res
}
