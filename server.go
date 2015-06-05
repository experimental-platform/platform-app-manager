package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	// for _, img := range imgs {
	// 	fmt.Println("ID: ", img.ID)
	// 	fmt.Println("RepoTags: ", img.RepoTags)
	// 	fmt.Println("Created: ", img.Created)
	// 	fmt.Println("Size: ", img.Size)
	// 	fmt.Println("VirtualSize: ", img.VirtualSize)
	// 	fmt.Println("ParentId: ", img.ParentID)
	// }
	fmt.Printf("%+v\n", client.List())
}

var client *Dokku
var CONTAINER_NAME = "dokku"

func main() {
	var port int
	flag.IntVar(&port, "port", 3001, "server port")
	flag.Parse()
	fmt.Println("Port: ", port)

	var err error
	client, err = NewDokku()

	if err != nil {
		panic("Container '" + CONTAINER_NAME + "' is not running!")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
