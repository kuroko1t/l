package main

import (
	"flag"
	"reflect"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	//"time"
	"sort"
)

type Color int

const (
	Red Color = iota
	Magenta
	Blue
	Green
	Cyan
	White
)

type Info struct {
	color Color
	fileinfo fs.FileInfo
}

var fileinfoMap map[string]Info

var detailFlag bool

func init() {
  flag.BoolVar(&detailFlag, "l", false, "show detail")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		dirwalk("./")
	} else {
		dirwalk(args[0])
	}
	PrintFiles()
}

func checkGitBranch(file fs.FileInfo, dirname string) bool {
	r, err := git.PlainOpen(dirname)
	if err != nil {
		return true
	}
	ref, err := r.Head()
	if err != nil {
		return true
	}
	printString := dirname + "(" + ref.Name().Short() + ")"
	fileinfoMap[printString] = Info{Blue, file}
	return false
}

func checkExt(file fs.FileInfo, filename string) {
	printString := filename
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		fileinfoMap[printString] = Info{Magenta, file}
	case ".jpg":
		fileinfoMap[printString] = Info{Magenta, file}
	case ".svg":
		fileinfoMap[printString] = Info{Magenta, file}
	case ".mp4":
		fileinfoMap[printString] = Info{Magenta, file}
	case ".mp3":
		fileinfoMap[printString] = Info{Magenta, file}
	case ".gz":
		fileinfoMap[printString] = Info{Magenta, file}
	case ".zip":
	default:
		fileinfoMap[printString] = Info{White, file}
	}
}

func fileModeCheck(file fs.FileInfo, filename string) bool {
	printString := filename
	if file.IsDir() {
		if checkGitBranch(file, filename) {
			fileinfoMap[printString] = Info{Blue, file}
		}
	} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		fileinfoMap[printString] = Info{Cyan, file}
	} else if file.Mode()&0100 == 0100 {
		fileinfoMap[printString] = Info{Green, file}
	} else {
		return true
	}
	return false
}

func dirwalk(dir string) []string {
	fileinfoMap = make(map[string]Info)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	var paths []string
	for _, file := range files {
		var filename string
		if strings.HasSuffix(dir, "/") {
			filename = dir + file.Name()
		} else {
			filename = dir + "/" + file.Name()
		}
		if fileModeCheck(file, filename) {
			checkExt(file, filename)
		}
	}
	return paths
}

func sortedKeys(mapInt interface{}) []string {
	values := reflect.ValueOf(mapInt).MapKeys()
	result := make([]string, len(values))
	for i, value1 := range values {
		result[i] = value1.String()
	}
	sort.Strings(result)
	return result
}

func PrintFiles() {
	for _, key := range sortedKeys(fileinfoMap) {
		printString := ""
		if detailFlag {
			fileModeString := fmt.Sprintf("%v", fileinfoMap[key].fileinfo.Mode())
			fileModeString += fmt.Sprintf(" %v", fileinfoMap[key].fileinfo.Size())
			fileModeString += fmt.Sprintf(" %v ", fileinfoMap[key].fileinfo.ModTime().Format("2006/01/02 15:04"))
			printString =  fileModeString + key
		} else {
			printString = key
		}
		switch fileinfoMap[key].color {
		case Blue:
			color.HiBlue(printString)
		case Magenta:
			color.Magenta(printString)
		case Green:
			color.Green(printString)
		case Cyan:
			color.Cyan(printString)
		default:
			fmt.Println(printString)
		}
	}
}
