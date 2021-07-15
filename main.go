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
		color.Blue(filename)
	} else if file.Mode() &os.ModeSymlink == os.ModeSymlink {
		color.Cyan(filename)
	} else if file.Mode() & 0100 == 0100 {
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
