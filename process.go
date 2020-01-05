package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	// "bytes"
	"io"
)

const (
	chunksize int = 10
)

func main() {
	sourceDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	destDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.Mkdir(destDir, 0777); err != nil {
			log.Fatal(err)
		}
	}

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	countFileLines := func(fileName, sourceDir, destDir string) {
		defer wg.Done()
		data, err := os.Open(filepath.Join(sourceDir, fileName))
		if err != nil {
			log.Fatal(err)
		}

		defer data.Close()
		var sentence string

		reader := bufio.NewReader(data)
		buffer := make([]byte, chunksize)
		var count int
		for {
			if count, err = reader.Read(buffer); err != nil {
				break
			}
			sentence += string(buffer[:count])
		}

		if err != io.EOF {
			log.Fatal("Error Reading ", fileName, ": ", err)
		} else {
			err = nil
		}

		newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".res"
		f, err := os.Create(filepath.Join(destDir, newFileName))

		if err != nil {
			log.Fatal("create file ", err)
		}

		defer f.Close()
		if _, err = f.WriteString(sentence); err != nil {
			log.Fatal(err)
		}
	}

	for _, file := range files {
		wg.Add(1)
		go countFileLines(file.Name(), sourceDir, destDir)
	}

	wg.Wait()
	fmt.Println("Total number of processed files:", len(files))
}
