package main

// import (
//   "fmt"
//   "os"
//   "path/filepath"
// )
import (
  "fmt"
  "os"
  "path/filepath"
)
func main() {
// Scan the current directory (can also be set to a certain directory)   
// filepath.Walk("/home/kenny/Downloads/", func(path string, info fs.FileInfo, err error) error {
//     fmt.Println(path),
//
//     return nil
//   }
//
// package main
//


  filepath.Walk("/home/kenny/Downloads", func(path string, info os.FileInfo, err error) error {
    fmt.Println(path)
    return nil
  })
  }
