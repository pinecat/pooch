/* package */
package handle

/* imports */
import (
    "net/http"
    "html/template" // for setting up html files
    "github.com/gorilla/securecookie" // for handling session info and security
    "github.com/pinecat/pooch/mgopooch"
    "github.com/gorilla/websocket"
)

/* gloabals */
var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
var upgrader = websocket.Upgrader {
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("html/index.html")
    t.Execute(w, nil)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type == "admin" {
        t, _ := template.ParseFiles("html/admin.html")
        t.Execute(w, u)
    } else {
        http.Redirect(w, r, "/", 302)
    }
}

func AdminRoomsHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    t, _ := template.ParseFiles("html/rooms.html")
    t.Execute(w, u)
}


func AdminGroupsHandler(w http.ResponseWriter, r *http.Request) {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return
    }

    t, _ := template.ParseFiles("html/groups.html")
    t.Execute(w, u)
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
            redirect_target = "/admin"
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
