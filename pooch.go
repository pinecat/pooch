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
    "log" // for logging errors
    "os" // for getting cmdline args
    "net/http" // for hosting webapp server
    "github.com/gorilla/mux" // gorilla mux router for handling funcs
    "github.com/pinecat/pooch/mgopooch"
    "github.com/pinecat/pooch/handle"
    "github.com/pinecat/pooch/confrdr"
)

/* globals */
var router = mux.NewRouter()

func help() {
    fmt.Printf("POOCH\n")
    fmt.Printf("\tThe friendly podium check helper.  Web app to help with podium checks at Elizabethtown College.\n\n")
    fmt.Printf("USAGE\n")
    fmt.Printf("\tUsage: pooch [-f filepath | -h]\n\n")
    fmt.Printf("OPTIONS\n")
    fmt.Printf("\t-f: specify the filepath of the config file (\"./pooch.conf\" by default)\n")
    fmt.Printf("\t-h: this menu\n\n")
    fmt.Printf("AUTHOR\n")
    fmt.Printf("\tRory Dudley (aka pinecat: https://github.com/pinecat/pooch)\n\n")
    fmt.Printf("EOF\n")
}

/*
    main:       main function of the program

    params:     na

    returns:    void
*/
func main() {
    filepath := "pooch.conf"
    if len(os.Args) == 3 && os.Args[1] == "-f" { // read in different filepath if specified by user at cmdline
        filepath = os.Args[2] // update the filepath
    } else if len(os.Args) > 1 && os.Args[1] == "-f" { // if format for -f flag is not correct...
        fmt.Printf("Usage: %s [-f filepath | -h]\n", os.Args[0]) // print a usage message
        return // and exit the program
    } else if len(os.Args) > 1 && os.Args[1] == "-h" { // if flag is -h...
        help() // ...print a help menu
        return // and exit the program
    }

    confrdr.ReadConfFile(filepath)
    mgopooch.SetupSession()

    router.HandleFunc("/", handle.IndexHandler) // handle the index page
    router.HandleFunc("/admin", handle.AdminHandler) // handle the admin page
    router.HandleFunc("/admin/createuser", handle.AdminCreateuserHandler)
    router.HandleFunc("/admin/removeuser", handle.AdminRemoveuserHandler)
    router.HandleFunc("/admin/chpswd", handle.AdminChangeuserpasswordHandler)

    router.HandleFunc("/admin/rooms", handle.AdminRoomsHandler)
    router.HandleFunc("/admin/groups", handle.AdminGroupsHandler)

    router.HandleFunc("/admin/addbdng", handle.AdminAddbuildingHandler)
    router.HandleFunc("/admin/rmbdng", handle.AdminRmbuildingHandler)
    router.HandleFunc("/admin/addroom", handle.AdminAddroomHandler)
    router.HandleFunc("/admin/rmroom", handle.AdminRmroomHandler)
    router.HandleFunc("/admin/saveusergroups", handle.AdminSaveusergroupsHandler)
    router.HandleFunc("/admin/saveroomgroups", handle.AdminSaveroomgroupsHandler)
    router.HandleFunc("/admin/exportcsv", handle.AdminRoomsCSVHandler)
    router.HandleFunc("/admin/resetrooms", handle.AdminResetroomsHandler)
    router.HandleFunc("/admin/resetindroom", handle.AdminResetindroomHandler)
    router.HandleFunc("/admin/resetusergroups", handle.AdminResetusergroupsHandler)
    router.HandleFunc("/admin/resetroomgroups", handle.AdminResetroomgroupsHandler)
    router.HandleFunc("/admin/totask", handle.AdminTotaskHandler)

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
