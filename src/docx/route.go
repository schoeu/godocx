package main

import (
	"net/http"
    "regexp"
    "fmt"
)

type regRoute struct {
    pattern string
    handler func (w http.ResponseWriter, r *http.Request)
}

var routes = make([]regRoute, 5)




func initial() {
    setRegRoute("*.md", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "*.md")
    })
}


func allRoutes(w http.ResponseWriter, r *http.Request) {

    // 添加路由
    initial()

    url := r.URL

    for _, v := range routes {
        reg, err := regexp.Compile(v.pattern)
        if err != nil {
            continue
        }
        if reg.MatchString(url.Path) {
            v.handler(w, r)
        }
    }

    fmt.Fprint(w, "404 Page Not Found!")

}

func setRegRoute(p string, h func(http.ResponseWriter, *http.Request)) {
    r := regRoute{p, h}
    routes = append(routes, r)
}
