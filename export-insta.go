// +build ignore

package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"gopkg.in/ahmdrz/goinsta.v2"
)

func main() {
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	inst := goinsta.New(env["IG_USERNAME"], env["IG_PASSWORD"])
	err = inst.Login()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	err = inst.Export(env["IG_COOKIE_PATH"])
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
}
