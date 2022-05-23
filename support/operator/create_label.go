package main

import (
	"bufio"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Message the json output for label
type Message struct {
	Originalref string `json:originalref`
	FileName    string `json:"filename"`
	Speaker     string `json:"speaker"`
	Topic       string `json:"topic"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	Md5Sum      string `json:"md5sum"`
	Book        string `json:"book"`
	Chapter     int    `json:"chapter"`
	VerseStart  int    `json:"verseStart"`
	VerseEnd    int    `json:"verseEnd"`
	Part        int    `json:"part"`
}

func pullExisting() []Message {
	url := "https://helpersofyourjoy.com/format_mapping.json"
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: transCfg}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to get %s", url)
		os.Exit(1)
	}
	defer resp.Body.Close()
	log.Printf("getting %s\n", url)
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read %s", url)
		os.Exit(1)
	}
	result := []Message{}
	// fmt.Println(string(htmlData))
	err = json.Unmarshal(htmlData, &result)
	if err != nil {
		log.Printf("Failed to unmarshal %s: \n %s", url, err)
	}
	return result
}

func (m *Message) getSpeaker(existing *[]Message) {
	fmt.Println("Speaker: ")
	// Taking input from user
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	speakerCandidate := scanner.Text()
	if speakerCandidate == "" {
		log.Fatal("well who spoke then?")
	}
	// check that the speaker already exists..
	for _, v := range *existing {
		if v.Speaker == speakerCandidate {
			log.Printf("found speaker: %s\n", v.Speaker)
			m.Speaker = speakerCandidate
			return
		}
	}
	fmt.Printf("%s does not exist.. are you sure you want to create this speaker? [y/n] ", speakerCandidate)
	var answer string
	fmt.Scanf("%s", &answer)
	fmt.Println("answer: ", answer)
	if strings.ToLower(answer) == "y" {
		fmt.Println("ok you got it")
		m.Speaker = speakerCandidate
	} else {
		fmt.Println("aborting")
		os.Exit(0)
	}

}

func (m *Message) setFileName(existing *[]Message) {
	fmt.Println("filepath: ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	filepath := s.Text()
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Fatalf("%s not exists", filepath)
	}
	if info.IsDir() {
		log.Fatalf("%s is a directory not a file.", filepath)
	}
	// check if duplicate
	// break down to just filename
	log.Println("checking for duplicate filename")
	basepath := strings.Split(filepath, "/")
	for _, v := range *existing {
		if v.FileName == basepath[len(basepath)-1] {
			log.Fatalf("%s is duplicate filename, change filename and try again", basepath)
		}
	}
	// md5sum the file
	m.FileName = filepath

}

func (m *Message) setMD5Sum() {
	log.Println("running setMD5Sum")
	f, err := os.Open(m.FileName)
	if err != nil {
		log.Fatal("Failed to read (want to md5sum) %s", m.FileName)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatalf("Failed md5sum %s", m.FileName)
	}
	m.Md5Sum = fmt.Sprintf("%x", h.Sum(nil))
	// check if the md5 already exists??
}

func (m *Message) setTopic(existing *[]Message) {
	log.Println("running setTopic")
	// Taking input from user
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	topicCandidate := scanner.Text()
	if topicCandidate == "" {
		log.Fatal("topic name required, even if its \"misc\"")
	}
	// check that the topic already exists..
	for _, v := range *existing {
		if v.Topic == topicCandidate {
			log.Printf("found topic: %s\n", v.Topic)
			m.Topic = topicCandidate
			return
		}
	}
	fmt.Printf("%s does not exist.. are you sure you want to create this topic? [y/n] ", topicCandidate)
	var answer string
	fmt.Scanf("%s", &answer)
	fmt.Println("answer: ", answer)
	if strings.ToLower(answer) == "y" {
		fmt.Println("ok you got it")
		m.Topic = topicCandidate
	} else {
		fmt.Println("aborting")
		os.Exit(0)
	}
}

func (m *Message) setYear() {
	log.Println("running setYear")
	t := time.Now()
	year := t.Year()
	// Taking input from user
	fmt.Printf("If %d ok, press enter, else enter the year: ", year)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	yearCandidate := scanner.Text()
	if yearCandidate == "" {
		log.Printf("using year: %s\n", year)
		m.Year = year
		return
	}
	holder, err := strconv.Atoi(yearCandidate)
	if err != nil {
		log.Fatal("%s is not a year")
	}
	if holder > 2050 {
		log.Fatal("lol we havent been raptured by 2050?")
	}
	if holder < 1984 {
		log.Fatal("Thats too old? right?")
	}
	m.Year = holder

}

func (m *Message) setMonth() {
	log.Println("running setMonth")
	t := time.Now()
	month := t.Month()
	// Taking input from user
	fmt.Printf("If %d ok, press enter, else enter the month: ", month)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	monthCandidate := scanner.Text()
	if monthCandidate == "" {
		log.Printf("using month: %s\n", month)
		m.Month = int(t.Month())
		return
	}
	holder, err := strconv.Atoi(monthCandidate)
	if err != nil {
		log.Fatalf("%s is not a month")
	}
	if holder > 12 || holder < 1 {
		log.Fatalf("invalid month: %d", holder)
	}
	m.Month = holder
}

func (m *Message) setDay() {
	log.Println("running setDay")
	t := time.Now()
	day := t.Day()
	// Taking input from user
	fmt.Printf("If %d ok, press enter, else enter the day: ", day)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dayCandidate := scanner.Text()
	if dayCandidate == "" {
		log.Printf("using day: %s\n", day)
		m.Day = int(t.Day())
		return
	}
	holder, err := strconv.Atoi(dayCandidate)
	if err != nil {
		log.Fatalf("%s is not a day")
	}
	if holder > 31 || holder < 1 {
		log.Fatalf("invalid day: %d", holder)
	}
	m.Day = holder
}

func main() {
	existing := pullExisting()
	m := Message{}
	// m.setFileName(&existing)
	// m.setMD5Sum()
	// m.getSpeaker(&existing)
	// m.setTopic(&existing)
	_ = existing
	m.setYear()
	m.setMonth()
	m.setDay()
	fmt.Printf("final result: %#v\n", m)
}
