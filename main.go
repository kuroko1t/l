package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	//"github.com/djherbis/times"
)

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
	color.HiBlue(printString)
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
		color.Magenta(printString)
	case ".jpg":
		color.Magenta(printString)
	case ".svg":
		color.Magenta(printString)
	case ".mp4":
		color.Magenta(printString)
	case ".mp3":
		color.Magenta(filename)
	case ".gz":
		color.Red(printString)
	case ".zip":
		color.Red(printString)
	default:
		fmt.Println(printString)
	}
}

func fileModeCheck(file fs.FileInfo, filename string) bool {
	printString := ""
	if detailFlag {
		//fileModeString := fmt.Sprintf("%v", file.Mode())
		//printString = fileModeString + " " + filename
		fileModeString := fmt.Sprintf("%v", file.Mode())
		fileModeString += fmt.Sprintf(" %v", file.Size())
		fileModeString += fmt.Sprintf(" %v ", file.ModTime().Format("2006/01/02 15:04"))
		printString = fileModeString + filename
	} else {
		printString = filename
	}
	if file.IsDir() {
		if checkGitBranch(file, filename) {
			color.HiBlue(printString)
		}
	} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		color.Cyan(printString)
	} else if file.Mode()&0100 == 0100 {
		color.Green(printString)
	} else {
		return true
	}
	return false
}

func dirwalk(dir string) []string {
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
