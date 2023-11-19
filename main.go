package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"
)

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func stringify(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}

	return string(b)
 }

func main() {
	for _, s := range find("/home/kenny/Downloads", ".epub") {

		s, err := epub.GetMetadataFromFile(s)
		 if err != nil {
		 	
      errors.Cause(err)
      fmt.Printf("%+s:%d\n",err,err)
		 	   os.Exit(1)
		 }

		fmt.Printf("%s\n", (s.Title))
	}
  }
