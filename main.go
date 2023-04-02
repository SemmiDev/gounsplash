package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
	"path/filepath"
	"sync"
    "strconv"
)

func main() {
    var total int
    var folderName, searchTerm string

    if len(os.Args) > 1 {
        totalInStr := os.Args[1]
        totalInInt, err := strconv.Atoi(totalInStr)
        if err != nil {
            fmt.Println("Please enter a valid number for the total of images to download.")
            os.Exit(1)
        }

        total = totalInInt
        folderName = os.Args[2]
        searchTerm = os.Args[3]
    } else {
        total = 52
        folderName = "images"
        searchTerm = "drink"
    }

    downloadImages(total, folderName, searchTerm, 10)
    fmt.Printf("Downloaded %d images to %s folder\n", total, folderName)
}

func downloadImages(count int, folderPath string, searchTerm string, workers int) {
    folderLocation := createFolderIfNotExists(folderPath)
    wg := &sync.WaitGroup{}
    jobs := make(chan int, count)
    for i := 1; i <= count; i++ {
        jobs <- i
    }
    close(jobs)

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                imageUrl := fmt.Sprintf("https://source.unsplash.com/random/300x300/?%s", searchTerm)
                fileName := fmt.Sprintf("%02d.jpg", job)
                fileLocation := filepath.Join(folderLocation, fileName)
                err := downloadImage(imageUrl, fileLocation)
                if err != nil {
                    fmt.Printf("Failed to download image %s: %v\n", imageUrl, err)
                }
                fmt.Printf("Downloaded image %d of %d\n", job, count)
            }
        }()
    }
    wg.Wait()
}

func createFolderIfNotExists(folderPath string) string {
    _, err := os.Stat(folderPath)
    if os.IsNotExist(err) {
        os.MkdirAll(folderPath, os.ModePerm)
    }
    return folderPath
}

func downloadImage(imageUrl string, filePath string) error {
    response, err := http.Get(imageUrl)
    if err != nil {
        return err
    }
    defer response.Body.Close()
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = io.Copy(file, response.Body)
    if err != nil {
        return err
    }
    return nil
}
