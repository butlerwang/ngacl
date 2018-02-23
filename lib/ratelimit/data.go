/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : data.go

* Purpose :

* Creation Date : 03-23-2017

* Last Modified : Tue 16 May 2017 09:48:58 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package ratelimit

import (
	"encoding/json"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Data struct {
	Meter  map[Host]map[Ip]int64
	start  time.Time
	Uptime time.Duration
	lock   *sync.RWMutex
	Count  uint64
}

type Host string
type Ip string

type ReqChan struct {
	*http.Request
	av int64
}

var (
	data *Data
	ch   chan *http.Request
)

func cleanMapBlackList() {
	data.lock.Lock()
	var list []Ip
	for host := range data.Meter {
		for k := range data.Meter[host] {
			data.Meter[host][k] += int64(-5)
			if data.Meter[host][k] < 0 {
				list = append(list, k)
			}
		}
		for _, ip := range list {
			delete(data.Meter[host], ip)
		}
		list = []Ip{}
	}
	data.lock.Unlock()
}

func server() {
	r := mux.NewRouter()
	r.HandleFunc("/list", indexHandler)
	r.HandleFunc("/list/{host}", hostHandler)
	http.ListenAndServe(":8106", gziphandler.GzipHandler(r))
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

func writeData(r chan *http.Request) {
	for val := range r {
		data.Count++
		h := Host(val.Host)
		ip := val.RemoteAddr
		if v, ok := val.Context().Value("xff0").(string); ok {
			ip = v
		}
		i := Ip(ip)

		_, ok := data.Meter[h]
		if !ok {
			m := make(map[Ip]int64)
			data.lock.Lock()
			data.Meter[h] = m
			data.lock.Unlock()
		}

		data.lock.Lock()
		data.Meter[h][i] += 1
		data.lock.Unlock()
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

// /list/host.com
func hostHandler(w http.ResponseWriter, r *http.Request) {
	val := mux.Vars(r)
	h := val["host"]
	host := Host(h)
	html := NewHtml()

	html.Body = toJson(data.Meter[host])

	html.Title = h
	html.Back = "/list"
	t, _ := template.New("index").Parse(IndexTemplate)
	t.Execute(w, html)
}
