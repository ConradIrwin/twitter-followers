package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	if os.Getenv("TWITTER_CONSUMER_KEY") == "" ||
		os.Getenv("TWITTER_CONSUMER_SECRET") == "" ||
		os.Getenv("TWITTER_ACCESS_KEY") == "" ||
		os.Getenv("TWITTER_ACCESS_SECRET") == "" {

		fmt.Println(`Please sign up for a twitter app, and export the following environment variables:

			export TWITTER_CONSUMER_KEY=""
			export TWITTER_CONSUMER_SECRET=""
			export TWITTER_ACCESS_KEY=""
			export TWITTER_ACCESS_SECRET=""
			twitter-followers SuperhumanCo
		`)

		os.Exit(1)
	}

	cursorFile := flag.String("cursor-file", "", "a file to save the state of the cursor in")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		fmt.Printf("%#v", args)
		fmt.Println("Usage: twitter-followers [ --cursor-file .cursor] <screen_name>")
		os.Exit(1)
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	client := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_KEY"), os.Getenv("TWITTER_ACCESS_SECRET"))

	cursorStr := "-1"

	if *cursorFile != "" {
		bytes, err := ioutil.ReadFile(*cursorFile)
		if err == nil {
			cursorStr = string(bytes)
		}
	}

	v := url.Values{}
	v.Set("screen_name", args[0])
	v.Set("count", "200")
	v.Set("skip_status", "true")
	v.Set("cursor", cursorStr)

	cursor, err := client.GetFollowersList(v)
	if err != nil {
		panic(err)
	}

	for {
		if len(cursor.Users) == 0 {
			break
		}

		for _, u := range cursor.Users {
			bytes, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}
			os.Stdout.Write(bytes)
			os.Stdout.Write([]byte{'\n'})

			time.Sleep(70 * time.Second / time.Duration(len(cursor.Users)))
		}

		if *cursorFile != "" {
			err := ioutil.WriteFile(*cursorFile, []byte(cursor.Next_cursor_str), 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "writing to cursor file failed: %v %v", cursor.Next_cursor_str, err)
			}
		}

		v.Set("cursor", cursor.Next_cursor_str)

		cursor, err = client.GetFollowersList(v)
		if err != nil {
			panic(err)
		}
	}

}
