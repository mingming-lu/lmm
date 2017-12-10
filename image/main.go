package main

import (
    "os"

    "github.com/akinaru-lu/elesion"
)

const Res = "image/res"
const Special = "image/special"

func handler(c *elesion.Context) {
    c.File(Res + c.Request.URL.Path)
}

func avatar(c *elesion.Context) {
	c.File(Special + "/avatar.jpg")
}

func ensureDir(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.Mkdir(path, os.ModePerm)
    }
}

func init() {
	ensureDir(Res)
	ensureDir(Special)
}

func main() {
    el := elesion.Default("[image]")
    el.Handle("/", handler)
    el.Handle("/avatar", avatar)
    el.Run(":8082")
}
