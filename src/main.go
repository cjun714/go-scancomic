package main

import (
	"archive/zip"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"webp"
)

func listFiles(file *zip.File) error {
	rd, e := file.Open()
	if e != nil {
		msg := "Failed to open zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, e)
	}
	defer rd.Close()

	fi := file.FileInfo()
	fmt.Fprintf(os.Stdout, "%s, isdir:%t", file.Name, fi.IsDir())

	if e != nil {
		msg := "Failed to read zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, e)
	}

	fmt.Println()

	return nil
}

func main() {
	path := "//192.168.2.106/movie/mamic-e/Zatanna by Paul Dini (2017) (Digital) (Son of Ultron-Empire).cbz"
	e := getCover(path, "z:/")
	if e != nil {
		panic(e)
	}
}

func list(path string) {
	rd, e := zip.OpenReader(path)
	if e != nil {
		log.Fatalf("Failed to open: %s", e)
	}
	defer rd.Close()

	names := make([]string, len(rd.File))
	i := 0
	for _, f := range rd.File {
		names[i] = f.Name
		i++
	}

	sort.Strings(names)
	fmt.Println(names[0])
}

func convert() {
	src := os.Args[1]
	target := os.Args[2]
	quality, _ := strconv.Atoi(os.Args[3])
	scale, _ := strconv.ParseFloat(os.Args[4], 32)

	e := webp.ToWEBP(src, target, quality, float32(scale))
	if e != nil {
		panic(e)
	}

}

func getCover(src, targetDir string) error {
	rd, e := zip.OpenReader(src)
	if e != nil {
		return e
	}
	defer rd.Close()

	names := make([]string, len(rd.File))
	files := make(map[string]*zip.File, len(rd.File))
	i := 0
	for _, f := range rd.File {
		if f.FileInfo().IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue
		}

		names[i] = f.Name
		files[f.Name] = f
		i++
	}

	sort.Strings(names)
	ext := filepath.Ext(names[0])
	coverName := filepath.Base(src)
	coverName = strings.Replace(coverName, ".cbz", ext, 1)
	coverPath := filepath.Join(targetDir, coverName)

	f := files[names[0]]
	r, e := f.Open()
	if e != nil {
		return e
	}
	wr, e := os.Create(coverPath)
	if e != nil {
		return e
	}
	defer wr.Close()
	_, e = io.Copy(wr, r)
	if e != nil {
		return e
	}

	return nil
}
