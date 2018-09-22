/* package */
package mgopooch

/* imports */
import (
    _"fmt"
    "log"
    "golang.org/x/crypto/bcrypt" // for hashing passwords
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "github.com/pinecat/pooch/confrdr"
    "strconv"
    "errors"
)

/* structs */
type User struct {
    ID          bson.ObjectId `bson:"_id,omitempty"`
    Username    string
    Password    string
    Fname       string
    Lname       string
    Type        string
    Group       int
}

type AllRooms struct {
    RoomMap map[string]Building
}

type Building struct {
    ID      bson.ObjectId `bson:"_id,omitempty"`
    Name    string
    Rooms   map[string]Room
    Abrv    string
}

type Room struct {
    Status  string
    Group   int
    Lamph   Projector
    Probcat []int
    Notes   string
}

type Projector struct {
    Standard    int
    Interactive int
}

/* globals */
var (
    Session *mgo.Session
    DB      *mgo.Database
)

func SetupSession() {
    Session, err := mgo.Dial(confrdr.DBConf.Url)
    if err != nil {
        log.Fatal("[ERR] Could not connect to database!!!")
        log.Fatal(err)
    }
    DB = Session.DB("pooch")
    DB.Login(confrdr.DBConf.User, confrdr.DBConf.Pass)
}

func GetUser(username string) (User, error) {
    c := DB.C("users")
    result := User{}
    err := c.Find(bson.M{"username":username}).One(&result)
    return result, err
}

func GetAllUsers() ([]User, error) {
    c := DB.C("users")
    result := []User{}
    err := c.Find(nil).All(&result)
    return result, err
}

func IsPassValid(password string, user *User) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return false
    }
    return true
}

func InsertUser(user *User) error {
    hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hash)
    err := DB.C("users").Insert(&user)
    return err
}

func RemoveUser(username string) error {
    err := DB.C("users").Remove(bson.M{"username":username})
    return err
}

func InsertBuilding(name string, abrv string) error {
    building := Building{Name: name, Abrv: abrv}
    err := DB.C("buildings").Insert(&building)
    return err
}

func RemoveBuilding(name string) error {
    err := DB.C("buildings").Remove(bson.M{"name":name})
    return err
}

func InsertRoom(name string, num string, proj string, group string) error {
    var r Room
    if proj == "stanplusinter" {
        r.Lamph.Interactive = 0;
    } else if proj == "standard" {
        r.Lamph.Interactive = -1;
    }
    r.Status = "unchecked"
    r.Group, _ = strconv.Atoi(group)
    r.Probcat = []int{0}
    r.Notes = ""
    err := DB.C("buildings").Update(bson.M{"name":name}, bson.M{"$set": bson.M{"rooms." + num + "": r}})
    return err
}

func RemoveRoom(name string, num string) error {
    r, err := GetRoom(name, num)
    if err != nil {
        return err
    }
    err = DB.C("buildings").Update(bson.M{"name":name}, bson.M{"$unset": bson.M{"rooms." + num + "": r}})
    return err
}

func GetRoom(name string, num string) (Room, error) {
    c := DB.C("buildings")
    var result Building
    err := c.Find(bson.M{"name":name}).One(&result)
    if err != nil {
        return Room{}, err
    }

    var r Room
    for i, data := range result.Rooms {
        if i == num {
            r = data
        }
    }
    if r.Status != "unchecked" && r.Status != "checked" {
        err = errors.New("Room not found!")
    }
    return r, err
}

func GetRooms() ([]Building, error) {
    c := DB.C("buildings")
    var result []Building
    err := c.Find(nil).All(&result)
    return result, err
}

func UpdateUserGroup(username string, group string) error {
    gr, _ := strconv.Atoi(group)
    err := DB.C("users").Update(bson.M{"username":username}, bson.M{"$set": bson.M{"group": gr}})
    return err
}

func UpdateBuildingGroup(building string, room string, group string) error {
    gr, _ := strconv.Atoi(group)
    err := DB.C("buildings").Update(bson.M{"name":building}, bson.M{"$set": bson.M{"rooms." + room + ".group": gr}})
    return err
}

func UpdateRoomStatus(building string, num string, room *Room) error {
    err := DB.C("buildings").Update(bson.M{"name":building}, bson.M{"$set": bson.M{"rooms." + num + "": room}})
    return err
}

func ResetAllUserGroups() error {
    users, _ := GetAllUsers()
    for index, _ := range users {
        users[index].Group = 0
    }

    _, err := DB.C("users").RemoveAll(nil)
    if err != nil {
        log.Println(err)
        return err
    }

    for _, data := range users {
        err := DB.C("users").Insert(&data)
        if err != nil {
            log.Println(err)
            return err
        }
    }

    return nil
}

func ResetAllRoomGroups() error {
    b, _ := GetRooms()

    reset := make([]Building, len(b))

    for index, data := range b {
        reset[index].Name = data.Name
        reset[index].Abrv = data.Abrv
        reset[index].Rooms = make(map[string]Room)
        for num, info := range data.Rooms {
            var temp Room
            temp.Status = info.Status
            temp.Probcat = info.Probcat
            temp.Notes = info.Notes
            temp.Group = 0
            temp.Lamph.Standard = info.Lamph.Interactive
            if info.Lamph.Interactive != -1 {
                temp.Lamph.Interactive = info.Lamph.Interactive
            } else {
                temp.Lamph.Interactive = -1
            }
            reset[index].Rooms[num] = temp
        }
    }

    _, err := DB.C("buildings").RemoveAll(nil)
    if err != nil {
        log.Println(err)
        return err
    }

    for _, data := range reset {
        err := DB.C("buildings").Insert(&data)
        if err != nil {
            log.Println(err)
            return err
        }
    }

    return nil
}

func ResetRoom(bdngName string, num string) {
    room, _ := GetRoom(bdngName, num)

    reset := Room {Status: "unchecked", Probcat: []int{0}, Notes: "", Group: room.Group}
    reset.Lamph.Standard = 0
    if room.Lamph.Interactive != -1 {
        reset.Lamph.Interactive = 0
    } else {
        reset.Lamph.Interactive = -1
    }

    UpdateRoomStatus(bdngName, num, &reset)
}

func ResetAllRooms() error{
    b, _ := GetRooms()

    reset := make([]Building, len(b))

    for index, data := range b {
        reset[index].Name = data.Name
        reset[index].Abrv = data.Abrv
        reset[index].Rooms = make(map[string]Room)
        for num, info := range data.Rooms {
            var temp Room
            temp.Status = "unchecked"
            temp.Probcat = []int{0}
            temp.Notes = ""
            temp.Group = info.Group
            temp.Lamph.Standard = 0
            if info.Lamph.Interactive != -1 {
                temp.Lamph.Interactive = 0
            } else {
                temp.Lamph.Interactive = -1
            }
            reset[index].Rooms[num] = temp
        }
    }

    _, err := DB.C("buildings").RemoveAll(nil)
    if err != nil {
        log.Println(err)
        return err
    }

    for _, data := range reset {
        err := DB.C("buildings").Insert(&data)
        if err != nil {
            log.Println(err)
            return err
        }
    }

    return nil
}
