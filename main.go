package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/imroc/req/v3"
)

type Config struct {
	Wordlist     string
	User_Agent   string
	Ratelimit    int
	RequestDelay int
}

func main() {
	conf := configinit()
	filecontent, err := os.ReadFile(conf.Wordlist)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	wordlist := strings.Split(string(filecontent), "\n")
	client := req.C()

	for _, words := range wordlist {
		checkign(strings.ReplaceAll(words, "\r", ""), client, conf)
		time.Sleep(time.Duration(conf.RequestDelay) * time.Second)
	}
}

func configinit() *Config {

	var conf Config

	TomlData, err := os.ReadFile("./config.toml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	toml.Decode(string(TomlData), &conf)
	return &conf
}

func checkign(username string, client *req.Client, conf *Config) {
	resp, err := client.R().SetHeader("Accept", "*/*").SetHeader("User-Agent", conf.User_Agent).Get(fmt.Sprintf("https://github.com/signup_check/username?value=%s", username))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch resp.StatusCode {
	case 200:
		fmt.Printf("[FREE] %s\n", username)
	case 422:
		/* For Debug Purposes
		Ignore this Line Otherwise */

		//fmt.Printf("[NOT FREE] %s\n", username)
		return
	case 429:
		fmt.Println("==== RATE LIMIT COOLDOWN... ====")
		time.Sleep(time.Duration(conf.Ratelimit) * time.Second)
		fmt.Println("==== CONTINUE... ====")

	default:
		fmt.Println(resp.StatusCode)
		fmt.Println(resp)
		return
	}

}
