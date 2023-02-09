// Copyright 2022 Red Hat, Inc.
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package myquery

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

type UserInfo struct {
	Name string
	Pass string
}

var User UserInfo

// Display login screen and get user info.
func GetUserInfo() {

	greetings()
	for {
		fmt.Print("Username: ")

		_, err := fmt.Scanln(&User.Name)
		if err == nil {
			break
		}
		logrus.Error("Username cannot be empty, please try again")
	}
	fmt.Print("Password: ")
	bytepw, passerr := term.ReadPassword(int(syscall.Stdin))

	fmt.Print("\n")
	s := cli.IniSpinner(" Getting user credentials")
	if passerr != nil {
		logrus.Fatal("Error getting password", passerr)
	}
	s.Stop()
	User.Pass = string(bytepw)
}

func greetings() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[H")

	greetings := "Please, enter your ldap credentials below."
	hi := `
  ███████╗██╗   ██╗███╗   ██╗ ██████╗██████╗  ██████╗ ███╗   ██╗
  ██╔════╝╚██╗ ██╔╝████╗  ██║██╔════╝██╔══██╗██╔═══██╗████╗  ██║
  ███████╗ ╚████╔╝ ██╔██╗ ██║██║     ██████╔╝██║   ██║██╔██╗ ██║
  ╚════██║  ╚██╔╝  ██║╚██╗██║██║     ██╔══██╗██║   ██║██║╚██╗██║
  ███████║   ██║   ██║ ╚████║╚██████╗██║  ██║╚██████╔╝██║ ╚████║
  ╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
         Your tool for everything usage patterns related                                       
`
	fmt.Println(hi)
	fmt.Println("Welcome!")
	fmt.Println(greetings)
	fmt.Println(strings.Repeat("-", len(greetings)))
}
