package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/dropbox/dropbox-sdk-go-unofficial/files"
)

const (
	token      = "insert token here"
	notesPath  = "/home/andrew/Notes/"
	remotePath = "/nightly-dropbox-sync/"
)

func main() {
	start := time.Now()

	var options dropbox.Options
	client := dropbox.Client(token, options)

	fileNameCh := make(chan string)
	outputCh := make(chan string)

	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		fmt.Printf("%+v", err.Error())
	}

	for i := 0; i < 10; i++ {
		go uploadWorker(i, client, fileNameCh, outputCh)
	}

	for _, file := range files {
		fileNameCh <- file.Name()
	}

	for i := 0; i < len(files); i++ {
		fmt.Printf("%s\n", <-outputCh)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func uploadWorker(i int, client dropbox.Api, fileNameCh chan string, outputCh chan string) {
	for {
		fileName := <-fileNameCh

		srcPath := notesPath + fileName
		dstPath := remotePath + fileName

		err := upload(client, srcPath, dstPath)
		if err != nil {
			outputCh <- fmt.Sprintf("%d: ✗ Error uploading %s to %s: %+v\n", i, srcPath, dstPath, err)
		} else {
			outputCh <- fmt.Sprintf("%d: ✓ Uploaded %s to %s\n", i, srcPath, dstPath)
		}
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func upload(client dropbox.Api, srcPath string, dstPath string) (err error) {
	contents, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer contents.Close()

	commitInfo := files.NewCommitInfo(dstPath)
	commitInfo.ClientModified = time.Now().UTC().Truncate(time.Second) // requires format '%Y-%m-%dT%H:%M:%SZ'
	commitInfo.Mode.Tag = "overwrite"                                  // dangerous!
	// commitInfo.Autorename = true                                    // set if change to 'add' or 'update'

	_, err = client.Upload(commitInfo, contents)
	return
}
