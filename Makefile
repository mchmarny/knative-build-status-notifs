
.PHONY: event mod


mod:
	go mod tidy
	go mod vendor


image: mod
	gcloud builds submit \
		--project cloudylabs-public \
		--tag gcr.io/cloudylabs-public/build-notif:0.4.1

source:
	kubectl apply -f config/source.yaml

service:
	kubectl apply -f config/service.yaml


cleanup:
	kubectl delete -f config/source.yaml
	kubectl delete -f config/service.yaml
