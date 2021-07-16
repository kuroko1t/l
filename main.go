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
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		dirwalk("./")
	} else {
		dirwalk(args[0])
	}
}

func checkGitBranch(dirname string) bool {
	r, err := git.PlainOpen(dirname)
	if err != nil {
		return true
	}
	ref, err := r.Head()
	if err != nil {
		return true
	}
	color.HiBlue(dirname + "(" + ref.Name().Short() + ")")
	return false
}

func checkPrint(filename string) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		color.Magenta(filename)
	case ".jpg":
		color.Magenta(filename)
	case ".svg":
		color.Magenta(filename)
	case ".mp4":
		color.Magenta(filename)
	case ".mp3":
		color.Magenta(filename)
	case ".gz":
		color.Red(filename)
	case ".zip":
		color.Red(filename)
	default:
		fmt.Println(filename)
	}
}

func fileModeCheck(file fs.FileInfo, filename string) bool {
	if file.IsDir() {
		if checkGitBranch(filename) {
			color.HiBlue(filename)
		}
	} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		color.Cyan(filename)
	} else if file.Mode()&0100 == 0100 {
		color.Green(filename)
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
			checkPrint(filename)
		}
	}
	return paths
}
