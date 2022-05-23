package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
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
			break
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

func main() {
	existing := pullExisting()
	m := Message{}
	m.getSpeaker(&existing)
	fmt.Printf("final result: %#v\n", m)
}
