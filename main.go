package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Structure used to save the folders with images
// I save the folder path and the list of the images name
type imageFolder struct {
	path      string
	imageList []string
}

// Return a new imageFolder
func newImageFolder(path string) *imageFolder {
	return &imageFolder{path, make([]string, 0)}
}

// Report the error if needed
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Test if the file is an image
func isImageFile(extension string) bool {
	switch extension {
	case "bmp", "jpg", "jpeg", "gif", "png":
		return true
	}
	return false
}

// Return the list of sub-directory with images.
func listFilesInSubDir(path string) []imageFolder {
	var result []imageFolder

	files, err := ioutil.ReadDir(path)
	check(err)

	folder := newImageFolder(path)

	for _, f := range files {
		if f.IsDir() {
			result = append(result, listFilesInSubDir(path+"/"+f.Name())...)
		} else {
			nameSplitted := strings.Split(f.Name(), ".")
			extension := nameSplitted[len(nameSplitted)-1]

			if isImageFile(extension) {
				folder.imageList = append(folder.imageList, f.Name())
			}

		}
	}

	//fmt.Println(folder)
	if len(folder.imageList) != 0 {
		result = append(result, *folder)
	}

	return result
}

// Write the path of each image of all sub-directories in the file
func fillFile(fileName string, folders []imageFolder) {
	f, err := os.Create(fileName)
	check(err)

	for _, folder := range folders {
		path := folder.path + "/"

		for _, imageName := range folder.imageList {
			f.WriteString(path + imageName + "\n")
		}
	}

}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Error: missing parameters")
		fmt.Println("USAGE: ImageUnifier <folderPath> <outputFile>")
		return
	}

	path := os.Args[1]
	fileName := os.Args[2]

	folderList := listFilesInSubDir(path)

	fillFile(fileName, folderList)
}
