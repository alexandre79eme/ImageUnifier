package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Structure used to saved the given parameters
type parameters struct {
	valid      bool
	shuffle    bool
	listMode   bool
	path       string
	outputFile string
}

func parseArguments(arguments []string) parameters {
	param := parameters{false, false, false, "", ""}
	pathSet := false
	outputSet := false

	for i := 1; i < len(arguments); i++ {
		if strings.HasPrefix(arguments[i], "-") {

			if strings.Contains(arguments[i], "r") {
				// If the option was already set then there is a input parameter problem
				if param.shuffle {
					break
				}
				param.shuffle = true
			}

			if strings.Contains(arguments[i], "l") {
				// If the option was already set then there is a input parameter problem
				if param.listMode {
					break
				}
				param.listMode = true
			}

		} else {
			if !pathSet {
				param.path = arguments[i]
				pathSet = true
				param.valid = !param.listMode
			} else if !outputSet && param.listMode {
				param.outputFile = arguments[i]
				outputSet = true
				param.valid = true
			} else {
				param.valid = false
			}

		}

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
		regex, _ := regexp.Compile("[0-9]")

		for k := 0; k < len(first); k++ {
			if k >= len(second) {
				return false
			}

			fChar := strings.ToLower(first[k])
			sChar := strings.ToLower(second[k])
			fCharInt := regex.MatchString(fChar)
			sCharInt := regex.MatchString(sChar)

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
	regex, _ := regexp.Compile("[0-9]")

	for i := start; i < len(array); i++ {
		charInt := regex.MatchString(array[i])

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

	res = strings.Replace(s, "[", "%5B", -1)
	res = strings.Replace(res, "]", "%5D", -1)

	return
}

// Open the first image of each folder successively
func interactiveImageReading(folders []imageFolder) {
	for _, folder := range folders {

		cmd := exec.Command("cmd.exe", "/C", folder.imageList[0])
		cmd.Dir = folder.path
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Type 'Enter' to continue: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
	}
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

	parameters := parseArguments(os.Args)
	if !parameters.valid {
		fmt.Println("Error: missing parameters")
		fmt.Println("USAGE: ImageUnifier [-rl] <folderPath> [<outputFile>]")
		return
	}

	folderList := listFilesInSubDir(parameters.path)

	if parameters.shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(folderList), func(i, j int) {
			folderList[i], folderList[j] = folderList[j], folderList[i]
		})
	}

	if parameters.listMode {
		fillFile(parameters.outputFile, folderList)
	} else {
		interactiveImageReading(folderList)
	}
}
