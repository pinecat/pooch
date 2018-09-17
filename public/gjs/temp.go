package main

import (
    "net"
    "github.com/gopherjs/gopherjs/js"
)

func main() {
    //js.Global.Get("document").Call("write", "Hello World!")
    conn, err := net.Dial("tcp", "192.168.200.132:8081")
    if err != nil {
        println(err)
    }
    conn.Write([]byte("Hello Server!"))

    doc := js.Global.Get("document")
    resDiv := doc.Call("getElementById", "mydiv")
	resDiv.Set("innerHTML", "Hello World!")
}
