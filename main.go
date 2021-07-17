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
	"time"
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
	modTime time.Time
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
	printString := ""
	if detailFlag {
		fileModeString := fmt.Sprintf("%v", file.Mode())
		fileModeString += fmt.Sprintf(" %v", file.Size())
		fileModeString += fmt.Sprintf(" %v ", file.ModTime().Format("2006/01/02 15:04"))
		printString = fileModeString + " " + dirname + "(" + ref.Name().Short() + ")"
	} else {
		printString = dirname + "(" + ref.Name().Short() + ")"
	}
	fileinfoMap[printString] = Info{Blue, file.ModTime()}
	return false
}

func checkExt(file fs.FileInfo, filename string) {
	printString := ""
	if detailFlag {
		fileModeString := fmt.Sprintf("%v", file.Mode())
		fileModeString += fmt.Sprintf(" %v", file.Size())
		fileModeString += fmt.Sprintf(" %v ", file.ModTime().Format("2006/01/02 15:04"))
		printString = fileModeString + filename
	} else {
		printString = filename
	}
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		fileinfoMap[printString] = Info{Magenta, file.ModTime()}
	case ".jpg":
		fileinfoMap[printString] = Info{Magenta, file.ModTime()}
	case ".svg":
		fileinfoMap[printString] = Info{Magenta, file.ModTime()}
	case ".mp4":
		fileinfoMap[printString] = Info{Magenta, file.ModTime()}
	case ".mp3":
		fileinfoMap[printString] = Info{Magenta, file.ModTime()}
	case ".gz":
		fileinfoMap[printString] = Info{Magenta, file.ModTime()}
	case ".zip":
	default:
		fileinfoMap[printString] = Info{White, file.ModTime()}
	}
}

func fileModeCheck(file fs.FileInfo, filename string) bool {
	printString := ""
	if detailFlag {
		fileModeString := fmt.Sprintf("%v", file.Mode())
		fileModeString += fmt.Sprintf(" %v", file.Size())
		fileModeString += fmt.Sprintf(" %v ", file.ModTime().Format("2006/01/02 15:04"))
		printString = fileModeString + filename
	} else {
		printString = filename
	}
	if file.IsDir() {
		if checkGitBranch(file, filename) {
			fileinfoMap[printString] = Info{Blue, file.ModTime()}
		}
	} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		fileinfoMap[printString] = Info{Cyan, file.ModTime()}
	} else if file.Mode()&0100 == 0100 {
		fileinfoMap[printString] = Info{Green, file.ModTime()}
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
		switch fileinfoMap[key].color {
		case Blue:
			color.HiBlue(key)
		case Magenta:
			color.Magenta(key)
		case Green:
			color.Green(key)
		case Cyan:
			color.Cyan(key)
		default:
			fmt.Println(key)
		}
	}
}
