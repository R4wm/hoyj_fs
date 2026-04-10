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
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Name}} - Helpers of Your Joy</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      min-height: 100vh;
      padding: 20px;
    }
    .wrapper { max-width: 800px; margin: 0 auto; }
    header { text-align: center; margin-bottom: 30px; }
    header h1 { color: white; font-size: 2rem; font-weight: 300; text-shadow: 0 2px 4px rgba(0,0,0,0.2); }
    header a { color: rgba(255,255,255,0.9); text-decoration: none; }
    header a:hover { text-decoration: underline; }
    .content-box {
      background: white;
      border-radius: 12px;
      padding: 30px;
      box-shadow: 0 10px 40px rgba(0,0,0,0.2);
    }
    .back-link {
      display: inline-block;
      margin-bottom: 20px;
      color: #667eea;
      text-decoration: none;
      font-weight: 500;
    }
    .back-link:hover { text-decoration: underline; }
    .file-list { list-style: none; }
    .file-list li {
      padding: 14px 16px;
      border-radius: 8px;
      margin-bottom: 8px;
      background: #f8f9fa;
      transition: all 0.2s;
    }
    .file-list li:hover { background: #e9ecef; transform: translateX(4px); }
    .file-list a {
      color: #333;
      text-decoration: none;
      display: block;
      word-break: break-word;
    }
    .file-list a:hover { color: #667eea; }
    footer { text-align: center; margin-top: 30px; color: rgba(255,255,255,0.8); font-size: 0.9rem; }
    footer a { color: white; text-decoration: none; font-weight: 500; }
    footer a:hover { text-decoration: underline; }
  </style>
</head>
<body>
  <div class="wrapper">
    <header>
      <h1><a href="https://helpersofyourjoy.com">Helpers of Your Joy</a></h1>
    </header>
    <div class="content-box">
      <a href="/media" class="back-link">← Back to Series</a>
      <h2 style="margin-bottom: 20px; color: #333;">{{ .Name }}</h2>
      <ul class="file-list">
        {{ range $index, $results := .Series }}
        <li><a href="{{ $results }}">{{ createLink $results }}</a></li>
        {{ end }}
      </ul>
    </div>
    <footer>
      <a href="/mp3/form.html">Search MP3 Library</a> · <a href="/media">Browse Series</a>
    </footer>
  </div>
</body>
</html>
`

const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}} - Helpers of Your Joy</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      min-height: 100vh;
      padding: 20px;
    }
    .wrapper { max-width: 800px; margin: 0 auto; }
    header { text-align: center; margin-bottom: 30px; }
    header h1 { color: white; font-size: 2rem; font-weight: 300; text-shadow: 0 2px 4px rgba(0,0,0,0.2); }
    header a { color: rgba(255,255,255,0.9); text-decoration: none; }
    .content-box {
      background: white;
      border-radius: 12px;
      padding: 30px;
      box-shadow: 0 10px 40px rgba(0,0,0,0.2);
    }
    .back-link {
      display: inline-block;
      margin-bottom: 20px;
      color: #667eea;
      text-decoration: none;
      font-weight: 500;
    }
    .back-link:hover { text-decoration: underline; }
    .results-count { color: #666; margin-bottom: 16px; }
    .file-list { list-style: none; max-height: 600px; overflow-y: auto; }
    .file-list li {
      padding: 14px 16px;
      border-radius: 8px;
      margin-bottom: 8px;
      background: #f8f9fa;
      transition: all 0.2s;
    }
    .file-list li:hover { background: #e9ecef; transform: translateX(4px); }
    .file-list a {
      color: #333;
      text-decoration: none;
      display: block;
      word-break: break-word;
    }
    .file-list a:hover { color: #667eea; }
    .no-results { text-align: center; padding: 40px; color: #888; }
    footer { text-align: center; margin-top: 30px; color: rgba(255,255,255,0.8); font-size: 0.9rem; }
    footer a { color: white; text-decoration: none; font-weight: 500; }
  </style>
</head>
<body>
  <div class="wrapper">
    <header>
      <h1><a href="https://helpersofyourjoy.com">Helpers of Your Joy</a></h1>
    </header>
    <div class="content-box">
      <a href="/mp3/form.html" class="back-link">← New Search</a>
      <h2 style="margin-bottom: 10px; color: #333;">Search Results</h2>
      {{if .Items}}
      <p class="results-count">Found {{len .Items}} result{{if ne (len .Items) 1}}s{{end}}</p>
      <ul class="file-list">
        {{range .Items}}<li><a href="{{ . }}" target="_blank">{{ basenameMP3Files . }}</a></li>{{end}}
      </ul>
      {{else}}
      <div class="no-results">No results found. Try different keywords.</div>
      {{end}}
    </div>
    <footer>
      <a href="/mp3/form.html">Search</a> · <a href="/media">Browse Series</a>
    </footer>
  </div>
</body>
</html>`

const chapterButtonsTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Media Library - Helpers of Your Joy</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      min-height: 100vh;
      padding: 20px;
    }
    .wrapper { max-width: 900px; margin: 0 auto; }
    header { text-align: center; margin-bottom: 30px; }
    header h1 { color: white; font-size: 2rem; font-weight: 300; text-shadow: 0 2px 4px rgba(0,0,0,0.2); }
    header a { color: rgba(255,255,255,0.9); text-decoration: none; }
    header a:hover { text-decoration: underline; }
    .search-bar {
      background: white;
      border-radius: 12px;
      padding: 20px 30px;
      margin-bottom: 20px;
      box-shadow: 0 10px 40px rgba(0,0,0,0.2);
      display: flex;
      gap: 12px;
      align-items: center;
    }
    .search-bar input {
      flex: 1;
      padding: 12px 16px;
      font-size: 1rem;
      border: 2px solid #e0e0e0;
      border-radius: 8px;
      outline: none;
    }
    .search-bar input:focus { border-color: #667eea; }
    .search-bar a {
      padding: 12px 24px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      text-decoration: none;
      border-radius: 8px;
      font-weight: 600;
      white-space: nowrap;
    }
    .search-bar a:hover { opacity: 0.9; }
    .content-box {
      background: white;
      border-radius: 12px;
      padding: 30px;
      box-shadow: 0 10px 40px rgba(0,0,0,0.2);
    }
    .series-count { color: #666; margin-bottom: 20px; }
    .series-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
      gap: 12px;
      list-style: none;
    }
    .series-grid li a {
      display: block;
      padding: 16px 20px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      text-decoration: none;
      border-radius: 8px;
      font-weight: 500;
      text-align: center;
      transition: all 0.2s;
      font-size: 0.9rem;
    }
    .series-grid li a:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
    }
    footer { text-align: center; margin-top: 30px; color: rgba(255,255,255,0.8); font-size: 0.9rem; }
    footer a { color: white; text-decoration: none; font-weight: 500; }
    footer a:hover { text-decoration: underline; }
    @media (max-width: 600px) {
      .search-bar { flex-direction: column; }
      .search-bar a { width: 100%; text-align: center; }
      .series-grid { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <div class="wrapper">
    <header>
      <h1><a href="https://helpersofyourjoy.com">Helpers of Your Joy</a></h1>
    </header>
    <div class="search-bar">
      <input type="text" id="filterInput" placeholder="Filter series..." onkeyup="filterSeries()">
      <a href="/mp3/form.html">Search MP3s</a>
    </div>
    <div class="content-box">
      <h2 style="margin-bottom: 10px; color: #333;">Media Library</h2>
      <p class="series-count" id="seriesCount">{{ len .Series }} series available</p>
      <ul class="series-grid" id="seriesList">
        {{ range $index, $results := .Series }}
        <li><a href="{{ add "./media/" $results }}">{{ $results }}</a></li>
        {{ end }}
      </ul>
    </div>
    <footer>
      <a href="https://helpersofyourjoy.com">Home</a> · <a href="/mp3/form.html">Search</a>
    </footer>
  </div>
  <script>
    function filterSeries() {
      const filter = document.getElementById('filterInput').value.toLowerCase();
      const items = document.querySelectorAll('#seriesList li');
      let visible = 0;
      items.forEach(item => {
        const text = item.textContent.toLowerCase();
        if (text.includes(filter)) {
          item.style.display = '';
          visible++;
        } else {
          item.style.display = 'none';
        }
      });
      document.getElementById('seriesCount').textContent = visible + ' series' + (filter ? ' matching "' + filter + '"' : ' available');
    }
  </script>
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
	// search function returning html with result
	r.HandleFunc("/mp3/search", mp3Search)
	// html page of categories to select from
	r.HandleFunc("/media", listMedia)
	// json output listing all resources
	r.HandleFunc("/media/info", mediaInfo)
	// returns html page with resources per category selected
	r.HandleFunc("/media/{category}", topicDirectory)
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
