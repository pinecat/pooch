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
)

/* structs */
type User struct {
    ID          bson.ObjectId `bson:"_id,omitempty"`
    Username    string
    Password    string
    Fname       string
    Lname       string
    Type        string
    Team        int
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
}

type Assign struct {
    Group   int
    Team    int
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
    if err != nil {
        return result, err
    }
    return result, nil
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
    if err != nil {
        return err
    }
    return nil
}

func RemoveUser(username string) error {
    err := DB.C("users").Remove(bson.M{"username":username})
    if err != nil {
        return err
    }
    return nil
}

func GetRooms() ([]Building, error) {
    c := DB.C("buildings")
    var result []Building
    err := c.Find(nil).All(&result)
    if err != nil {
        return result, err
    }
    return result, nil
}

func GetAssignments() ([]Assign, error) {
    c := DB.C("assign")
    var result []Assign
    err := c.Find(nil).All(&result)
    if err != nil {
        return result, err
    }
    return result, nil
}
