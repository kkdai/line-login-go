#!/usr/bin/env bash
# Deploy to Cloud Run from source (Cloud Build uses the Dockerfile).
# Usage:
#   export PROJECT_ID=my-project
#   export LINE_LOGIN_CHANNEL_ID=... LINE_LOGIN_CHANNEL_SECRET=...
#   export LINE_BOT_CHANNEL_SECRET=... LINE_BOT_CHANNEL_TOKEN=...
#   ./deploy-cloudrun.sh
set -euo pipefail

: "${PROJECT_ID:?set PROJECT_ID}"
SERVICE="${SERVICE:-line-login-go}"
REGION="${REGION:-asia-east1}"

gcloud config set project "$PROJECT_ID"
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com

# First deploy: SERVERURL unknown yet, use placeholder, then update with the real URL.
ENV_VARS="LINECORP_PLATFORM_CHANNEL_CHANNELID=${LINE_LOGIN_CHANNEL_ID:?},LINECORP_PLATFORM_CHANNEL_CHANNELSECRET=${LINE_LOGIN_CHANNEL_SECRET:?},LINECORP_PLATFORM_CHATBOT_CHANNELSECRET=${LINE_BOT_CHANNEL_SECRET:?},LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN=${LINE_BOT_CHANNEL_TOKEN:?}"

gcloud run deploy "$SERVICE" \
  --source . \
  --region "$REGION" \
  --allow-unauthenticated \
  --set-env-vars "$ENV_VARS,LINECORP_PLATFORM_SERVERURL=${SERVER_URL:-https://placeholder.invalid}"

URL=$(gcloud run services describe "$SERVICE" --region "$REGION" --format 'value(status.url)')
gcloud run services update "$SERVICE" --region "$REGION" \
  --update-env-vars "LINECORP_PLATFORM_SERVERURL=$URL"

echo
echo "Deployed: $URL"
echo "Set in LINE Developers Console:"
echo "  LINE Login Callback URL: $URL/auth"
echo "  Messaging API Webhook:   $URL/callback"
