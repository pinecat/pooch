package main

import (
    "honnef.co/go/js/dom"
)

func main() {
    println("Hi.")

    d := dom.GetWindow().Document()
    //clickme := d.GetElementByID("clickme").(*dom.HTMLButtonElement)
    display := d.GetElementByID("content").(*dom.HTMLDivElement)
    display.SetInnerHTML("Hello World")
    // clickme.AddEventListener("click", false, func(event dom.Event) {
    //     display.SetInnerHTML("Hello World")
    // })
}
