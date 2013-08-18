package main

import (
	"github.com/dgryski/go-identicon"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {

	args := strings.Split(r.URL.Path, "/")
	args = args[1:]

	if len(args) != 1 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := args[0]

	// support jpg too?
	if !strings.HasSuffix(item, ".png") {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item = strings.TrimSuffix(item, ".png")

	key := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	icon := identicon.New7x7(key)

	log.Printf("creating identicon for '%s'\n", item)

	data := []byte(item)
	pngdata := icon.Render(data)

	w.Header().Set("Content-Type", "image/png")
	w.Write(pngdata)

	return
}

func main() {
	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	log.Println("Listening on port", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))
}
