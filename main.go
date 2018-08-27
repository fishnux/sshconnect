package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"

	"github.com/logrusorgru/aurora"
)

type config struct {
	Tags []tag
}

type tag struct {
	Name    string
	Servers []server
}

type server struct {
	Name    string
	Address string
}

var configFile config
var servers []server
var configLocation = "~/.sshconnect.json"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimRight(input, "\n")
}

func main() {
	configLocation, _ = homedir.Expand(configLocation)
	b, err := ioutil.ReadFile(configLocation)
	if err != nil {
		json := []byte("{}")
		err = ioutil.WriteFile(configLocation, json, 0700)
		check(err)
		b = json
	}
	err = json.Unmarshal(b, &configFile)
	check(err)
	serverCount := 0
	for _, el := range configFile.Tags {
		fmt.Println(aurora.Red(fmt.Sprintf(" - %s", el.Name)))
		for _, el2 := range el.Servers {
			servers = append(servers, el2)
			serverCount++
			fmt.Printf("%v.    %s (%s)\n", serverCount, el2.Name, el2.Address)
		}
	}
	if serverCount != 0 {
		fmt.Print("What server would you like to connect to? ")
	} else {
		fmt.Println("Looks like you don't have any servers at all. Please type \"add\":")
	}
	input := readInput()
	inputNumeric, err := strconv.ParseInt(input, 10, 16)
	if err == nil {
		if inputNumeric < 1 || inputNumeric > int64(len(servers)) {
			fmt.Println("Invalid number")
			os.Exit(1)
		}
		inputNumeric-- // Fix off-by-one error
		str := fmt.Sprintf("Connecting to %s (%s)...", servers[inputNumeric].Name, servers[inputNumeric].Address)
		fmt.Println(aurora.Red(str))
		binary, lookErr := exec.LookPath("ssh")
		check(lookErr)

		args := []string{"ssh", servers[inputNumeric].Address}

		// `Exec` also needs a set of [environment variables](environment-variables)
		// to use. Here we just provide our current
		// environment.
		env := os.Environ()

		execErr := syscall.Exec(binary, args, env)
		check(execErr)
	} else {
		if input == "add" {
			fmt.Println("What's the server name?")
			serverName := readInput()
			fmt.Println("What's the server address?")
			serverAddress := readInput()
			fmt.Println("What tag does it belongs to? (e.g.: Work)")
			serverTag := readInput()
			newServer := server{serverName, serverAddress}
			serverTagFound := -1
			for i := range configFile.Tags {
				if configFile.Tags[i].Name == serverTag {
					// Found!
					serverTagFound = i
					break
				}
			}
			if serverTagFound == -1 {
				configFile.Tags = append(configFile.Tags, tag{serverTag, []server{}})
				serverTagFound = len(configFile.Tags) - 1
			}
			configFile.Tags[serverTagFound].Servers = append(configFile.Tags[serverTagFound].Servers, newServer)
			json, err := json.MarshalIndent(configFile, "", "    ")
			check(err)
			err = ioutil.WriteFile(configLocation, json, 0700)
			check(err)
			fmt.Println("Done, please re-run this program")
		} else {
			fmt.Println("Invalid command, exiting...")
			os.Exit(0)
		}
	}
}
