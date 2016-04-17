package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/dropbox/dropbox-sdk-go-unofficial/files"
)

const (
	token = "insert token here"
	notesPath = "/home/andrew/Notes/"
)

func main() {

	
	files, err := retrieveLocalFileList(notesPath)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	for _, file := range files {
		fmt.Printf("%s\n", file)
	}

	// var options dropbox.Options
	// api := dropbox.Client(token, options)

	// fmt.Printf("\nList remote files once...\n")
	// listRemoteFiles(api, "/nightly-dropbox-sync")

	// fmt.Printf("\nUpload...\n")
	// upload(api, "./test.txt", "/nightly-dropbox-sync/test.txt")
	// upload(api, "./test2.txt", "/nightly-dropbox-sync/test2.txt")
	// upload(api, "./test3.txt", "/nightly-dropbox-sync/test3.txt")

	// fmt.Printf("\nList remote files twice...\n")
	// listRemoteFiles(api, "/nightly-dropbox-sync")
}

func retrieveLocalFileList(dir string) (fileList []string, err error) {
	fileList = make([]string, 5) // what is the best starting length?
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	return
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

func upload(client dropbox.Api, srcPath string, dstPath string) {
	contents, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		return
	}
	defer contents.Close()

	commitInfo := files.NewCommitInfo(dstPath)

	//WIP options below

	// commitInfo.ClientModified = time.Now() // FIXME  Error in call to API function "files/upload": HTTP header "Dropbox-API-Arg": client_modified: time data '2016-04-16T17:02:37.15034098-05:00' does not match format '%Y-%m-%dT%H:%M:%SZ'

	// commitInfo.Mode.Tag = "update"
	// commitInfo.Mode.Update = "49c1a0a8e89e1" // FIXME when figuring out what's new
	// commitInfo.Autorename = true

	// commitInfo.Mode.Tag = "overwrite"
	// commitInfo.Autorename = true

	metadata, err := client.Upload(commitInfo, contents)
	if err != nil {
		fmt.Printf("✗ Error uploading %s to %s: %+v\n", srcPath, dstPath, err)
	} else {
		fmt.Printf("✓ Uploaded %s to %s:\n%+v\n", srcPath, dstPath, metadata)
	}
}
