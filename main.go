package main

import (
	"flag"
	//"reflect"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	color       Color
	fileinfo    fs.FileInfo
	printString string
}

var allFileInfo []Info

type ByTime struct {
	infos []Info
}

func (by ByTime) Len() int {
	return len(by.infos)
}

func (by ByTime) Swap(i, j int) {
	by.infos[i], by.infos[j] = by.infos[j], by.infos[i]
}

func (by ByTime) Less(i, j int) bool {
	return by.infos[i].fileinfo.ModTime().Unix() < by.infos[j].fileinfo.ModTime().Unix()
}

type ByStr struct {
	infos []Info
}

func (by ByStr) Len() int {
	return len(by.infos)
}

func (by ByStr) Swap(i, j int) {
	by.infos[i], by.infos[j] = by.infos[j], by.infos[i]
}

func (by ByStr) Less(i, j int) bool {
	return strings.ToLower(by.infos[i].printString) < strings.ToLower(by.infos[j].printString)
}

//var fileinfoMap map[string]Info

var detailFlag bool

func main() {
	flag.BoolVar(&detailFlag, "l", false, "show detail")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		dirwalk("./")
	} else {
		dirwalk(args[0])
	}
	if detailFlag {
		PrintDetailFiles()
	} else {
		PrintFiles()
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
	printString := dirname + "(" + ref.Name().Short() + ")"
	allFileInfo = append(allFileInfo, Info{Blue, file, printString})
	return false
}

func checkExt(file fs.FileInfo, filename string) {
	printString := filename
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	case ".jpg":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	case ".svg":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	case ".mp4":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	case ".mp3":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	case ".gz":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	case ".zip":
		allFileInfo = append(allFileInfo, Info{Magenta, file, printString})
	default:
		allFileInfo = append(allFileInfo, Info{White, file, printString})
	}
}

func fileModeCheck(file fs.FileInfo, filename string) bool {
	printString := filename
	if file.IsDir() {
		if checkGitBranch(file, filename) {
			allFileInfo = append(allFileInfo, Info{Blue, file, printString})
		}
	} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		allFileInfo = append(allFileInfo, Info{Cyan, file, printString})
	} else if file.Mode()&0100 == 0100 {
		allFileInfo = append(allFileInfo, Info{Green, file, printString})
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

func PrintFiles() {
	sort.Sort(ByStr{allFileInfo})
	for _, info := range allFileInfo {
		switch info.color {
		case Blue:
			color.HiBlue(info.printString)
		case Magenta:
			color.Magenta(info.printString)
		case Green:
			color.Green(info.printString)
		case Cyan:
			color.Cyan(info.printString)
		default:
			fmt.Println(info.printString)
		}
	}
}

func PrintDetailFiles() {
	sort.Sort(ByTime{allFileInfo})
	for _, info := range allFileInfo {
		fileModeString := fmt.Sprintf("%v", info.fileinfo.Mode())
		fileModeString += fmt.Sprintf(" %v", info.fileinfo.Size())
		fileModeString += fmt.Sprintf(" %v ", info.fileinfo.ModTime().Format("2006/01/02 15:04"))
		switch info.color {
		case Blue:
			color.HiBlue(fileModeString + info.printString)
		case Magenta:
			color.Magenta(fileModeString + info.printString)
		case Green:
			color.Green(fileModeString + info.printString)
		case Cyan:
			color.Cyan(fileModeString + info.printString)
		default:
			fmt.Println(fileModeString + info.printString)
		}
	}
}
