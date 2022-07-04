package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type reference struct {
	Book         string `json:"book"`
	StartChapter string `json:"startChapter"`
	EndChapter   string `json:"endChapter"`
	StartVerse   string `json:"startVerse"`
	EndVerse     string `json:"endVerse"`
}

type mediaData struct {
	Date    string      `json:"date"`
	Subject string      `json:"subject"`
	Title   string      `json:"title"`
	Speaker string      `json:"speaker"`
	PartNum string      `json:"partNum"`
	Refs    []reference `json:"refs"`
}

func (m *mediaData) createFileName() string {
	return ""
}

func (m *mediaData) clean() {
	m.Date = strings.TrimSuffix(m.Date, "\n")
	m.Date = strings.ToLower(m.Date)
	m.Subject = strings.TrimSuffix(m.Subject, "\n")
	m.Subject = strings.ToLower(m.Subject)
	m.Title = strings.TrimSuffix(m.Title, "\n")
	m.Title = strings.ToLower(m.Title)
	m.Speaker = strings.TrimSuffix(m.Speaker, "\n")
	m.Speaker = strings.ToLower(m.Speaker)
	m.PartNum = strings.TrimSuffix(m.PartNum, "\n")
	m.PartNum = strings.ToLower(m.PartNum)
}

func (r *reference) clean() {
	r.Book = strings.TrimSuffix(r.Book, "\n")
	r.Book = strings.ToLower(r.Book)
	r.StartChapter = strings.TrimSuffix(r.StartChapter, "\n")
	r.EndChapter = strings.TrimSuffix(r.EndChapter, "\n")
	r.StartVerse = strings.TrimSuffix(r.StartVerse, "\n")
	r.EndVerse = strings.TrimSuffix(r.EndVerse, "\n")
}

func (r *reference) sanitize() {
	r.Book = strings.ReplaceAll(r.Book, " ", "_")
}

func createRef() *reference {
	reader := bufio.NewReader(os.Stdin)
	r := reference{}
	// book
	fmt.Printf("book name: ")
	r.Book, _ = reader.ReadString('\n')
	// startChapter
	fmt.Printf("starting chapter: ")
	r.StartChapter, _ = reader.ReadString('\n')
	// endChapter
	fmt.Printf("ending chapter: ")
	r.EndChapter, _ = reader.ReadString('\n')
	// startVerse
	fmt.Printf("starting verse: ")
	r.StartVerse, _ = reader.ReadString('\n')
	// endVerse
	fmt.Printf("ending verse: ")
	r.EndVerse, _ = reader.ReadString('\n')

	r.clean()
	return &r
}

// replace space for _
func (m *mediaData) sanitize() {
	m.Subject = strings.ReplaceAll(m.Subject, " ", "_")
	m.Title = strings.ReplaceAll(m.Title, " ", "_")
	m.Speaker = strings.ReplaceAll(m.Speaker, " ", "_")
}

func (m *mediaData) makeFileName() string {
	result := ""
	// date
	result = result + fmt.Sprintf("date.%s", m.Date)
	// series
	if m.Subject != "" {
		result = result + fmt.Sprintf("-series.%s", m.Subject)
	} else {
		result = result + fmt.Sprintf("-series.null")
	}
	if m.Speaker != "" {
		result = result + fmt.Sprintf("-speaker.%s", m.Speaker)
	} else {
		result = result + fmt.Sprintf("-speaker.null")
	}
	if m.PartNum != "" {
		result = result + fmt.Sprintf("-part.%s", m.PartNum)
	} else {
		result = result + fmt.Sprintf("-part.null")
	}

	return result
}

func (m *mediaData) addRefToFileName(fileName string) string {
	// book         string
	// startChapter string
	// endChapter   string
	// startVerse   string
	// endVerse     string
	refs := ".--"
	for i := 0; i < len(m.Refs); i++ {
		refs = refs + fmt.Sprintf("book.%s", m.Refs[i].Book)

		if m.Refs[i].StartChapter != "" {
			_, err := strconv.Atoi(m.Refs[i].StartChapter)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-startchapter.%s", m.Refs[i].StartChapter)
		}

		if m.Refs[i].EndChapter != "" {
			_, err := strconv.Atoi(m.Refs[i].EndChapter)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-startchapter.%s", m.Refs[i].EndChapter)
		}

		if m.Refs[i].StartVerse != "" {
			_, err := strconv.Atoi(m.Refs[i].StartVerse)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-startverse.%s", m.Refs[i].StartVerse)
		}

		if m.Refs[i].EndVerse != "" {
			_, err := strconv.Atoi(m.Refs[i].EndVerse)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-endverse.%s", m.Refs[i].EndVerse)
		}
	}
	// double period is our end marker
	refs = refs + "--."
	fileName = fileName + refs
	return fileName
}

func getFileType(fileName string) string {
	broken := strings.Split(fileName, ".")
	fileType := broken[len(broken)-1]
	fmt.Println("fileType: ", fileType)
	return fileType

}
func main() {
	reader := bufio.NewReader(os.Stdin)
	skipStr := "(\"Enter\" to skip) "
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	m := mediaData{}
	// Date
	fmt.Println("date format: YYYYMMDD")
	m.Date, _ = reader.ReadString('\n')
	// Series
	fmt.Printf("series name? %s", skipStr)
	m.Subject, _ = reader.ReadString('\n')
	// Speaker
	fmt.Printf("speaker name: %s", skipStr)
	m.Speaker, _ = reader.ReadString('\n')
	// Title
	fmt.Printf("title name?: ")
	m.Title, _ = reader.ReadString('\n')
	// Part number
	fmt.Printf("part number?: %s ", skipStr)
	m.PartNum, _ = reader.ReadString('\n')
	m.clean()
	m.sanitize()
	// scripture reference
	fmt.Printf("Add book references?: %s ", skipStr)
	addRefs, _ := reader.ReadString('\n')
	fmt.Println("addRefs:", addRefs)
	if addRefs == "\n" {
		fmt.Println("all done")
	} else {
		m.Refs = append(m.Refs, *createRef())
	}

	fmt.Println("location of the media file?: ")
	originalFileLocation, _ := reader.ReadString('\n')
	originalFileLocation = strings.TrimSuffix(originalFileLocation, "\n")
	fmt.Printf("mediaData: %#v\n", m)
	fileName := m.makeFileName()
	if addRefs != "\n" {
		fileName = m.addRefToFileName(fileName)
	}

	fmt.Println("filename: ", fileName)
	fileType := getFileType(originalFileLocation)
	fmt.Printf("fileType: -%s-\n", fileType)
	fileName = fileName + fileType

	fmt.Println("getFileType: ", getFileType(originalFileLocation))
	// rename the file
	dest := filepath.Dir(originalFileLocation) + "/" + fileName
	dest = filepath.FromSlash(dest)

	fmt.Println("dest: ", dest)
	os.Rename(originalFileLocation, dest)

	// json output
	jsonData, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
	fmt.Println("EOF")
}
