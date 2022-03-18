package main

import (
    "github.com/yonikosiner/go-wiki/server"
    "github.com/yonikosiner/go-wiki/utils"
)

func main() {
    utils.SearchWiki("")
    server.Run()
}
