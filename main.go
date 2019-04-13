package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var GlobalUserId int64

type Config struct {
	ConsumerKey      string `yaml:"consumer_key"`
	ConsumerSecret   string `yaml:"consumer_secret"`
	OauthToken       string `yaml:"oauth_token"`
	OauthTokenSecret string `yaml:"oauth_token_secret"`
}

type Tweet struct {
	IdStr         string `json:"id_str"`
	FavoriteCount string `json:"favorite_count"`
	RetweetCount  string `json:"retweet_count"`
	FullText      string `json:"full_text"`
}

type Tweets []Tweet

func main() {
	// mode is required
	mode := flag.String("mode", "clean", "[clean|purge] mode")
	tweetJsFile := flag.String("file", "", "the file path of tweet.js")
	flag.Parse()

	var configBase Config

	// read config first
	buffer, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic("./config.yaml is needed!")
	}
	err = yaml.Unmarshal(buffer, &configBase)
	if err != nil {
		panic("invalid ./config.yaml file")
	}

	// connect
	config := oauth1.NewConfig(configBase.ConsumerKey, configBase.ConsumerSecret)
	token := oauth1.NewToken(configBase.OauthToken, configBase.OauthTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		panic("Invalid credentials, please check ./config.yaml value")
	}
	fmt.Printf("Your Account:\nID: %+v\nHandle: @%+v\nName: %+v\n", user.ID, user.ScreenName, user.Name)
	GlobalUserId = user.ID

	// the whole thing
	modeStr := *mode
	if modeStr == "clean" {

	} else if modeStr == "purge" {
		purge(client, *tweetJsFile)
	}
}

// clean will cleanup all tweets older thant delta number of days
func clean(client *twitter.Client, delta int) {
	var maxId int64

	// initial tweets run
	for {
		params := &twitter.UserTimelineParams{
			UserID: GlobalUserId,
			Count:  200,
		}
		if maxId != 0 {
			params.MaxID = maxId
		}
		tweets, _, err := client.Timelines.UserTimeline(params)
		if err != nil {
			fmt.Println("\tcannot lookup timeline", err)
			break
		}
		if len(tweets) > 0 {
			for idx, tweet := range tweets {
				_, _ , err := client.Statuses.Destroy(tweet.ID, nil)
				if err != nil {
					fmt.Println("\tcannot destroy", tweet.ID, err)
				} else {
					fmt.Println("\tdeleting ", tweet.ID)
				}
			}
			maxId = tweets[len(tweets) -1 ].ID
		} else {
			fmt.Println("\treach end of max tweet to lookup")
			break
		}
	}
}

// purge will read tweet.js on -file argument and iterate each status and try to delete it
func purge(client *twitter.Client, tweetJsFile string) {
	fmt.Println(tweetJsFile)
	if tweetJsFile == "" {
		fmt.Println("\t-file argument is required for purge mode")
		os.Exit(1)
	}
	buffer, err := ioutil.ReadFile(tweetJsFile)
	if err != nil {
		fmt.Println(err)
		panic("invalid tweet.js path on --file as purge source")
	}
	var tweetBase Tweets
	strBuffer := strings.Replace(string(buffer), "window.YTD.tweet.part0 = ", "", 1)
	err = json.Unmarshal([]byte(strBuffer), &tweetBase)
	if err != nil {
		panic("cannot parse tweet.js JSON")
	}

	fmt.Println("Tweet.js")
	fmt.Println("\tTweets: ", len(tweetBase))
	for _, el := range tweetBase {
		intId, _ := strconv.ParseInt(el.IdStr, 10, 64)
		_, _ , err := client.Statuses.Destroy(intId, nil)
		if err != nil {
			fmt.Println("\tcannot destroy", intId, err)
		} else {
			fmt.Println("\tdeleting ", intId)
		}
	}
}
