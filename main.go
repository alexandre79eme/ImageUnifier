package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"sort"
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

	sort.Slice(folder.imageList, func(i, j int) bool {
		first := strings.Split(folder.imageList[i], "")
		second := strings.Split(folder.imageList[j], "")

		for k := 0; k < len(first); k++ {
			if k >= len(second) {
				return false
			}

			fChar := strings.ToLower(first[k])
			sChar := strings.ToLower(second[k])
			fCharInt, _ := regexp.MatchString("[0-9]", fChar)
			sCharInt, _ := regexp.MatchString("[0-9]", sChar)

			if fCharInt && sCharInt {
				fCount := countSuccessiveInt(first, k)
				sCount := countSuccessiveInt(second, k)
				if fCount != sCount && fChar != "0" && sChar != "0" {
					return fCount < sCount
				} else if fCount != sCount && (fChar == "0" || sChar == "0") {
					return fChar == "0"
				}
			}

			if fChar != sChar {
				return fChar < sChar
			}

		}
		return true
	})

	//fmt.Println(folder)
	if len(folder.imageList) != 0 {
		result = append(result, *folder)
	}

	return result
}

// Count the count of successive integer starting by the choosen start included
func countSuccessiveInt(array []string, start int) (count int) {
	for i := start; i < len(array); i++ {
		charInt, _ := regexp.MatchString("[0-9]", array[i])

		if charInt {
			count++
		} else {
			break
		}
	}

	return
}

// Replace some characters by their asci hexa equivalent. Useful for programs that don't read
// some special characters.
func replaceUnsupportedCharacter(s string) (res string) {
	res = s
	strings.Replace(res, "[", "%5B", -1)
	strings.Replace(res, "]", "%5D", -1)

	return
}

// Write the path of each image of all sub-directories in the file
func fillFile(fileName string, folders []imageFolder) {
	f, err := os.Create(fileName)
	check(err)

	for _, folder := range folders {
		path := replaceUnsupportedCharacter(folder.path + "/")

		for _, imageName := range folder.imageList {
			f.WriteString(path + replaceUnsupportedCharacter(imageName) + "\n")
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
