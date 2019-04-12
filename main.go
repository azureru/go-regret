package main

import (
	"flag"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ConsumerKey string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	OauthToken string `yaml:"oauth_token"`
	OauthTokenSecret string `yaml:"oauth_token_secret"`
}

func main() {

	flag.Parse()

	var configStruct Config

	// read config first
	buffer, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic("./config.yaml is needed!")
	}
	err = yaml.Unmarshal(buffer, &configStruct)
	if err != nil {
		panic("invalid ./config.yaml file")
	}

	// connect
	config := oauth1.NewConfig(configStruct.ConsumerKey, configStruct.ConsumerSecret)
	token := oauth1.NewToken(configStruct.OauthToken, configStruct.OauthTokenSecret)
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

	// the whole thing
}
