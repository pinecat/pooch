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
    _"fmt" // for printing out text
    "log" // for logging errors
    "net/http" // for hosting webapp server
    "github.com/gorilla/mux" // gorilla mux router for handling funcs
    "github.com/pinecat/pooch/mgopooch"
    "github.com/pinecat/pooch/handle"
    "github.com/pinecat/pooch/confrdr"
)

/* globals */
var router = mux.NewRouter()

/*
    main:       main function of the program

    params:     na

    returns:    void
*/
func main() {
    confrdr.ReadConfFile("pooch.conf")
    mgopooch.SetupSession()

    router.HandleFunc("/", handle.IndexHandler) // handle the index page
    router.HandleFunc("/admin", handle.AdminHandler) // handle the admin page
    router.HandleFunc("/admin/createuser", handle.AdminCreateuserHandler)
    router.HandleFunc("/admin/removeuser", handle.AdminRemoveuserHandler)

    router.HandleFunc("/admin/rooms", handle.AdminRoomsHandler)
    router.HandleFunc("/admin/groups", handle.AdminGroupsHandler)

    router.HandleFunc("/admin/addbdng", handle.AdminAddbuildingHandler)
    router.HandleFunc("/admin/rmbdng", handle.AdminRmbuildingHandler)
    router.HandleFunc("/admin/addroom", handle.AdminAddroomHandler)
    router.HandleFunc("/admin/rmroom", handle.AdminRmroomHandler)
    router.HandleFunc("/admin/saveusergroups", handle.AdminSaveusergroupsHandler)
    router.HandleFunc("/admin/saveroomgroups", handle.AdminSaveroomgroupsHandler)

    router.HandleFunc("/task", handle.TaskHandler)
    router.HandleFunc("/task/submitroom", handle.TaskRoomHandler)

    router.HandleFunc("/temp", handle.TempHandler)

    router.HandleFunc("/login", handle.LoginHandler).Methods("POST")
    router.HandleFunc("/logout", handle.LogoutHandler).Methods("POST")

    http.Handle("/", router)
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

    log.Println("[LOG] Pooch web server running on: http://" + confrdr.PoochConf.IP + ":" + confrdr.PoochConf.Port + "/")
    http.ListenAndServe(":" + confrdr.PoochConf.Port, nil) // start the webapp on port 8080
    // err := http.ListenAndServeTLS(":8080", "ssl/pooch.cert", "ssl/pooch.key", nil) // start the webapp on port 8080
}
