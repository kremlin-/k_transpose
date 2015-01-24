package main

import "fmt"
import "net/http"

func main() {

    /* make sure to close() anything you need to (you need to) */
    fmt.Println("welp");
    resp, err := http.Get("http://kremlin.cc")

    if err != nil {}
    fmt.Println(resp)
}
