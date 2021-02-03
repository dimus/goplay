package main

import (
	"fmt"
	"io"
	"log"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	shellUrl := "localhost:5001"
	examplesHash := "QmR5bphF4SkG7SDcw8H1kXTUu4jyR1HAtEPRHybDnWM6U6"
	s := shell.NewShell(shellUrl)

	list, err := s.List(fmt.Sprintf("/ipfs/%s", examplesHash))
	if err != nil {
		log.Fatal(err)
	}
	for i := range list {
		fmt.Printf("list: %+v", *list[i])
		fileName := list[i].Name
		hash := list[i].Hash
		log.Printf("downloading %s", fileName)
		r, err := s.Cat("/ipfs/" + hash)
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		buffer := make([]byte, 1000)
		for {
			bytes, err := r.Read(buffer)
			_, err2 := f.Write(buffer[:bytes])
			if err2 != nil {
				log.Fatal(err)
			}
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				break
			}
		}
	}

}
