package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Structure used to saved the given parameters
type parameters struct {
	shuffle    bool
	path       string
	outputFile string
}

func parseArguments(arguments []string) parameters {
	param := parameters{false, "", ""}

	if strings.HasPrefix(arguments[1], "-") {
		if arguments[1] == "-r" {
			param.shuffle = true
		}

		param.path = arguments[2]
		param.outputFile = arguments[3]
	} else {
		param.path = arguments[1]
		param.outputFile = arguments[2]
	}

	return param
}

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
		fmt.Println("USAGE: ImageUnifier [-r] <folderPath> <outputFile>")
		return
	}

	parameters := parseArguments(os.Args)

	folderList := listFilesInSubDir(parameters.path)

	if parameters.shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(folderList), func(i, j int) {
			folderList[i], folderList[j] = folderList[j], folderList[i]
		})
	}

	fillFile(parameters.outputFile, folderList)
}
