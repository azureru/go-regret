# go-regret

This is a CLI tool to remove social media posts (twitter) -
to encourage people to live in the moment and quickly move on from the past
(and as we are getting older - to remove things that we may regret when we were younger)

## Usage
- You need to get customer key, secret, oauth token and oauth secret
- Copy `config.sample.yaml` to `config.yaml`
- Edit `config.yaml` and fill the appropriate fields
- Run the cli `go-regret`

## Clean 
Will clean tweet that are older than N number of days (by default 30 days)
This will use twitter search (so the max number of tweet are limited)

## Purge
Will clean tweets that listed on `tweet.js` (you can get this js on Twitter Archive Takeout)
