package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/go-redis/redis"
)

var (
	redisSetName = "HOYJ::MP3::MAP::DUMP"
)

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

// main: run the search engine
func main() {
	port := "127.0.0.1:8082"
	http.HandleFunc("/mp3/search", mp3Search)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
