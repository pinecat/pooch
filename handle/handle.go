/* package */
package handle

/* imports */
import (
    _"fmt"
    "time"
    "strings"
    "strconv"
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
const alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

/* structs */
type PageData struct {
    IP          string
    Port        string
    UserData    mgopooch.User
    BdngData    []mgopooch.Building
    Users       []mgopooch.User
    Msg         string
    Reloc       string
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
        disperr(w, "Invalid or missing information for user creation!", "/admin")
    } else {
        currentUsers, _ := mgopooch.GetAllUsers();
        for _, val := range currentUsers {
            if val.Username == username {
                disperr(w, "User already exists!", "/admin")
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
        disperr(w, "Please enter the username of the user you would like to remove!", "/admin")
        return
    } else if username == u.Username {
        disperr(w, "You may not remove your own account!", "/admin")
        return
    } else if username == "root" {
        disperr(w, "You may not remove the root account!", "/admin")
        return
    } else {
        err := mgopooch.RemoveUser(username)
        if err != nil {
            disperr(w, "Could not remove user!  Are you sure the user exists?", "/admin")
            return
        } else {
            http.Redirect(w, r, "/admin", 302)
        }
    }
}

func AdminChangeuserpasswordHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    username := r.FormValue("chusername")
    password := r.FormValue("chpassword")

    if username == "" {
        disperr(w, "Please enter the username of the user you would like to update the password for!", "/admin")
        return
    } else if password == "" {
        disperr(w, "Please enter a password to update to!", "/admin")
        return
    } else {
        err := mgopooch.UpdatePassword(username, password)
        if err != nil {
            disperr(w, "Could not update user password!  Are you sure the user exists?", "/admin")
            return
        } else {
            http.Redirect(w, r, "/admin", 302)
        }
    }
}

func AdminAddbuildingHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    name := r.FormValue("bdngname")
    abrv := r.FormValue("bdngabrv")

    if name == "" {
        disperr(w, "Please enter the name of the building you'd like to add!", "/admin")
        return
    } else if abrv == "" {
        disperr(w, "Please enter the abbreviation of the building name you'd like to add!", "/admin")
        return
    }

    b, _ := mgopooch.GetRooms()
    for _, data := range b {
        if data.Name == name || data.Abrv == abrv {
            disperr(w, "Building already exists!", "/admin")
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
        disperr(w, "Could not remove building!  Are you sure the building exists?", "/admin")
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

    if name == "" {
        disperr(w, "Please select a building to add a room to!", "/admin")
        return
    } else if num == "" {
        disperr(w, "Please enter the room number of the room you are trying to add!", "/admin")
        return
    } else if proj == "" {
        disperr(w, "Please select a projector type!", "/admin")
        return
    }

    if group == "" {
        group = "0"
    }

    b, _ := mgopooch.GetRooms()

    for _, data := range b {
        if data.Name == name {
            for roomnum, _ := range data.Rooms {
                if roomnum == num {
                    disperr(w, "Room already exists!", "/admin")
                    return
                }
            }
        }
    }
    mgopooch.InsertRoom(name, num, proj, group, "never checked")
    http.Redirect(w, r, "/admin", 302)
}

func AdminRmroomHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)

    name := r.FormValue("bdngname")
    num := r.FormValue("roomnum")

    err := mgopooch.RemoveRoom(name, num)
    if err != nil {
        disperr(w, "Could not remove room!  Are you sure the room exists?", "/admin")
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

func AdminRoomsCSVHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)
    mgopooch.Export2CSV()
    http.Redirect(w, r, "/public/roomdata.csv", 302)
    http.Redirect(w, r, "/admin/rooms", 302)
}

func AdminResetindroomHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)
    resetinfo := r.FormValue("hiddenreset")
    s := strings.Split(resetinfo, ".")
    mgopooch.ResetRoom(s[0], s[1])
    http.Redirect(w, r, "/admin/rooms", 302)
}

func AdminResetroomsHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)
    mgopooch.ResetAllRooms()
    http.Redirect(w, r, "/admin/rooms", 302)
}

func AdminResetusergroupsHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)
    mgopooch.ResetAllUserGroups()
    http.Redirect(w, r, "/admin/groups", 302)
}

func AdminResetroomgroupsHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)
    mgopooch.ResetAllRoomGroups()
    http.Redirect(w, r, "/admin/groups", 302)
}

func AdminTotaskHandler(w http.ResponseWriter, r *http.Request) {
    chadmin(w, r)
    http.Redirect(w, r, "/task", 302)
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
            if u.Group == info.Group && info.Status != "checked" {
                assign[i].Rooms[num] = info
            }
        }
    }

    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, UserData: u, BdngData: assign}

    t, _ := template.ParseFiles("html/task.html")
    t.Execute(w, pd)
}

func TaskRoomHandler(w http.ResponseWriter, r *http.Request) {
    chlogin(w, r)

    roomselect := r.FormValue("roomselect")
    problist := r.FormValue("probcatlist")
    lamph := r.FormValue("lamphours")
    notes := r.FormValue("notes")
    if roomselect == "" {
        disperr(w, "Please select a room to submit!", "/task")
        return
    }

    var abrv string = ""
    var roomnum string = ""
    for _, c := range roomselect {
        if strings.Contains(alpha, strings.ToUpper(string(c))) {
            abrv += string(c)
        } else {
            roomnum += string(c)
        }
    }

    b, _ := mgopooch.GetRooms()
    var bdngName string
    for _, data := range b {
        if abrv == data.Abrv {
            bdngName = data.Name
        }
    }

    room, _ := mgopooch.GetRoom(bdngName, roomnum)
    if lamph == "" {
        disperr(w, "Please enter lamphours for the room!", "/task")
        return
    }
    if room.Lamph.Interactive != -1 {
        s := strings.Split(lamph, ",")
        if len(s) == 1 {
            disperr(w, "Please enter lamphours for the interactive projector as well!", "/task")
            return
        } else {
            room.Lamph.Standard, _ = strconv.Atoi(s[0])
            room.Lamph.Interactive, _ = strconv.Atoi(s[0])
        }
    } else {
        s := strings.Split(lamph, ",")
        room.Lamph.Standard, _ = strconv.Atoi(s[0])
    }

    var problems []int
    if problist != "" {
        buffer := trim_space(problist)
        s := strings.Split(buffer, ",")
        for i := 0; i < len(s); i++ {
            p, _ := strconv.Atoi(s[i])
            problems = append(problems, p)
        }
        room.Probcat = problems
    }

    room.Status = "checked"
    room.Notes = notes;
    room.Last = time.Now().Format("2006-01-02 15:04:05")

    mgopooch.UpdateRoomStatus(bdngName, roomnum, &room)

    http.Redirect(w, r, "/task", 302)
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

func trim_space(buffer string) string {
    var str string = ""
    for _, c := range buffer {
        if c != ' ' {
            str += string(c)
        }
    }
    return str
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

func disperr(w http.ResponseWriter, msg string, reloc string) {
    t, _ := template.ParseFiles("html/error.html")
    pd := PageData{IP: confrdr.PoochConf.IP, Port: confrdr.PoochConf.Port, Msg: msg, Reloc: reloc}
    t.Execute(w, pd)
}
