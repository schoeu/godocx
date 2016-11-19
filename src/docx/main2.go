package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
    "regexp"
)

var (
	index = "reademe.md"
    port = "8910"
    mdReg = ".+.md$"
)

func main() {
    router := gin.Default()

    // This handler will match /user/john but will not match neither /user/ or /user

    router.GET("/", func(c *gin.Context) {
        urlPath := c.Request.URL.Path
        ism, _ := regexp.MatchString(mdReg, urlPath)
        c.JSON(200, gin.H{
            "path": urlPath,
            "message": ism,
        })
        // _, err := regexp.MatchString(mdReg, urlPath)
        // if err != nil {
        //     c.String(http.StatusOK, "yes, it is md")
        // }
    })

    router.GET(index, func(c *gin.Context) {
        c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"mdHtml": "Main website",
		})
    })

    //util.GetRsHTML()

    router.Static("/static", "../../themes/default/static")
    router.StaticFile("/favicon.ico", "../../themes/default/static/favicon.ico")
	

    router.Run(":" + port)
}