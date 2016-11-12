package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"html/template"
	"encoding/json"
	"walker"
)

var (
	index    = "readme.md"
	path     = "/Users/memee/Downloads/svn/ps-fe"
	docxConf = "./docx-conf.json"
)

func main() {
	http.HandleFunc("/", routeHandle)

	walker.ReadDirRs()

	err := http.ListenAndServe(":8910", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type person struct {
	name string
	age  int
}

var e interface{}

func routeHandle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/re", http.StatusFound)

	content, _ := ioutil.ReadFile(docxConf)

	json.Unmarshal(content, &e)

	dinfo, ok := e.(map[string]interface{})
	if ok {
		for k, v := range dinfo {
			fmt.Println(k, "%s--------\n")
			fmt.Println(v)
		}

	}
	//fmt.Fprintf(w, "re")

	// mdStr := []byte(`# title标题`)
	//output := blackfriday.MarkdownBasic(mdStr)

	//fmt.Println(string(output))

	t := person{
		name: "name",
		age:  20,
	}

	fmt.Println(t)
	fmt.Fprintf(w, "re")

}
