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
	api := dropbox.Client(token, options)

	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		fmt.Printf("%+v", err.Error())
	}

	for _, file := range files {
		srcPath := notesPath + file.Name()
		dstPath := remotePath + file.Name()
		upload(api, srcPath, dstPath)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func upload(client dropbox.Api, srcPath string, dstPath string) {
	contents, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		return
	}
	defer contents.Close()

	commitInfo := files.NewCommitInfo(dstPath)
	commitInfo.ClientModified = time.Now().UTC().Truncate(time.Second) // requires format '%Y-%m-%dT%H:%M:%SZ'
	commitInfo.Mode.Tag = "overwrite"                                  // dangerous!
	// commitInfo.Autorename = true                                    // set if change to 'add' or 'update'

	_, err = client.Upload(commitInfo, contents)
	if err != nil {
		fmt.Printf("✗ Error uploading %s to %s: %+v\n", srcPath, dstPath, err)
	} else {
		fmt.Printf("✓ Uploaded %s to %s\n", srcPath, dstPath)
	}
}
