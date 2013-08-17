package main

import (
	"github.com/dgryski/identicon"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("usage: demo \"text to hash\"")
	}

	key := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	icon := identicon.New7x7(key)

        log.Printf("creating identicon for '%s'\n", os.Args[1])

	data := []byte(os.Args[1])
        pngdata := icon.Render(data)
        
        log.Println("writing output.png")
        f, err := os.Create("output.png")
        defer f.Close()
        if err != nil {
            log.Fatalf("cannot create output.png: %s\n", err)

        }

        f.Write(pngdata)
}
