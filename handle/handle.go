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
    Users       []mgopooch.User
    Msg         string
}

func TempHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("html/temp.html")
    t.Execute(w, nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port}

    t, _ := template.ParseFiles("html/index.html")
    t.Execute(w, pd)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
    u := chadmin(w, r)

    b, _ := mgopooch.GetRooms()
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u, BdngData: b}

    t, _ := template.ParseFiles("html/admin.html")
    t.Execute(w, pd)
}

func AdminRoomsHandler(w http.ResponseWriter, r *http.Request) {
    u := chadmin(w, r)

    b, _ := mgopooch.GetRooms()
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u, BdngData: b}

    t, _ := template.ParseFiles("html/rooms.html")
    t.Execute(w, pd)
}

func AdminGroupsHandler(w http.ResponseWriter, r *http.Request) {
    u := chadmin(w, r)

    b, _ := mgopooch.GetRooms()
    ul, _ := mgopooch.GetAllUsers()
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u, BdngData: b, Users: ul}

    t, _ := template.ParseFiles("html/groups.html")
    t.Execute(w, pd)
}


func AdminCreateuserHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    username := r.FormValue("mkusername")
    password := r.FormValue("password")
    fname := r.FormValue("fname")
    lname := r.FormValue("lname")
    acc := r.FormValue("acctype")

    if (username == "" || password == "" || fname == "" || lname == "" || acc == "") {
        disperr(w, "Invalid or missing information for user creation!")
    } else {
        currentUsers, _ := mgopooch.GetAllUsers();
        for _, val := range currentUsers {
            if val.Username == username {
                disperr(w, "User already exists!")
                return
            }
        }
        u := mgopooch.User{Username: username, Password: password, Fname: fname, Lname: lname, Type: acc}
        mgopooch.InsertUser(&u)
        http.Redirect(w, r, "/admin", 302)
    }
}

func AdminRemoveuserHandler(w http.ResponseWriter, r *http.Request) {
    u := chadmin(w, r)

    username := r.FormValue("rmusername")

    if username == "" {
        disperr(w, "Please enter the username of the user you would like to remove!")
    } else if username == u.Username {
        disperr(w, "You may not remove your own account!")
    } else {
        err := mgopooch.RemoveUser(username)
        if err != nil {
            disperr(w, "Could not remove user!  Are you sure the user exists?")
        } else {
            http.Redirect(w, r, "/admin", 302)
        }
    }
}

func AdminAddbuildingHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    name := r.FormValue("bdngname")
    abrv := r.FormValue("bdngabrv")

    b, _ := mgopooch.GetRooms()
    for _, data := range b {
        if data.Name == name || data.Abrv == abrv {
            disperr(w, "Building already exists!")
            return
        }
    }
    mgopooch.InsertBuilding(name, abrv)
    http.Redirect(w, r, "/admin", 302)
}

func AdminRmbuildingHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    name := r.FormValue("bdngname")
    err := mgopooch.RemoveBuilding(name)
    if err != nil {
        disperr(w, "Could not remove building!  Are you sure the building exists?")
    } else {
        http.Redirect(w, r, "/admin", 302)
    }
}

func AdminAddroomHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    name := r.FormValue("bdngname")
    num := r.FormValue("roomnum")
    proj := r.FormValue("projtype")
    group := r.FormValue("group")

    if group == "" {
        group = "0"
    }

    b, _ := mgopooch.GetRooms()

    for _, data := range b {
        if data.Name == name {
            for roomnum, _ := range data.Rooms {
                if roomnum == num {
                    disperr(w, "Room already exists!")
                    return
                }
            }
        }
    }
    mgopooch.InsertRoom(name, num, proj, group)
    http.Redirect(w, r, "/admin", 302)
}

func AdminRmroomHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    name := r.FormValue("bdngname")
    num := r.FormValue("roomnum")

    err := mgopooch.RemoveRoom(name, num)
    if err != nil {
        disperr(w, "Could not remove room!  Are you sure the room exists?")
    } else {
        http.Redirect(w, r, "/admin", 302)
    }
}

func AdminSaveusergroupsHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    ul, _ := mgopooch.GetAllUsers()
    for _, data := range ul {
        newGroup := r.FormValue(data.Username + "group")
        if newGroup != "" {
            mgopooch.UpdateUserGroup(data.Username, newGroup)
        }
    }
    http.Redirect(w, r, "/admin/groups", 302)
}

func AdminSaveroomgroupsHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    b, _ := mgopooch.GetRooms()
    for _, data := range b {
        for num, _ := range data.Rooms {
            newGroup := r.FormValue(data.Abrv + num + "group")
            if newGroup != "" {
                mgopooch.UpdateBuildingGroup(data.Name, num, newGroup)
            }
        }
    }

    http.Redirect(w, r, "/admin/groups", 302)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
    u := chlogin(w, r)
    b, _ := mgopooch.GetRooms()

    assign := make([]mgopooch.Building, len(b))
    for i := 0; i < len(b); i++ {
        assign[i].Name = b[i].Name
        assign[i].Abrv = b[i].Abrv
        assign[i].Rooms = make(map[string]mgopooch.Room)
        for num, info := range b[i].Rooms {
            if u.Group == info.Group {
                assign[i].Rooms[num] = info
            }
        }
    }

    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u, BdngData: assign}

    t, _ := template.ParseFiles("html/task.html")
    t.Execute(w, pd)
}

func TaskRoomHandler(w http.ResponseWriter, r *http.Request) {
    //u := chlogin(w, r)


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
            if u.Type == "regular" {
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

func chadmin(w http.ResponseWriter, r *http.Request) mgopooch.User {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return mgopooch.User{}
    }
    return u
}

func chlogin(w http.ResponseWriter, r *http.Request) mgopooch.User {
    u, _ := mgopooch.GetUser(get_username(r))
    if u.Type != "regular" && u.Type != "admin" {
        http.Redirect(w, r, "/", 302)
        return mgopooch.User{}
    }
    return u
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

func disperr(w http.ResponseWriter, msg string) {
    t, _ := template.ParseFiles("html/error.html")
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port}
    pd.Msg = msg
    t.Execute(w, pd)
}
