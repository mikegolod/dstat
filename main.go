package main

import "fmt"
import "os"
import "path/filepath"
import "io"
import "time"
import "sort"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ", os.Args[0], "dirname")
	}

	path, err := filepath.Abs(os.Args[1])
	checkAndExit(err, -1)

	pstat, err := os.Stat(path)
	checkAndExit(err, -1)
	if !pstat.IsDir() {
		fmt.Println(path, "is not a directory")
		os.Exit(0)
	}

	fmt.Println("Scanning", path, "please wait...")

	dir, err := os.Open(path)
	checkAndExit(err, -1)
	defer dir.Close()

	const n = 100
	m := make(map[int64]int64)
	for files, err := dir.Readdir(n); err != io.EOF; {
		checkAndExit(err, -1)
		for _, file := range files {
			modDateUnix := file.ModTime().UTC().Round(24 * time.Hour).Unix()
			if !file.Mode().IsDir() {
				m[modDateUnix] += file.Size()
			}
		}
		files, err = dir.Readdir(n)
	}
	out(m)
}

func checkAndExit(err error, code int) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(code)
	}
}

func out(m map[int64]int64) {
	var keys []int64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i int, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		val := m[key]
		fmt.Println(time.Unix(key, 0).UTC().Format("2006-01-02"), val)
	}
}
