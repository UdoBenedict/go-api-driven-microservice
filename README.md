# API Driven Microservice
Managing sensitive data and configuration securely is a fundamental aspect of operating applications within a Kubernetes cluster. This project tackles that security vulnerability head-on. It features an API-driven microservice written in Go that fetches data from the GitHub REST API. The project demonstrates how to securely inject authentication tokens using Secrets and route configurations via ConfigMaps as environment variables. Kubernetes Secrets and ConfigMaps are crucial tools for this, allowing you to store confidential data and configuration without hardcoding them directly into container images.

## Architectural Diagram
+---------------------------------------------------------+
|                  Kubernetes Cluster                     |
|                                                         |
|  +-------------------+       +-----------------------+  |
|  |  api-secret.yaml  |       |    api-configmap.yaml |  |
|  | (API_TOKEN base64)|       |    (API_BASE_URL)     |  |
|  +-------------------+       +-----------------------+  |
|           |                              |              |
|           |      Environment Variables   |              |
|           +--------------+---------------+              |
|                          |                              |
|                          v                              |
|  +---------------------------------------------------+  |
|  |                 Deployment Pod                    |  |
|  |                                                   |  |
|  |  +---------------------------------------------+  |  |
|  |  |            Docker Container                 |  |  |
|  |  |             (Alpine Linux)                  |  |  |
|  |  |                                             |  |  |
|  |  |   User: appuser (Non-Root)                  |  |  |
|  |  |                                             |  |  |
|  |  |   [ Go Microservice Binary ]                |  |  |
|  |  |   - Reads API_BASE_URL                      |  |  |
|  |  |   - Reads API_TOKEN                         |  |  |
|  |  |   - Serves HTTP on :8080                    |  |  |
|  |  +----------------------+----------------------+  |  |
|  +-------------------------|-------------------------+  |
+----------------------------|----------------------------+
                             |
                   HTTP GET  |  (Injects Bearer Token)
                             v
+---------------------------------------------------------+
|               GitHub REST API (External)                |
+---------------------------------------------------------+

## Tech Stack

`Golang` - `Docker` - `Kubernetes` - `GitHub REST API` 

+ Golang: The native language of Kubernetes, used to build the lightweight, highly performant microservice.
+ Docker: Used for a multi-stage build process to compile the Go code and package the binary into a minimal Alpine Linux image.
+ Kubernetes: The container orchestration platform used to manage deployments, ConfigMaps, and Secrets.
+ GitHub REST API: The easily accessible, authenticated third-party service used as a realistic mock external API for testing.

## What does this project do?
**Key Capabilities**

+ Secure Multi-Stage Containerization: Compiles a statically linked Go binary and runs it within an Alpine operating system to keep the image size small.
+ Non-Root Execution: Implements critical security best practices by creating a dedicated system group and user so the container does not run as the default root user.
+ Dynamic Configuration Management: Reads the external API's base URL dynamically from a Kubernetes ConfigMap, instructing Go to look for environment variables rather than hardcoded URLs.
+ Secure Token Injection: Ingests a GitHub Personal Access Token (PAT) securely from a Kubernetes Secret and attaches it to the external request headers as a Bearer token.

## Core Components
| Component | Tool/Service | Purpose |
| :--- | :---: | ---: |
| Microservice | Golang | Acts as a web server to securely fetch data from an external API using configurations provided by Kubernetes. |
| Containerization | Docker | Packages only the compiled binary into a minimal image and restricts permissions by running as a non-root user. |
| Routing data | K8s Config Map | Safely saves non-sensitive configuration data (the API base URL) that the app can use without hardcoding it. |
| Authentication | K8s Secret | Securely stores the sensitive GitHub API token/key as base64 encoded data. |
| External Data Source | GitHub REST API | Serves as the reliable, authenticated mock third-party API that the microservice interacts with. |

## Deployment Guide
### Step 1:  Write a dockerfile of the application. Build and push the image to Docker Hub.
```
$ docker build -t udonwaigwe/api-service-golang:v1.0 .

$ docker tag udonwaigwe/api-service-golang:v1.0 udonwaigwe/api-service-golang:v1.0

$ docker push udonwaigwe/api-service-golang:v1.0

#The push refers to repository [docker.io/udonwaigwe/api-service-golang]
```

## Step 2: Docker Login. Start Minikube.


```
$ docker login -u <username>

$ start minikube

```

## Step 3: Create a Namespace and Deploy Resources
```
$ kubectl create namespace demo
$  kubectl -n demo apply -f api-configmap.yaml
$ kubectl -n demo apply -f deployment.yaml
```

### Step 4: Port forward for External Access
Wait for your pod to start running, then:

```
kubectl -n demo port-forward deployment/go-api-service 8080:8080
```

### Step 5: Testing the Application
On a new terminal, run:

```
curl http://localhost:8080/fetch-data
```

The Go app reads the encrypted token from the mounted environment variable, attaches it to the request header, and returns my raw GitHub profile data in JSON format.

## Step 6: Cleanup
```
# stop the port forwarding
Ctrl + C

# delete the namespace
$ kubectl delete namespace demo

# stop minikube
$ minikube stop
```



