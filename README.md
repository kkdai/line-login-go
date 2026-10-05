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

Deploy on Google Cloud Run
=============

Prerequisites: the same LINE channels as above, plus `gcloud` installed and logged in (`gcloud auth login`).

```
export PROJECT_ID=your-gcp-project
export LINE_LOGIN_CHANNEL_ID=... LINE_LOGIN_CHANNEL_SECRET=...
export LINE_BOT_CHANNEL_SECRET=... LINE_BOT_CHANNEL_TOKEN=...
./deploy-cloudrun.sh   # optional: REGION (default asia-east1), SERVICE (default line-login-go)
```

The script builds from the `Dockerfile` with Cloud Build, deploys to Cloud Run, and sets `LINECORP_PLATFORM_SERVERURL` to the service URL. Afterwards, set these in the LINE Developers Console:

- LINE Login Callback URL: `<service URL>/auth`
- Messaging API Webhook URL: `<service URL>/callback`

### Deploy manually with gcloud

The script above wraps these commands:

```
gcloud config set project $PROJECT_ID
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com

# 1. Build from the Dockerfile with Cloud Build and deploy
gcloud run deploy line-login-go \
  --source . \
  --region asia-east1 \
  --allow-unauthenticated \
  --set-env-vars LINECORP_PLATFORM_CHANNEL_CHANNELID=your_channel_id,LINECORP_PLATFORM_CHANNEL_CHANNELSECRET=your_channel_secret,LINECORP_PLATFORM_CHATBOT_CHANNELSECRET=your_bot_secret,LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN=your_bot_token,LINECORP_PLATFORM_SERVERURL=https://placeholder.invalid

# 2. Get the service URL and set it as SERVERURL (used to build the /auth redirect URI)
URL=$(gcloud run services describe line-login-go --region asia-east1 --format 'value(status.url)')
gcloud run services update line-login-go --region asia-east1 \
  --update-env-vars LINECORP_PLATFORM_SERVERURL=$URL
```

Note: `status.url` may differ from the other URL form shown in the deploy output (`https://<service>-<project-number>.<region>.run.app`). Use the same URL in `LINECORP_PLATFORM_SERVERURL` and in the LINE console, otherwise the redirect URI will not match.

To update a single variable later:

```
gcloud run services update line-login-go --region asia-east1 --update-env-vars KEY=VALUE
```

To view logs: `gcloud run services logs read line-login-go --region asia-east1`

For production, consider storing secrets in Secret Manager and using `--set-secrets` instead of `--set-env-vars`.

Run In Docker
=============

```
docker build -t line-login-go:latest .

docker run \
-e LINECORP_PLATFORM_CHATBOT_CHANNELSECRET=your_secret \
-e LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN=your_token \
-e LINECORP_PLATFORM_SERVERURL=http://localhost:8080 \
-e LINECORP_PLATFORM_CHANNEL_CHANNELID=your_channel_id \
-e LINECORP_PLATFORM_CHANNEL_CHANNELSECRET=your_channel_secret \
-e PORT=8080 \
-p 8080:8080 \
line-login-go
```

License
=============

MIT License
