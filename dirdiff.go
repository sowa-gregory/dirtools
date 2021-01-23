package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

type fileent struct {
	name string
	size int64
	hash int
}

type dirent struct {
	path       string
	totalFiles int
	dirs       []dirent
	files      []fileent
}

func scanDir(path string) (dirent, error) {
	currentDir := dirent{path, 0, nil, nil}
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return currentDir, err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			subDir, err := scanDir(path + file.Name() + "/")
			if err != nil {
				return currentDir, err
			}
			currentDir.dirs = append(currentDir.dirs, subDir)
			currentDir.totalFiles += subDir.totalFiles
		} else {
			currentFile := fileent{file.Name(), file.Size(), 0}
			currentDir.files = append(currentDir.files, currentFile)
			currentDir.totalFiles++
		}
	}
	sort.Slice(currentDir.files, func(i, j int) bool {
		return currentDir.files[i].name < currentDir.files[j].name
	})

	return currentDir, nil
}

func printDir(dir dirent) {
	for _, file := range dir.files {
		fmt.Printf("%s\n", dir.path+"/"+file.name)
	}
	for _, dir := range dir.dirs {
		printDir(dir)
	}
}

func mergeSlices(dst []string, src []string) []string {
	for _, item := range src {
		dst = append(dst, item)
	}
	return dst
}

func findFiles(dir dirent, re *regexp.Regexp) []string {
	var foundFiles []string
	for _, file := range dir.files {
		fullPath := dir.path + "/" + file.name

		if re.MatchString(fullPath) {
			foundFiles = append(foundFiles, fullPath)
		}
	}
	for _, dir := range dir.dirs {
		foundFiles = mergeSlices(foundFiles, findFiles(dir, re))
	}
	if re.MatchString(dir.path) {
		foundFiles = append(foundFiles, dir.path)
	}
	return foundFiles
}

func removeFiles(files []string) {
	for _, path := range files {
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("Cannot delete file:%s\n", path)
		}
	}
}

func main() {
	//exclude := "*/node_modules/*"
	re := regexp.MustCompile(".*/node_modules/.*")
	fmt.Println(re.MatchString("asd/node_modules/asdf"))
	root := "/Users/sowisz/Programming/"

	dirs, err := scanDir(root)
	if err != nil {
		panic(err)
	}
	//printDir(dirs)
	//fmt.Printf("Total files:%d\n", dirs.totalFiles)
	o := findFiles(dirs, re)
	for _, s := range o {
		fmt.Println(s)
	}

	removeFiles(o)
}
