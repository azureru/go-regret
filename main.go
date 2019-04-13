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
	"time"
)

var Mode string
var GlobalUserId int64
var DryRun bool
var DeltaDays int

type Config struct {
	ConsumerKey      string `yaml:"consumer_key"`
	ConsumerSecret   string `yaml:"consumer_secret"`
	OauthToken       string `yaml:"oauth_token"`
	OauthTokenSecret string `yaml:"oauth_token_secret"`

	RetweetCount int `yaml:"retweet_count"`
	LikeCount    int `yaml:"like_count"`
	PurgeReply   int `yaml:"purge_reply"`
}

type TweetFromArchive struct {
	IdStr         string `json:"id_str"`
	FavoriteCount string `json:"favorite_count"`
	RetweetCount  string `json:"retweet_count"`
	FullText      string `json:"full_text"`
	CreatedAt     string `json:"created_at"`
}

type TweetsFromArchive []TweetFromArchive

func main() {
	// mode is required
	mode := flag.String("mode", "", "available mode are `clean` or `purge`")
	tweetJsFile := flag.String("file", "", "complete file path of tweet.js")
	deltaDays := flag.Int("delta", 0, "number of days - to delete tweets that are older than this value")
	dryRun := flag.Bool("dry", false, "when true - will only show tweets - will not do the actual delete operation")
	flag.Parse()

	Mode = *mode
	DeltaDays = *deltaDays
	DryRun = *dryRun

	// validate some required specifiers
	if Mode == "" {
		fmt.Println("\t-mode is required, use `clean` or `purge` mode")
		os.Exit(1)
	}
	if Mode == "clean" {
		// on clean mode -delta is required to be more than 0
		if DeltaDays <= 0 {
			fmt.Println("\t-delta need to be set with number of days")
			os.Exit(1)
		}
	}
	if Mode == "purge" {
		// on purge mode -file is required
		if *tweetJsFile=="" {
			fmt.Println("-file tweet.js path is required for purge mode")
			os.Exit(1)
		}
	}

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

	if DryRun {
		fmt.Println("\tDry Run Mode - No Actual Deletion will be executed")
	}

	// the whole thing
	if Mode == "clean" {
		clean(client, DeltaDays)
	} else if Mode == "purge" {
		purge(client, *tweetJsFile)
	}
}

// deleteTweet - to delete a tweet (or show the tweet if dry-run is specified)
func deleteTweet(client *twitter.Client, tweet twitter.Tweet) {
	whenTweet, _ := tweet.CreatedAtTime()

	fmt.Println()
	fmt.Println(whenTweet.Format(time.RubyDate)+"\t", tweet.ID)
	fmt.Println("\t" + tweet.FullText)
	fmt.Println("\t❤FAV:", tweet.FavoriteCount, "RT:", tweet.RetweetCount)

	// do the actual deletion
	if !DryRun {
		_, _, err := client.Statuses.Destroy(tweet.ID, nil)
		if err != nil {
			fmt.Println("\tcannot destroy", tweet.ID, err)
		} else {
			fmt.Println("\tdeleting ", tweet.ID)
		}
	}
}

// clean will cleanup all tweets older thant delta number of days
func clean(client *twitter.Client, delta int) {
	var maxId int64

	// number of days
	var deltaDays int64
	deltaDays = int64(delta) * 86400
	deltaUnixTime := time.Now().Unix() - deltaDays

	// initial tweets run
	for {
		params := &twitter.UserTimelineParams{
			UserID: GlobalUserId,
			Count:  200,
		}
		if maxId != 0 {
			params.MaxID = maxId
		}
		fmt.Println(params)
		tweets, _, err := client.Timelines.UserTimeline(params)
		if err != nil {
			fmt.Println("\tcannot lookup timeline", err)
			break
		}
		if len(tweets) > 0 {
			for _, tweet := range tweets {
				createdAt, _ := tweet.CreatedAtTime()
				if createdAt.Unix() < deltaUnixTime {
					// the tweet is older than specified delta - remove it
					fmt.Println(tweet)
				}
			}
			maxId = tweets[len(tweets)-1 ].ID - 1
		} else {
			// loop until the end of lookup
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
	var tweetBase TweetsFromArchive
	strBuffer := strings.Replace(string(buffer), "window.YTD.tweet.part0 = ", "", 1)
	err = json.Unmarshal([]byte(strBuffer), &tweetBase)
	if err != nil {
		panic("cannot parse tweet.js JSON")
	}

	fmt.Println("Tweet.js")
	fmt.Println("\tTweets: ", len(tweetBase))
	for _, el := range tweetBase {
		tw := tweetArchiveToTweet(el)
		deleteTweet(client, tw)
	}
}

// tweetArchiveToTweet convert tweet from archive to twitter.Tweet object
func tweetArchiveToTweet(el TweetFromArchive) twitter.Tweet {

	favInt, _ := strconv.Atoi(el.FavoriteCount)
	retweetInt, _ := strconv.Atoi(el.RetweetCount)
	idInt64, _ := strconv.ParseInt(el.IdStr, 10, 64)

	// Fri Jun 08 18:21:18 +0000 2018
	//createdAt, _ := time.Parse(time.RubyDate,el.CreatedAt)

	twit := twitter.Tweet{
		ID: idInt64,
		IDStr : el.IdStr,
		FullText: el.FullText,
		FavoriteCount: favInt,
		RetweetCount: retweetInt,
		CreatedAt: el.CreatedAt,
	}

	return twit
}
