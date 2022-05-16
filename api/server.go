package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var (
	redisSetName = "HOYJ::MP3::MAP::DUMP"
)

const chapterTemplate = `
<html>
<title>{{.Name}}</title>
<style>
.btn-group button {
  background-color: gold; /* Green background */
  border: 1px solid green; /* Green border */
  color: black;
  padding: 10px 24px; /* Some padding */
  cursor: pointer; /* Pointer/hand icon */
  float: center; /* Float the buttons side by side */
}
/* Clear floats (clearfix hack) */
.btn-group:after {
  content: "";
  clear: both;
  display: table;
}
.btn-group button:not(:last-child) {
  border-right: none; /* Prevent double borders */
}
/* Add a background color on hover */
.btn-group button:hover {
  background-color: #3e8e41;
}
</style>
  <body style="background-color:{{ .Color }};">
    <h1><center><a href=https://helpersofyourjoy.com/media>{{ .Name }}</a></h1>
  <body>
    {{ range $index, $results := .Series }}
    <p><b><left><a href={{ $results }}> {{ createLink $results }}</a>  </b></p>
    {{ end }}
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
    <div class="w3-bar">
    <div class="btn-group">
    </div>
  </body>
</html>
`

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<p><a href="{{ . }}">{{ basenameMP3Files . }}</a></p>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`

const chapterButtonsTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      .block {
      display: block;
      width: 100%;
      border: none;
      background-color: #4CAF50;
      color: white;
      padding: 14px 28px;
      font-size: 16px;
      cursor: pointer;
      text-align: center;
      }
      .block:hover {
      background-color: #ddd;
      color: black;
      }
    </style>
    <title>{{ .Name }}</title>
  </head>
  <body style="background-color:{{ .Color }};">
    <p><center><h1> {{ .Name }} </h1><center></p>
    <button onclick="window.location.href='https://helpersofyourjoy.com';" class="w3-bar-item w3-button" style="width:33.3%">Return to helpersofyourjoy.com</button>
    {{ range $index, $results := .Series }}
    <p><button class="block" onclick="window.location.href = '{{ add "./directory/" $results }}'">{{ $results }}</button></p>
    {{ end }}
  </body>
</html>
`

func sanitizeArgs(args string) []string {
	parsedArgs := strings.Split(args, " ")
	for i, val := range parsedArgs {
		parsedArgs[i] = strings.TrimSpace(strings.ToLower(val))
		// remove trailing whitespace
	}
	return parsedArgs
}

// mp3Search: search pre populated redis for files based on search strings
func mp3Search(w http.ResponseWriter, r *http.Request) {
	var relevantFiles []string

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	// TODO: write issues back to html page here
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "PraseForm() err: %v\n", err)
	}

	queryStr := r.FormValue("q")
	if len(queryStr) <= 0 {
		w.Write([]byte("i think your missing search argument.."))
		return
	}
	// Testing
	parsedArgs := sanitizeArgs(queryStr)
	fmt.Println("parsed Args: ", parsedArgs)

	fmt.Printf("typeQuery: %s\n", queryStr)
	// Get everything from redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	val, err := rdb.SMembers(redisSetName).Result()
	if err != nil {
		fmt.Println("Failed to parse redis: ", err)
	}
	sort.Strings(val)

	// sanitize args: remove trailing newline, whitespace,
	for _, v := range val {
		ok := true
		// for _, vv := range args {
		for _, vv := range parsedArgs {
			if !strings.Contains(strings.ToLower(v), strings.ToLower(vv)) {
				ok = false
				break
			}
		}
		// Passed all checks
		if ok {
			relevantFiles = append(relevantFiles, v)
		}
	}
	sort.Strings(relevantFiles)

	// func to create basename so results are super long
	funcs := template.FuncMap{"basenameMP3Files": func(mp3File string) string {
		parsedPaths := strings.Split(mp3File, "/")
		return parsedPaths[len(parsedPaths)-1]
	}}

	t, err := template.New("webpage").Funcs(funcs).Parse(tpl)
	check(err)

	mp3Payload := struct {
		Title string
		Items []string
	}{
		Title: "mp3 search results",
		Items: relevantFiles,
	}

	err = t.Execute(w, mp3Payload)
	check(err)
}

func listMedia(w http.ResponseWriter, r *http.Request) {

	fmt.Println("listing series")
	options := redis.Options{}
	redisCli := redis.NewClient(&options)

	result := redisCli.Keys("HOYJ::MEDIA::*").Val()
	sort.Strings(result)
	seriesAvailable := []string{}
	for _, seriesTopic := range result {
		seriesTopic = strings.SplitAfterN(seriesTopic, "::", 3)[2]
		seriesAvailable = append(seriesAvailable, seriesTopic)
	}
	w.WriteHeader(http.StatusOK)

	funcs := template.FuncMap{"add": func(x, y string) string { return x + y }}
	t, err := template.New("seriesList").Funcs(funcs).Parse(chapterButtonsTemplate)
	if err != nil {
		panic(err)
	}

	Listing := struct {
		Name   string
		Series []string
		Color  string
	}{"Series", seriesAvailable, "White"}

	t.Execute(w, Listing)

}

func mediaInfo(w http.ResponseWriter, r *http.Request) {
	// Get everything from redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	fmt.Println(redisSetName)
	val, err := rdb.SMembers(redisSetName).Result()
	if err != nil {
		fmt.Println("Failed to parse redis: ", err)
	}
	sort.Strings(val)
	haha, err := json.Marshal(val)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("failed: %s", err)))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(haha)
}

func topicDirectory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	options := redis.Options{}
	redisCli := redis.NewClient(&options)
	// get keys from redis for topic
	vars["category"] = strings.ToUpper(vars["category"])
	redisKey := "HOYJ::MEDIA::" + vars["category"]
	fmt.Println("redisKey: ", redisKey)

	redisResult := redisCli.SMembers(redisKey).Val()
	sort.Strings(redisResult)
	// fmt.Println("this is redis result: ", redisResult)

	funcs := template.FuncMap{"createLink": func(b string) string {
		breakit := strings.Split(b, "/")
		return breakit[len(breakit)-1]
	}}

	t, err := template.New("seriesLinks").Funcs(funcs).Parse(chapterTemplate)
	if err != nil {
		panic(err)
	}

	Listing := struct {
		Name   string
		Series []string
		Color  string
	}{vars["category"], redisResult, "White"}

	t.Execute(w, Listing)
}

// main: run the search engine
func main() {
	addr := "127.0.0.1:8082"
	r := mux.NewRouter()
	// returns html page with resources per category selected
	r.HandleFunc("/directory/{category}", topicDirectory)
	// html page of categories to select from
	r.HandleFunc("/directory", listMedia)
	// search function returning html with result
	r.HandleFunc("/mp3/search", mp3Search)
	// json output listing all resources
	r.HandleFunc("/directory/info", mediaInfo)
	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("running server")
	log.Fatal(srv.ListenAndServe())
}
