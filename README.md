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
