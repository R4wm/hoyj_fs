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
	"strings"
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
	speakerFound := false
	for _, v := range *existing {
		if v.Speaker == speakerCandidate {
			log.Printf("found speaker: %s\n", v.Speaker)
			speakerFound = true
			m.Speaker = speakerCandidate
			return
		}
	}
	if !speakerFound {
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

func main() {
	existing := pullExisting()
	m := Message{}
	m.setFileName(&existing)
	m.setMD5Sum()
	m.getSpeaker(&existing)

	fmt.Printf("final result: %#v\n", m)
}
