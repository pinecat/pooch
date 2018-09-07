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
    "html/template" // for setting up html files
    "github.com/gorilla/mux" // for handling webapp functions
    "github.com/gorilla/securecookie" // for handling session info and security
    "pooch/mgopooch"
)

/* globals */
var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
var router = mux.NewRouter()

/*
    index_handler:  handles the index page of the web app

    params:         w - webpage writer
                    r - request getter

    returns:        void
*/
func index_handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("html/index.html")
    t.Execute(w, nil)
}

/*
    admin_handler:  handles the admin page of the webapp

    params:         w - webpage writer
                    r - request getter

    returns:        void
*/
func admin_handler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type == "admin" {
        t, _ := template.ParseFiles("html/admin.html")
        t.Execute(w, u)
    } else {
        http.Redirect(w, r, "/", 302)
    }
}

/*
    login_handler:  handles login requests

    params:         w - webpage writer
                    r - request getter

    returns:        void
*/
func login_handler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    redirect_target := "/"
    if username != "" && password != "" {
        // .. check credentials ..

        u, err := mgopooch.GetUser(username)
        if err != nil {
            http.Redirect(w, r, redirect_target, 302)
        }

        if mgopooch.IsPassValid(password, &u) == true {
            set_session(username, w)
            redirect_target = "/admin"
        }
    }
    http.Redirect(w, r, redirect_target, 302)
}

/*
    logout_handler:  handles logout requests

    params:         w - webpage writer
                    r - request getter

    returns:        void
*/
func logout_handler(w http.ResponseWriter, r *http.Request) {
    clear_session(w)
    http.Redirect(w, r, "/", 302)
}

/*
    set_session:    sets the user session for the webapp when logging in

    params:         username - username of person logging in
                    w - webpage writer

    returns:        void
*/
func set_session(username string, w http.ResponseWriter) {
    value := map[string]string{
        "name": username,
    }
    if encoded, err := cookieHandler.Encode("session", value); err == nil {
        cookie := &http.Cookie{
            Name:  "session",
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(w, cookie)
    }
}

/*
    get_username:   gets username from cookie/session of user logging in

    params:         r - request getter

    returns:        username - username of whoever is in the current session
*/
func get_username(r *http.Request) (username string) {
    if cookie, err := r.Cookie("session"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
            username = cookieValue["name"]
        }
    }
    return username
}

/*
    clear_session:  used to log users out of a session

    params:         w - webpage writer

    returns:        void
*/
func clear_session(w http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(w, cookie)
}

/*
    main:       main function of the program

    params:     na

    returns:    void
*/
func main() {
    mgopooch.SetupSession()

    router.HandleFunc("/", index_handler) // handle the index page
    router.HandleFunc("/admin", admin_handler) // handle the admin page

    router.HandleFunc("/login", login_handler).Methods("POST")
    router.HandleFunc("/logout", logout_handler).Methods("POST")

    http.Handle("/", router)
    fmt.Println("pooch web server running on port 8080....")
    http.ListenAndServe(":8080", nil) // start the webapp on port 8080
}
