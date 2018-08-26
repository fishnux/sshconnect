package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//c := config{[]tag{{"Work", []server{{"Box1", "127.0.0.1"}}}}}
	//b, err := json.Marshal(c)
	//fmt.Println(string(b), err)
	b, err := ioutil.ReadFile("./config.json")
	check(err)
	err = json.Unmarshal(b, &configFile)
	check(err)
	fmt.Println("Here are your servers:")
	serverCount := 0
	for _, el := range configFile.Tags {
		fmt.Println(aurora.Red(fmt.Sprintf(" - %s", el.Name)))
		for _, el2 := range el.Servers {
			servers = append(servers, el2)
			serverCount++
			fmt.Printf("%v.    %s (%s)\n", serverCount, el2.Name, el2.Address)
		}
	}
	fmt.Println("What server would you like to connect to?")
	var input string
	fmt.Scanln(&input)
	inputNumeric, err := strconv.ParseInt(input, 10, 16)
	check(err)
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

	// Here's the actual `syscall.Exec` call. If this call is
	// successful, the execution of our process will end
	// here and be replaced by the `/bin/ls -a -l -h`
	// process. If there is an error we'll get a return
	// value.
	execErr := syscall.Exec(binary, args, env)
	check(execErr)
}
