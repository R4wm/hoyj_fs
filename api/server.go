package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/go-redis/redis"
)

var ctx = context.Background

func mp3Search(w http.ResponseWriter, r *http.Request) {
	var redisSetName = "HOYJ::MP3::MAP::DUMP"
	var relevantFiles []string

	fmt.Println("MP3 Searching..")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "PraseForm() err: %v\n", err)
	}

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

	args := []string{"verstegen", "mystery", "CROSS"}
	// sanitize args: remove trailing newline, whitespace,
	for _, v := range val {
		ok := true
		for _, vv := range args {
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

	funcs := template.FuncMap{"basenameMP3Files": func(mp3File string) string {
		parsedPaths := strings.Split(mp3File, "/")
		return parsedPaths[len(parsedPaths)-1]
		return "hi"
	}}

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Funcs(funcs).Parse(tpl)
	check(err)

	noItems := struct {
		Title string
		Items []string
	}{
		Title: "mp3 search results",
		Items: relevantFiles,
	}

	err = t.Execute(w, noItems)
	check(err)
}

func main() {

	port := "127.0.0.1:8082"

	http.HandleFunc("/mp3/search", mp3Search)

	fmt.Printf("Starting web server")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
