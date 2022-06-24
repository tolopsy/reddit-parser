reddit-parser is an api server that parses posts in a subreddit into json response.

## What you should know
### How to set the port
 - By default, The server listens to port :9000
 - Use the `port` flag to specify a different port to listen to
 
 ### How to specify subreddit to parse
 - In GET request to server's root url, pass the subreddit's url as a query `url` parameter.
	 For example: `localhost:9000?url=https://www.reddit.com/r/news/`
 - In POST request to server's root url, pass the subreddit's url via the request's body in json format
	For example:
	```
	{
		"url": "https://www.reddit.com/r/news/"
	}
	```
	