package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
	"path/filepath"
	"sync"
)

func main() {
    downloadImages(52, "images", 25)
    fmt.Println("Downloaded 52 random Unsplash images to the 'images' folder.")
}

func downloadImages(count int, folderPath string, workers int) {
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
                imageUrl := fmt.Sprintf("https://source.unsplash.com/random/300x300?sig=%d", job)
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
