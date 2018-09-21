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
        if i != num {
            err = errors.New("Room not found!")
        } else {
            r = data
        }
    }
    return r, err
}

func GetRooms() ([]Building, error) {
    c := DB.C("buildings")
    var result []Building
    err := c.Find(nil).All(&result)
    return result, err
}

func UpdateGroup(building string, room string, group int) error {
    err := DB.C("buildings").Update(bson.M{"name":building}, bson.M{"$set": bson.M{"rooms." + room + ".group": 0}})
    return err
}
