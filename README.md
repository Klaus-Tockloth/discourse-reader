# discourse-reader

## Context
Go, forum, discourse, download data, site, category, topic, posts

## Purpose
Retrieves data (e.g. site, category, topic) from Discourse forum on behalf of user.

## General concept
'discourse-reader' retrieves data from an arbitrary forum based on Discourse. E.g. [OpenStreetMap](https://community.openstreetmap.org) or [Meta Discourse](https://meta.discourse.org).

The typical usage of 'discourse-reader' is to retrieve all posts of a topic (thread) for further processing.

'discourse-reader' does functionally no different than a user via a browser. However, the data is retrieved somewhat faster.
    
## Prerequisites
In order to retrieve data, you must have a user and a user-api-key for the forum. All actions are performed on behalf of the user. The companion program [discourse-user-api-key](https://github.com/Klaus-Tockloth/discourse-user-api-key) helps to create a user-api-key. The user-api-key can be set as command line option or via environment variable.

## Binaries
Precompiled binaries for many operating systems can be found here : **Releases -> Assets**

## Read all posts of a topic (thread)
Often the posts of a topic are very valuable. The main purpose of 'discourse-reader' is to download all posts for subsequent processing. A topic can have much very posts. This is the download algorithms:

```
first read  : meta data
second read : 50 posts
third read  : 50 posts
... 
```
Downloading can be controlled with the options '-pages' and '-sleeptime'. A forums server typically limits downloads (user rate limits). The '-sleeptime' option waits some time (in seconds) between subsequent read operations. Example: If the server allows 20 read operations per minute, the option '-sleeptime=6' would be a good choice in order to prevent user rate limit rejections. You limit the read operations (60/6=10) to 10 per minute. Be fair to the server and therefore also to all other users.

## Internet proxy
'discourse-reader' also works behind a (company) internet proxy. See remarks section in the help output.

## Program usage
```
discourse-reader 
Program:
  Name    : discourse-reader
  Release : v1.0.0 - 2022/11/18
  Purpose : Discourse Reader
  Info    : Retrieves data (e.g. site, category, topic) from Discourse forum.

Usage:
  discourse-reader -forum=string -query=string -category=int -topic=int -pages=int -output=string -userapikey -sleeptime=int

Examples for general query:
  discourse-reader
  discourse-reader -query=community.openstreetmap.org/site.json -output=community.openstreetmap.org.json
  discourse-reader -query=community.openstreetmap.org/site.json -output=community.openstreetmap.org.json -userapikey=bd38603815e3f2562c3eb3988c69eb77
  discourse-reader -query=meta.discourse.org/site.json -output=meta.discourse.org.json
  discourse-reader -query=meta.discourse.org/session/current.json -output=session-current.json

Examples for category:
  discourse-reader -forum=community.openstreetmap.org -category=56 -output=category-56.json
  discourse-reader -forum=community.openstreetmap.org -category=56 -output=category-56.json -userapikey=bd38603815e3f2562c3eb3988c69eb77
  discourse-reader -forum=meta.discourse.org -category=67 -pages=99 -sleeptime=6 -output=category-67.json

Examples for topic:
  discourse-reader -forum=community.openstreetmap.org -topic=4120 -output=topic-4120.json
  discourse-reader -forum=community.openstreetmap.org -topic=4120 -pages=99 -sleeptime=6 -output=topic-4120.json
  discourse-reader -forum=community.openstreetmap.org -topic=4120 --output=topic-4120.json -userapikey=bd38603815e3f2562c3eb3988c69eb77
  discourse-reader -forum=meta.discourse.org -topic=112837 -output=topic-112837.json

Options:
  -category int
    	retrieve data (list of topics) for category with identifier (default -1)
  -forum string
    	Discourse forum URL
  -output string
    	name of JSON output file
  -pages int
    	pages of data to retrieve (default 19)
  -query string
    	general data retrieve query (full URL)
  -sleeptime int
    	sleep time in seconds before retrieving the next page (avoids user rate limiting) (default 2)
  -topic int
    	retrieve data (list of posts) for topic with identifier (default -1)
  -userapikey string
    	personal user API key (can also be set as environment var 'USER_API_KEY')

Remarks:
  - User API key can be set as environment variable [USER_API_KEY].
  - Internet proxy can be set as environment variable [HTTPS_PROXY].
  - Examples for Linux:
    export USER_API_KEY=bd38603815e3f2562c3eb3988c69eb77
    export HTTPS_PROXY=http://user:password@194.114.63.23:8080
  - Examples for Windows:
    set USER_API_KEY=bd38603815e3f2562c3eb3988c69eb77
    set HTTPS_PROXY=http://user:password@194.114.63.23:8080

Rate limiting by forum service:
  - This program does functionally no different than a user via a browser. However, the
    data is retrieved somewhat faster. This can lead to rejections (rate limiting) by the
    service. To prevent this, the program can pause between fetching pages. The pause time
    can be specified with the option '-sleeptime=int'.
  - Typical user rate limit settings are:
    - requests per minute : 20
    - requests per day    : 2880
```

