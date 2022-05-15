package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type reference struct {
	book         string
	startChapter string
	endChapter   string
	startVerse   string
	endVerse     string
}

type mediaData struct {
	date    string
	subject string
	title   string
	speaker string
	partNum string
	refs    []reference
}

func (m *mediaData) createFileName() string {
	return ""
}

func (m *mediaData) clean() {
	m.date = strings.TrimSuffix(m.date, "\n")
	m.date = strings.ToLower(m.date)
	m.subject = strings.TrimSuffix(m.subject, "\n")
	m.subject = strings.ToLower(m.subject)
	m.title = strings.TrimSuffix(m.title, "\n")
	m.title = strings.ToLower(m.title)
	m.speaker = strings.TrimSuffix(m.speaker, "\n")
	m.speaker = strings.ToLower(m.speaker)
	m.partNum = strings.TrimSuffix(m.partNum, "\n")
	m.partNum = strings.ToLower(m.partNum)
}

func (r *reference) clean() {
	r.book = strings.TrimSuffix(r.book, "\n")
	r.book = strings.ToLower(r.book)
	r.startChapter = strings.TrimSuffix(r.startChapter, "\n")
	r.endChapter = strings.TrimSuffix(r.endChapter, "\n")
	r.startVerse = strings.TrimSuffix(r.startVerse, "\n")
	r.endVerse = strings.TrimSuffix(r.endVerse, "\n")
}

func (r *reference) sanitize() {
	r.book = strings.ReplaceAll(r.book, " ", "_")
}

func createRef() *reference {
	reader := bufio.NewReader(os.Stdin)
	r := reference{}
	// book
	fmt.Printf("book name: ")
	r.book, _ = reader.ReadString('\n')
	// startChapter
	fmt.Printf("starting chapter: ")
	r.startChapter, _ = reader.ReadString('\n')
	// endChapter
	fmt.Printf("ending chapter: ")
	r.endChapter, _ = reader.ReadString('\n')
	// startVerse
	fmt.Printf("starting verse: ")
	r.startVerse, _ = reader.ReadString('\n')
	// endVerse
	fmt.Printf("ending verse: ")
	r.endVerse, _ = reader.ReadString('\n')

	r.clean()
	return &r
}

// replace space for _
func (m *mediaData) sanitize() {
	m.subject = strings.ReplaceAll(m.subject, " ", "_")
	m.title = strings.ReplaceAll(m.title, " ", "_")
	m.speaker = strings.ReplaceAll(m.speaker, " ", "_")
}

func (m *mediaData) makeFileName() string {
	result := ""
	// date
	result = result + fmt.Sprintf("date.%s", m.date)
	// series
	if m.subject != "" {
		result = result + fmt.Sprintf("-series.%s", m.subject)
	} else {
		result = result + fmt.Sprintf("-series.null")
	}
	if m.speaker != "" {
		result = result + fmt.Sprintf("-speaker.%s", m.speaker)
	} else {
		result = result + fmt.Sprintf("-speaker.null")
	}
	if m.partNum != "" {
		result = result + fmt.Sprintf("-part.%s", m.partNum)
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
	for i := 0; i < len(m.refs); i++ {
		refs = refs + fmt.Sprintf("book.%s", m.refs[i].book)

		if m.refs[i].startChapter != "" {
			_, err := strconv.Atoi(m.refs[i].startChapter)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-startchapter.%s", m.refs[i].startChapter)
		}

		if m.refs[i].endChapter != "" {
			_, err := strconv.Atoi(m.refs[i].endChapter)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-startchapter.%s", m.refs[i].endChapter)
		}

		if m.refs[i].startVerse != "" {
			_, err := strconv.Atoi(m.refs[i].startVerse)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-startverse.%s", m.refs[i].startVerse)
		}

		if m.refs[i].endVerse != "" {
			_, err := strconv.Atoi(m.refs[i].endVerse)
			if err != nil {
				panic(err)
			}
			refs = refs + fmt.Sprintf("-endverse.%s", m.refs[i].endVerse)
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
	m.date, _ = reader.ReadString('\n')
	// Series
	fmt.Printf("series name? %s", skipStr)
	m.subject, _ = reader.ReadString('\n')
	// Speaker
	fmt.Printf("speaker name: %s", skipStr)
	m.speaker, _ = reader.ReadString('\n')
	// Title
	fmt.Printf("title name?: ")
	m.title, _ = reader.ReadString('\n')
	// Part number
	fmt.Printf("part number?: %s ", skipStr)
	m.partNum, _ = reader.ReadString('\n')
	m.clean()
	m.sanitize()
	// scripture reference
	fmt.Printf("Add book references?: %s ", skipStr)
	addRefs, _ := reader.ReadString('\n')
	fmt.Println("addRefs:", addRefs)
	if addRefs == "\n" {
		fmt.Println("all done")
	} else {
		m.refs = append(m.refs, *createRef())
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

	fmt.Println("EOF")
}
