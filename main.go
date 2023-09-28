package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Get the current directory.
    dir, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create a channel to store the .epub and .pdf files.
    epubsAndPdfs := make(chan string)

    // Go through all the files in the current directory.
    for _, file := range filepath.Glob("*") {
        // Check if the file is an .epub or .pdf file.
        if filepath.Ext(file) == ".epub" || filepath.Ext(file) == ".pdf" {
            // Send the file to the channel.
            epubsAndPdfs <- file
        }
    }

    // Close the channel.
    close(epubsAndPdfs)

    // Iterate over the channel and print the .epub and .pdf files.
    for file := range epubsAndPdfs {
        fmt.Println(file)
    }
}

