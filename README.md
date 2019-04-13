# go-regret

This is a CLI tool to remove old social media posts (twitter) 
This is to encourage people to live in the moment and quickly move on from the past
(and as we are getting older - to remove things that we may regret when we were younger)

Me personally, I use a raspberry pi - call `clean` daily - 
so my tweets are always cleaned up :) - I can't trust my 30-days previous me :P

## Usage
- You need to get customer key, secret, oauth token and oauth secret
- Copy `config.sample.yaml` to `config.yaml`
- Edit `config.yaml` and fill the appropriate fields
- Run the cli `go-regret` with specified mode below

## Clean 
```
    ./go-regret -mode=clean -delta=30
```
This mode will clean tweet that are older than N number of days specified on `-delta` param
This will use twitter search (so the max number of previous tweets are limited)

If you call this daily (e.g. using cron) this will basically keep your tweet always within
30 days

## Purge
```
    ./go-regret -mode=purge -file=/blah/blah/tweet.js
```
Will clean tweets that listed on `tweet.js` (you can get this tweet.js on Twitter Archive Takeout)
Since clean is limited - you will need this mode to purge super-duper old tweet

## Additional Params 
```
    -dry=1 
    When specified will only run the deletion on dry-mode, so you can test the effect of your config first     
```

## Installation for Cron 
On my machine I install the binary by:
```
 mkdir /opt/go-regret
 cd /opt/go-regret 
 wget https://github.com/azureru/go-regret/releases/download/0.1/go-regret_linux_amd64.zip
 unzip go-regret_linux_amd64.zip
 rm go-regret_linux_amd64.zip
```

Create a config.yaml file on the same directory, of course I need to change values in {REPLACE_ME}
by value that I get on twitter developer page

```
# twitter keys
consumer_key : "{REPLACE_ME}"
consumer_secret : "{REPLACE_ME}"
oauth_token : "{REPLACE_ME}"
oauth_token_secret : "{REPLACE_ME}"

# do not purge when retweet exceed this value
retweet_count : 1

# do not purge when like exceed this value
like_count : 1

# if 1 - we also purge our replies
purge_reply : 0

``` 

Create a cron-entry
```
0 18 * * * cd /opt/go-regret && ./go-regret_linux_amd64 -mode=clean -delta=30 -confirm=y
```

Profit!