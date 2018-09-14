/* package */
package handle

/* imports */
import (
    _"fmt"
    "net/http"
    "html/template" // for setting up html files
    "github.com/gorilla/securecookie" // for handling session info and security
    "github.com/pinecat/pooch/mgopooch"
    "github.com/gorilla/websocket"
    "github.com/pinecat/pooch/confrdr"
)

/* gloabals */
var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
var upgrader = websocket.Upgrader {
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

/* structs */
type PageData struct {
    IP          string
    Port        string
    UserData    mgopooch.User
    BdngData    []mgopooch.Building
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port}

    t, _ := template.ParseFiles("html/index.html")
    t.Execute(w, pd)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u}

    t, _ := template.ParseFiles("html/admin.html")
    t.Execute(w, pd)
}

func AdminRoomsHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    b, _ := mgopooch.GetRooms()
    pd := PageData{confrdr.PoochConf.IP, confrdr.PoochConf.Port, u, b}

    t, _ := template.ParseFiles("html/rooms.html")
    t.Execute(w, pd)
}


func AdminGroupsHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u}

    t, _ := template.ParseFiles("html/groups.html")
    t.Execute(w, pd)
}


func AdminCreateuserHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    username := r.FormValue("mkusername")
    password := r.FormValue("password")
    fname := r.FormValue("fname")
    lname := r.FormValue("lname")
    acc := r.FormValue("acctype")

    u = mgopooch.User{Username:username, Password:password, Fname:fname, Lname:lname, Type:acc}
    mgopooch.InsertUser(&u)

    http.Redirect(w, r, "/admin", 302)
}

func AdminRemoveuserHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    username := r.FormValue("rmusername")
    if username == u.Username {
        http.Redirect(w, r, "/admin", 302)
        return
    }
    mgopooch.RemoveUser(username)

    http.Redirect(w, r, "/admin", 302)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "reg" && u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u}

    t, _ := template.ParseFiles("html/task.html")
    t.Execute(w, pd)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    redirect_target := "/"
    if username != "" && password != "" { // if the username and password are not blank
        // check credentials....

        // get user and check if they exist (err is nil if user exists)
        u, err := mgopooch.GetUser(username)
        if err != nil {
            http.Redirect(w, r, redirect_target, 302)
        }

        // make sure the password matches the hash in database
        if mgopooch.IsPassValid(password, &u) == true {
            set_session(username, w)
            if u.Type == "reg" {
                redirect_target = "/task"
            } else if u.Type == "admin" {
                redirect_target = "/admin"
            }
        }
    }
    http.Redirect(w, r, redirect_target, 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    clear_session(w)
    http.Redirect(w, r, "/", 302)
}

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

func get_username(r *http.Request) (username string) {
    if cookie, err := r.Cookie("session"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
            username = cookieValue["name"]
        }
    }
    return username
}

func clear_session(w http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(w, cookie)
}
