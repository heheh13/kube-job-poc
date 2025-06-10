# File: server/Makefile

APP_NAME=job-server
AGENT_APP_NAME=agent-job
IMAGE_TAG=latest
IMAGE_REPO=heheh13/$(APP_NAME)

# Build the Go binary
build:
	go build -o server main.go

# Build the Docker image
docker-build:
	docker build -t $(IMAGE_REPO):$(IMAGE_TAG) .

# Push the Docker image
docker-push:
	docker push $(IMAGE_REPO):$(IMAGE_TAG)


# Build and push Docker image
build-push-server:
	docker build -t $(IMAGE_REPO):$(IMAGE_TAG) .
	docker push $(IMAGE_REPO):$(IMAGE_TAG)

# Build and push agent Docker image
build-push-agent:
	cd ./agents && docker build -t $(IMAGE_REPO):$(IMAGE_TAG) .
	docker push $(IMAGE_REPO):$(IMAGE_TAG)

# Apply Kubernetes manifests (create server Deployment & Service)
kube-deploy:
	kubectl apply -f manifests/RBAC.yaml
	kubectl apply -f manifests/deployment.yaml

# Delete Kubernetes resources
kube-clean:
	kubectl delete -f manifests/deployment.yaml

# Full cycle: build, push, and deploy
deploy: build-push-server build-push-agent kube-deploy

.PHONY: build docker-build docker-push kube-deploy kube-clean deploy
