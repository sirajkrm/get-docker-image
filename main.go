package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	memory "github.com/go-git/go-git/v5/storage/memory"
)

// Data defines the structure for a Reposity
type data struct {
	ID         int      `json:"id"`
	RepoURL    string   `json:"repository"`
	Commit     string   `json:"commit"`
	Dockerfile string   `json:"dockerfile"`
	Image      []string `json:"image"`
}

// Repositories is a collection of Data which is a collection of main struct data
type Repositories struct {
	Data []*data
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//func to read a string and check for image and add it to Slice of string
func getImage(term, s string) []string {
	//match exactly FROM of Dockerfile, nothing before it (\\b), after it a bunch of characters \\W

	/**
	Scanner provides a convenient interface for reading data such as a file of newline-delimited lines of text.
	using bufio.Scanner in Go which iterates over lines from an io.Reader.
	**/
	scanner := bufio.NewScanner(strings.NewReader(s))
	var img []string
	for scanner.Scan() {
		matchedFrom, err := regexp.MatchString("\\b"+term+"\\W", scanner.Text())
		check(err)
		if matchedFrom {
			//scanner.Text() is each line of the string of Dockerfile
			i := strings.Split(scanner.Text(), " ")[1]

			//appending every image found to a slice of string
			img = append(img, i)
		}
	}

	if err := scanner.Err(); err != nil {
		check(err)
	}

	return img

}

//func to read a URL and convert it to string
func readURL(URL string) string {
	response, err := http.Get(URL)

	if err != nil {
		check(err)
	}

	defer response.Body.Close()

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, response.Body)
		out <- buf.String()
		check(err)
	}()

	urlData := <-out
	return urlData

}

func checkArgs() {

	f := flag.String("url", "", "need a URL to run againt it\nUsage: [-url] https:// ...")
	flag.Parse()

	if len(*f) == 0 {
		fmt.Println("expected a flag -url")
		flag.PrintDefaults()
		os.Exit(1)
	}

}

var (
	storer     *memory.Storage
	repository []*data
)

func main() {
	fs := memfs.New()

	var u string
	if len(os.Args) < 2 {
		checkArgs()
	} else {
		u = os.Args[2]
	}

	str := readURL(u)
	var (
		dockerfilePath string
		finalResult    []string
		id             int
	)

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		if scanner.Text() != "" {
			id++
			repoUrl := strings.Split(scanner.Text(), " ")[0]
			repoCommit := strings.Split(scanner.Text(), " ")[1]

			storer = memory.NewStorage()
			r, err := git.Clone(storer, fs, &git.CloneOptions{
				URL: repoUrl,
			})

			// this will only get you the full commit ID of 40 characters
			h, err := r.ResolveRevision(plumbing.Revision(repoCommit))
			check(err)

			//commit info here
			commit, err := r.CommitObject(*h)
			check(err)

			tree, err := commit.Tree()
			check(err)

			//this is where listing all files of Tree
			tree.Files().ForEach(func(f *object.File) error {

				//get Dockerfile exactly not before (\\b) not after $
				matched, err := regexp.MatchString("\\bDockerfile$", f.Name)
				check(err)

				if matched {
					dockerfilePath = f.Name
					file, err := fs.Open(f.Name)

					//file is of type int64, so making a channel to convert it to string by using buffer bytes
					out := make(chan string)
					go func() {
						var buf bytes.Buffer
						io.Copy(&buf, file)
						out <- buf.String()
						check(err)
					}()

					//the content of Dockerfile is string here
					strFrom := <-out

					//getting only the images (as a Slice of string) of that of Dockerfile
					finalResult = getImage("FROM", strFrom)

				}
				return nil
			})
			repository = append(repository, &data{ID: id, RepoURL: repoUrl, Commit: repoCommit, Dockerfile: dockerfilePath, Image: finalResult})
		}
	}

	//having the data presented in JSON format
	j, _ := json.Marshal(Repositories{Data: repository})
	fmt.Printf("%s\n", j)

}
