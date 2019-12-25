package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"bufio"
	"bytes"
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
		fmt.Println(destDir)
		if err := os.Mkdir(destDir, 0644); err != nil {
			log.Fatal(err)
		}
	}

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	countFileLines := func(fileName, sourceDir, destDir string) {
		data, err := os.Open(filepath.Join(sourceDir,fileName))

		if err != nil {
			log.Fatal(err)
		}
		defer data.Close()

		reader := bufio.NewReader(data)
		buffer := bytes.NewBuffer(make([]byte, 0))
		part := make([]byte, chunksize)
		var count int
		for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}

	if err != io.EOF {
		log.Fatal("Error Reading ", fileName, ": ", err)
	} else {
		err = nil
	}
		lineCount := buffer.Len()

		newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".res"

		err = ioutil.WriteFile(filepath.Join(destDir, newFileName), []byte(strconv.Itoa(lineCount)), 0644)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}

	for _, file := range files {
		wg.Add(1)
		go countFileLines(file.Name(), sourceDir, destDir)
	}

	wg.Wait()
	fmt.Println("Total number of processed files:", len(files))
}
