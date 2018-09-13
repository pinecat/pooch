/************************************************************
*   title:  pooch
*   desc:   the friendly podium check helper
*   files:  pooch.go
*   author: rory dudley
************************************************************/

/* pacakge */
package main

/* imports */
import (
    "fmt" // for printing out text
    _"log" // for logging errors
    "net/http" // for hosting webapp server
    "github.com/gorilla/mux" // gorilla mux router for handling funcs
    "pooch/mgopooch"
    "pooch/handle"
    "labix.org/v2/mgo/bson"
)

/* globals */
var router = mux.NewRouter()

/*
    main:       main function of the program

    params:     na

    returns:    void
*/
func main() {
    mgopooch.SetupSession()

    router.HandleFunc("/", handle.IndexHandler) // handle the index page
    router.HandleFunc("/admin", handle.AdminHandler) // handle the admin page
    router.HandleFunc("/admin/createuser", handle.AdminCreateuserHandler)
    router.HandleFunc("/admin/removeuser", handle.AdminRemoveuserHandler)

    router.HandleFunc("/admin/rooms", handle.AdminRoomsHandler)
    router.HandleFunc("/admin/groups", handle.AdminGroupsHandler)

    router.HandleFunc("/login", handle.LoginHandler).Methods("POST")
    router.HandleFunc("/logout", handle.LogoutHandler).Methods("POST")

    http.Handle("/", router)
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

    fmt.Println("Pooch web server running on port 8080....")
    http.ListenAndServe(":8080", nil) // start the webapp on port 8080
}
