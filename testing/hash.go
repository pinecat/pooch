package main

import "fmt"
import "golang.org/x/crypto/bcrypt"

func main() {
    pass := "golang"
    hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    fmt.Println("Hashed pass: ", string(hash), "!")
}
