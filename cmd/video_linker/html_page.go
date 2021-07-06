package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Video struct {
	Link  string
	Title string
}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func create_link(directory, file string) Video {
	url := "https://helpersofyourjoy.com/media"
	result := strings.Join([]string{url, directory, file}, "/")
	return Video{
		Link:  result,
		Title: "something",
	}
}

// ignore returns true if file should be ignored
func ignore(file string) bool {
	ignore_list := []string{"#", "~"}
	for _, v := range ignore_list {
		if strings.Contains(file, v) {
			return true
		}
	}
	return false
}

func main() {
	pwd, err := os.Getwd()
	check_error(err)
	fmt.Println(pwd)

	// get files in the directory
	fileInfo, err := ioutil.ReadDir(pwd)
	check_error(err)

	directory_broken := strings.Split(pwd, "/")
	directory := directory_broken[len(directory_broken)-2]

	var files []string
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	videos := []Video{}

	for _, v := range files {
		fmt.Println("v: ", v)
		if ignore(v) {
			fmt.Printf("skipping: %s\n", v)
			continue
		}
		videoObj := create_link(directory, v)
		videos = append(videos, videoObj)
	}

	fmt.Println(videos)
	fmt.Println("End of program")
}
