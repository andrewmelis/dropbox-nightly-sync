package main

import (
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/dropbox/dropbox-sdk-go-unofficial/files"
)

const (
	token = "insert token here"
)

func main() {
	var options dropbox.Options
	api := dropbox.Client(token, options)

	fmt.Printf("\nList files once...\n")
	listFiles(api, "/nightly-dropbox-sync")

	fmt.Printf("\nUpload...\n")
	upload(api, "./test2.txt", "/nightly-dropbox-sync/test2.txt")

	fmt.Printf("\nList files twice...\n")
	listFiles(api, "/nightly-dropbox-sync")
}

func listFiles(client dropbox.Api, path string) {
	arg := files.NewListFolderArg(path)

	resp, err := client.ListFolder(arg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	for i, entry := range resp.Entries {
		fmt.Printf("%d: type %T\n %+v\n\n", i, entry.File, entry.File)
	}
}

func upload(client dropbox.Api, srcPath string, dstPath string) {
	contents, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	defer contents.Close()

	commitInfo := files.NewCommitInfo(dstPath)
	// commitInfo.Mode.Tag = "overwrite"

	metadata, err := client.Upload(commitInfo, contents)
	if err != nil {
		fmt.Printf("type: %T\n %+v\n", err, err)
	} else {
		fmt.Printf("type: %T\n %+v\n", metadata, metadata)
	}
}
