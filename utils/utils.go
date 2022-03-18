package utils

import (
    "fmt"
    "strings"
    "os/exec"
	"github.com/russross/blackfriday"
)

func SearchWiki(query string) string {
    cmd := exec.Command("find", "./wiki-files", "-name", "*" + query + "*.md")
    out, err := cmd.Output()

    if err != nil {
        fmt.Println(err)
        return "Sorry there was an error"
    }

    out = []byte(strings.Replace(string(out), "./wiki-files/", "", -1))
    out = []byte(strings.Replace(string(out), ".md", "", -1))

    //Split each item by newline
    var items []string
    for _, item := range strings.Split(string(out), "\n") {
        if item != "" {
            items = append(items, item)
        }
    }

    //Render each item to have a link
    var links []string
    for _, item := range items {
        links = append(links, fmt.Sprintf("[%s](/view/%s)", item, item))
    }

    //Convert markdown to html
    var html []byte

    for _, item := range links {
        html = append(html, []byte(item)...)
        html = append(html, []byte("<br>")...)
    }

    return string(blackfriday.MarkdownCommon(html))
}

