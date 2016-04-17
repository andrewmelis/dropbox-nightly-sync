package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/dropbox/dropbox-sdk-go-unofficial/files"
)

const (
	token = "insert token here"
	notesPath  = "/home/andrew/Notes/"
	remotePath = "/nightly-dropbox-sync/"
)

func main() {
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
}

func upload(client dropbox.Api, srcPath string, dstPath string) {
	contents, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		return
	}
	defer contents.Close()

	commitInfo := files.NewCommitInfo(dstPath)
	metadata, err := client.Upload(commitInfo, contents)
	if err != nil {
		fmt.Printf("✗ Error uploading %s to %s: %+v\n", srcPath, dstPath, err)
	} else {
		fmt.Printf("✓ Uploaded %s to %s:\n%+v\n", srcPath, dstPath, metadata)
	}
}

func listRemoteFiles(client dropbox.Api, path string) {
	arg := files.NewListFolderArg(path)

	resp, err := client.ListFolder(arg)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		return
	}

	for i, entry := range resp.Entries {
		fmt.Printf("%d: type %T\n %+v\n\n", i, entry.File, entry.File)
	}
}

