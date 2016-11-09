package main

import (
	"fmt"
	"log"
	"net/http"

	"./walker"
	//"html/template"
	//"github.com/russross/blackfriday"
)

var (
	index = "readme.md"
	path  = "/Users/memee/Downloads/svn/ps-fe"
)

func main() {
	http.HandleFunc("/", routeHandle)
	http.HandleFunc("/re", reRouteHandle)
	err := http.ListenAndServe(":8910", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func routeHandle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/re", http.StatusFound)
}

func reRouteHandle(w http.ResponseWriter, r *http.Request) {
	walker.ReadDirRs()
	fmt.Fprintf(w, "re")
}
