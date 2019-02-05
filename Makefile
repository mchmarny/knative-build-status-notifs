
.PHONY: event


secrets:
	kubectl create secret generic slacknotifs-secrets -n demo \
		--from-literal=APP_TOKEN=$(PUSHOVER_APP_TOKEN) \
		--from-literal=USR_TOKEN=$(PUSHOVER_USR_TOKEN)


deps:
	go mod tidy

run:
	go run cmd/service/*.go --sink=https://events.demo.knative.tech/

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/slacknotifs:latest

docker:
	docker build -t slacknotifs .

source:
	kubectl apply -f config/source.yaml

service:
	kubectl apply -f config/service.yaml


cleanup:
	kubectl delete -f config/source.yaml
	kubectl delete -f config/service.yaml
