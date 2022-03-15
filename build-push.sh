#!/bin/bash

if [[ "$1" -eq "" ]]
then 
    SERVICE=appsheet-docai
else
    SERVICE="$1"
fi

if [[ "$DEVSHELL_PROJECT_ID" -eq "" ]]
then 
    DEVSHELL_PROJECT_ID=holy-diver-297719
fi

docker build -t gcr.io/$DEVSHELL_PROJECT_ID/$SERVICE . 
docker push gcr.io/$DEVSHELL_PROJECT_ID/$SERVICE
gcloud alpha run deploy $SERVICE \
    --project $DEVSHELL_PROJECT_ID \
    --set-env-vars PROJECT_ID=$DEVSHELL_PROJECT_ID \
    --set-env-vars ENVIRONMENT=PRODUCTION \
    --image gcr.io/$DEVSHELL_PROJECT_ID/$SERVICE \
    --timeout 5m \
    --no-cpu-throttling \
    --region us-central1 \
    --platform managed \
    --min-instances 0 \
    --max-instances 5 \
    --allow-unauthenticated