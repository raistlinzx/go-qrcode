package main

import "github.com/go-martini/martini"
import "net/http"

func main() {
  m := martini.Classic()
  m.Get("/", func(res http.ResponseWriter, req *http.Request) {
    pngBytes := Generate(req.FormValue("url"), req.FormValue("fcolor"), req.FormValue("bcolor"), req.FormValue("size"), req.FormValue("logo"))
    res.Write(pngBytes)
  })
  m.Run()
}

