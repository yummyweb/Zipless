package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Invalid arguements. Must be 2 args.")
	}

	src := os.Args[1]
	dest := os.Args[2]

    r, err := zip.OpenReader(src)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            log.Fatal(err)
        }
        defer rc.Close()

        fpath := filepath.Join(dest, f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, f.Mode())
        } else {
            var fdir string
            if lastIndex := strings.LastIndex(fpath,string(os.PathSeparator)); lastIndex > -1 {
                fdir = fpath[:lastIndex]
            }

            err = os.MkdirAll(fdir, f.Mode())
            if err != nil {
                log.Fatal(err)
            }
            f, err := os.OpenFile(
                fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                log.Fatal(err)
            }
            defer f.Close()

            _, err = io.Copy(f, rc)
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}