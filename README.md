A tool to download all a user's twitter followers


# Installation

Ensure you have go installed:

```
# e.g. on a Mac
brew install go
```

Install the `twitter-followers` script:


```
go install github.com/ConradIrwin/twitter-followers
```

# Usage

```
twitter-followers <screen_name>
```

This will output a json object per line for each follower. It uses the same
ordering as the Twitter API, which is approximately reverse chronological order
of when they followed you.

I suggest pairing this with [jq](https://stedolan.github.io/jq/) to extract the
fields you need.

