
.PHONY: event


image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/build-notif

source:
	kubectl apply -f config/source.yaml

service:
	kubectl apply -f config/service.yaml


cleanup:
	kubectl delete -f config/source.yaml
	kubectl delete -f config/service.yaml
