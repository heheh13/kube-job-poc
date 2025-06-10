# Kube-job-poc
A Go-based HTTP server running in a Kubernetes pod receives requests to dynamically create and launch jobs using
a different Docker image.

# Build & Push Images

`make build-push-server` build the server docker images

`make build-push-agent` build the server docker images

# Deploy to Kubernetes

`make deploy` will build docker images and create rbac and deployments

# Usage
**Port-forward the service:**

`kubectl port-forward svc/job-server 8080:80`

**Call the /start-job endpoint:**

`curl "http://localhost:8080/start-job"`

# Custom command

`curl "http://localhost:8080/start-job?cmd=echo%20hello%20world%20%3B%20sleep%2010"`


# improvements
- cluster wide RBAC for create job in different namespace
- usage of base64 encoded command for simply maintaining custom command