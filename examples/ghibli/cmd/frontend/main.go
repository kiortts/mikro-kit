/*
Простой http сервер для директории "web"
*/

package main

import (
	"log"
	"net/http"

	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Println("Server started on port: 8090")
	log.Fatal(http.ListenAndServe(":8090", nil))

	utils.OpenInBrowser("localhost:")
}
