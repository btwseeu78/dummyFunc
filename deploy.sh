gcloud functions deploy go-pubsub-function \
--gen2 \
--runtime=go121 \
--region=us-central1 \
--entry-point=HelloPubSub \
--trigger-topic=gke-notify \
--set-env-vars DYNA_URL=https://www.google.com/api/v2,GKE_PROJECT_ID=irn-70740-dev,DT_API_TOKEN=dt.abc1w22233 \
--service-account=gkeupdatetriger@symmetric-ion-423517-n3.iam.gserviceaccount.com \
--trigger-service-account=gkeupdatetriger@symmetric-ion-423517-n3.iam.gserviceaccount.com
--docker-registry=us-central1-docker.pkg.dev/symmetric-ion-423517-n3/gcf-artifacts/go--pubsub--function:latest
gcloud functions add-invoker-policy-binding go-pubsub-function \
      --region="us-central1" \
      --member="gkeupdatetriger@symmetric-ion-423517-n3.iam.gserviceaccount.com"

      --source=. \