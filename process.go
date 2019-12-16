package main

import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
    "strings"
    "strconv"
    "sync"
)

func main() {
    args := os.Args
    sourceDir, err := filepath.Abs(args[1])
    if err != nil {
        fmt.Printf("error: %v\n", err)
        os.Exit(0);
    }
    destDir, err := filepath.Abs(args[2])
    if err != nil {
        fmt.Printf("error: %v\n", err)
        os.Exit(0);
    }

    if _, err := os.Stat(destDir); os.IsNotExist(err) {
        fmt.Println(destDir)
        os.Mkdir(destDir, 0777)
    }

    files, err := ioutil.ReadDir(sourceDir)
    if err != nil {
        fmt.Printf("error: %v\n", err)
        os.Exit(0);
    }

    var wg sync.WaitGroup

    countFileLines := func(fileName, sourceDir, destDir string, wg *sync.WaitGroup) {
        data, err := ioutil.ReadFile(filepath.Join(sourceDir,fileName))
        if err != nil {
            fmt.Printf("error: %v\n", err)
            os.Exit(0);
        }
        lineCount := 1
        for _, chank:= range data {
            if chank == 10 {
                lineCount++
            }
        }

        newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".res"

        ioutil.WriteFile(filepath.Join(destDir,newFileName), []byte(strconv.Itoa(lineCount)), 0777)
        wg.Done()
    }

    for _, file := range files {
        wg.Add(1)
        go countFileLines(file.Name(), sourceDir, destDir, &wg)
    }

    wg.Wait()
    fmt.Println("Total number of processed files:" + strconv.Itoa(len(files)))
}
