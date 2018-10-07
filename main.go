package main

import "fmt"
import "os"
import "path/filepath"
import "io"
import "time"
import "sort"
import "encoding/csv"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ", os.Args[0], "dirname")
		os.Exit(-1)
	}

	path, err := filepath.Abs(os.Args[1])
	checkAndExit(err, -1)

	pstat, err := os.Stat(path)
	checkAndExit(err, -1)
	if !pstat.IsDir() {
		fmt.Println(path, "is not a directory")
		os.Exit(0)
	}

	fmt.Println("Scanning", path, "please wait")

	startTime := time.Now()
	dir, err := os.Open(path)
	checkAndExit(err, -1)
	defer dir.Close()

	const n = 100
	m := make(map[int64]int64)
	for files, err := dir.Readdir(n); err != io.EOF; {
		checkAndExit(err, -1)
		for _, file := range files {
			modDateUnix := file.ModTime().UTC().Truncate(24 * time.Hour).Unix()
			if !file.Mode().IsDir() {
				m[modDateUnix] += file.Size()
			}
		}
		files, err = dir.Readdir(n)
	}
	writeScanDuration(startTime)
	writeScanResults(m)
}

func checkAndExit(err error, code int) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(code)
	}
}

func writeScanResults(m map[int64]int64) {
	var keys []int64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i int, j int) bool { return keys[i] < keys[j] })
	const gb = 1024 * 1024 * 1024
	var total float64
	outFilePath := "dstat.csv"
	outFile, err := os.Create("dstat.csv")
	checkAndExit(err, -2)
	defer outFile.Close()
	writer := csv.NewWriter(outFile)
	writer.Write([]string{"Date", "Bytes"})
	for _, key := range keys {
		val := m[key]
		total += float64(val)
		t := time.Unix(key, 0).UTC()
		writer.Write([]string{t.Format("2006-01-02"), fmt.Sprint(val)})
	}
	writer.Flush()
	fmt.Println("Stats written to", outFilePath)
	fmt.Println("Total size:", total/gb, "Gb")
}

func writeScanDuration(start time.Time) {
	scanTime := time.Since(start)
	fmt.Println("Scan completed in", scanTime)
}
