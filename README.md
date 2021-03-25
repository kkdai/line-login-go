LINE loing in Go: Sample code for LINE login in Go
==============

 [![GoDoc](https://godoc.org/github.com/kkdai/line-login-go.svg?status.svg)](https://godoc.org/github.com/kkdai/line-login-go)[![goreportcard.com](https://goreportcard.com/badge/github.com/kkdai/line-login-go)](https://goreportcard.com/report/github.com/kkdai/line-login-go)
 ![Go](https://github.com/kkdai/line-login-go/workflows/Go/badge.svg)


![](https://developers.line.biz/media/line-login/integrate-login-web/login-flow-web-0bc4c99d.png)

Refer LINE Developer Document "[Integrating LINE Login with your web app](https://developers.line.biz/en/docs/line-login/web/integrate-line-login/)" for more detail.


This sample code implement how to integrate LINE login to your website in Go. You can use this sample code to integrate LINE login in your web application. Also this also provide link a chatbot service when user use LINE login. Refer "[Linking a bot with your LINE Login channel](https://developers.line.biz/en/docs/line-login/web/link-a-bot/)".

Deploy on Heroku
=============

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Before deploy this to your Heroku, you will need complete as follows:

- Create a LINE login channel. Remember it's channel ID and channel secret.
- Crete a LINE Message API channel. Remember it's channel secret and token.
- Link the chatbot to the LINE login channel

License
=============

MIT License
