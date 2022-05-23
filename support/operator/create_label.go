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

var orderedBooks = []string{
	"GENESIS",
	"EXODUS",
	"LEVITICUS",
	"NUMBERS",
	"DEUTERONOMY",
	"JOSHUA",
	"JUDGES",
	"RUTH",
	"1SAMUEL",
	"2SAMUEL",
	"1KINGS",
	"2KINGS",
	"1CHRONICLES",
	"2CHRONICLES",
	"EZRA",
	"NEHEMIAH",
	"ESTHER",
	"JOB",
	"PSALMS",
	"PROVERBS",
	"ECCLESIASTES",
	"SONG OF SOLOMON",
	"ISAIAH",
	"JEREMIAH",
	"LAMENTATIONS",
	"EZEKIEL",
	"DANIEL",
	"HOSEA",
	"JOEL",
	"AMOS",
	"OBADIAH",
	"JONAH",
	"MICAH",
	"NAHUM",
	"HABAKKUK",
	"ZEPHANIAH",
	"HAGGAI",
	"ZECHARIAH",
	"MALACHI",
	"MATTHEW",
	"MARK",
	"LUKE",
	"JOHN",
	"ACTS",
	"ROMANS",
	"1CORINTHIANS",
	"2CORINTHIANS",
	"GALATIANS",
	"EPHESIANS",
	"PHILIPPIANS",
	"COLOSSIANS",
	"1THESSALONIANS",
	"2THESSALONIANS",
	"1TIMOTHY",
	"2TIMOTHY",
	"TITUS",
	"PHILEMON",
	"HEBREWS",
	"JAMES",
	"1PETER",
	"2PETER",
	"1JOHN",
	"2JOHN",
	"3JOHN",
	"JUDE",
	"REVELATION",
}

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
	m.Originalref = s.Text()
	filepathbreak := strings.Split(m.Originalref, "/")
	filepath := filepathbreak[len(filepathbreak)-1]
	info, err := os.Stat(m.Originalref)
	if os.IsNotExist(err) {
		log.Fatalf("%s not exists", m.Originalref)
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
	f, err := os.Open(m.Originalref)
	if err != nil {
		log.Fatalf("Failed to read (want to md5sum) %s", m.Originalref)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatalf("Failed md5sum %s", m.Originalref)
	}
	m.Md5Sum = fmt.Sprintf("%x", h.Sum(nil))
	// check if the md5 already exists??
}

func (m *Message) setTopic(existing *[]Message) {
	log.Println("running setTopic")
	// Taking input from user
	fmt.Println("topic: ")
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

func (m *Message) setBook() {
	log.Println("running setBook")
	// Taking input from user
	fmt.Printf("Book: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	bookCandidate := scanner.Text()
	for _, v := range orderedBooks {
		if strings.ToUpper(bookCandidate) == v {
			m.Book = strings.ToLower(v)
			return
		}
	}
	log.Fatalf("book: %s does not exist, these are available: %s",
		bookCandidate, orderedBooks)
}

func (m *Message) setChapter() {
	log.Println("running setChapter")
	// Taking input from user
	fmt.Printf("Chapter: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	chapterCandidate := scanner.Text()
	if chapterCandidate == "" {
		m.Chapter = 0
		return
	}
	t, err := strconv.Atoi(chapterCandidate)
	if err != nil {
		log.Fatalf("Failed to set chapter: %s", chapterCandidate)
	}
	m.Chapter = t
}

func (m *Message) setVerseStart() {
	log.Println("running setVerseStart")
	// Taking input from user
	fmt.Println("Verse Start: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	verseStartCandidate := scanner.Text()
	if verseStartCandidate == "" {
		m.VerseStart = 0
		return
	}
	t, err := strconv.Atoi(verseStartCandidate)
	if err != nil {
		log.Fatalf("Failed to set verseStart: %s", verseStartCandidate)
	}
	m.VerseStart = t
}

func (m *Message) setVerseEnd() {
	log.Println("running setVerseEnd")
	// Taking input from user
	fmt.Printf("VerseEnd: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	verseEndCandidate := scanner.Text()
	if verseEndCandidate == "" {
		m.VerseEnd = 0
		return
	}
	t, err := strconv.Atoi(verseEndCandidate)
	if err != nil {
		log.Fatalf("Failed to set verseEnd: %s", verseEndCandidate)
	}
	m.VerseEnd = t
}

func (m *Message) setPart() {
	log.Println("running setPart")
	// Taking input from user
	fmt.Println("Part number: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	PartCandidate := scanner.Text()
	if PartCandidate == "" {
		m.Part = 0
		return
	}
	t, err := strconv.Atoi(PartCandidate)
	if err != nil {
		log.Fatalf("Failed to set Part: %s", PartCandidate)
	}
	m.Part = t
}

func (m *Message) writeIt() {
	log.Printf("Writing %s\n", m.FileName+".json")
	file, _ := json.MarshalIndent(m, "", "")
	_ = ioutil.WriteFile(m.FileName+".json", file, 0644)
}

func main() {
	existing := pullExisting()
	m := Message{}
	m.setFileName(&existing)
	m.setMD5Sum()
	m.getSpeaker(&existing)
	m.setTopic(&existing)

	m.setYear()
	m.setMonth()
	m.setDay()
	m.setBook()
	m.setChapter()
	m.setVerseStart()
	m.setVerseEnd()
	m.setPart()
	m.writeIt()
	fmt.Printf("final result: %#v\n", m)
}
