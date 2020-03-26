# [k8s-label-rules-webhook](https://github.com/circa10a/k8s-label-rules-webhook)

Enforce standards for labels of resources being created in your k8s cluster

![image](https://i.imgur.com/yiqw1we.png)

[![Build Status](https://travis-ci.org/circa10a/k8s-label-rules-webhook.svg?branch=master)](https://travis-ci.org/circa10a/k8s-label-rules-webhook)
[![Go Report Card](https://goreportcard.com/badge/github.com/circa10a/k8s-label-rules-webhook)](https://goreportcard.com/report/github.com/circa10a/k8s-label-rules-webhook)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/circa10a/k8s-label-rules-webhook?style=plastic)
![Docker Pulls](https://img.shields.io/docker/pulls/circa10a/k8s-label-rules-webhook?style=plastic)

## Table of Contents

* [Usage](#usage)
  + [Docker](#docker)
    - [Volume mount your `rules.yaml` file](#volume-mount-your-rulesyaml-file)
    - [Build your own docker image](#build-your-own-docker-image)
  + [Kubernetes](#kubernetes)
    - [Deploy webhook application](#deploy-webhook-application)
    - [Deploy webhook application (TLS)](#deploy-webhook-application-tls)
    - [Deploy admission webhook](#deploy-admission-webhook)
* [Features](#features)
  + [Hot reloading of ruleset](#hot-reloading-of-ruleset)
  + [Rule validation](#rule-validation)
  + [Easily view loaded ruleset](#easily-view-loaded-ruleset)
  + [Prometheus Metrics](#prometheus-metrics)
* [Configuration](#configuration)
* [Development](#development)
  + [Build](#build)
  + [Run](#run)
  + [Test](#test)
* [Changelog](#changelog)
* [Troubleshooting](#troubleshooting)

## Usage

![gif](https://j.gifs.com/6X7Vvn.gif)

Start by creating a `rules.yaml` file containing rules for labels you require for your cluster resources to have along with a regex pattern for the values of the labels.

> Any rules specified in the rulseset will be required on resources to which you configure the admission webhook to fire on. View the kubernetes deployment section.

```yaml
rules:
  - name: require-phone-number
    key: phone-number
    value:
      regex: "[0-9]{3}-[0-9]{3}-[0-9]{4}" # 555-555-5555
  - name: require-owner
    key: owner
    value:
      regex: ".*" # Any pattern matches, Just ensure a label of "owner" is set
```

> IMPORTANT NOTE: Invalid regex for a given rule will make the rule default to .* which will allow any label value, but will still require the label to be present

Once you have your ruleset, you can deploy the webhook several different ways.

### Docker

#### Volume mount your `rules.yaml` file

```shell
docker run -d --name k8s-label-rules-webhook \
  -p 8080: 8080 \
  -v rules.yaml:/rules.yaml \
  circa10a/k8s-label-rules-webhook
```

#### Build your own docker image

```dockerfile
FROM circa10a/k8s-label-rules-webhook
COPY rules.yaml /
```

### Kubernetes

####  Deploy webhook application

> Kubernetes admission webhooks require https

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: label-rules-webhook
  labels:
    app: label-rules-webhook
spec:
  replicas: 1
  selector:
    matchLabels:
      app: label-rules-webhook
  template:
    metadata:
      labels:
        app: label-rules-webhook
    spec:
      containers:
      - name: label-rules-webhook
        image: circa10a/k8s-label-rules-webhook
        volumeMounts:
        - name: label-rules
          mountPath: /rules.yaml
          subPath: rules.yaml
        readinessProbe:
          httpGet:
            path: /rules
            port: gin-port
          initialDelaySeconds: 5
          periodSeconds: 15
        livenessProbe:
          httpGet:
            path: /rules
            port: gin-port
          initialDelaySeconds: 5
          periodSeconds: 15
        ports:
        - name: gin-port
          containerPort: 8080
      volumes:
        - name: label-rules
          configMap:
            name: label-rules
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: label-rules
data:
  rules.yaml: |
    rules:
      - name: require-phone-number
        key: phone-number
        value:
          regex: "[0-9]{3}-[0-9]{3}-[0-9]{4}"
      - name: require-number
        key: number
        value:
          regex: "[0-1]{1}"
---
apiVersion: v1
kind: Service
metadata:
  name: label-rules-webhook-service
spec:
  selector:
    app: label-rules-webhook
  ports:
  - protocol: TCP
    port: 8080
    targetPort: gin-port
  type: NodePort
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: label-rules-webhook-ingress
spec:
  backend:
    serviceName: label-rules-webhook-service
    servicePort: 8080
```

#### Deploy webhook application (TLS)

To have the web server listen on https, you need to supply a certificate and a key in conjunction with the appropriate [configuration options](#configuration).

Here's an [example deployment](examples/tls.yaml) which supplies a cert, key and enables TLS in the application.

> Note: The default https port is `8443`

#### Deploy admission webhook

[More info on kubernetes admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)

```yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: label-rules-webhook
webhooks:
- name: my.application.domain
  clientConfig:
    url: "https://<deployed-app-location>/"
  rules:
  - operations:
    - "CREATE"
    - "UPDATE"
    apiGroups:
      - "apps"
    apiVersions:
      - "v1"
      - "v1beta1"
    resources:
      - "deployments"
      - "replicasets"
  failurePolicy: Fail # Ignore, Fail
```

## Features

Checkout the swagger api docs at `/swagger/index.html`

### Hot reloading of ruleset

Update the `rules.yaml` file used by your deployed instance then send a `POST` request to `/reload` to reload the rules into memory without downtime.

### Rule validation

The regex supplied to each rule is compiled when the application starts and any problems are logged.

You can access the `/validate` endpoint via `GET` request to view any issues with the current ruleset that is loaded.

### Easily view loaded ruleset

Access the `/rules` endpoint via `GET` request to see the current rules loaded.

### Prometheus Metrics

Prometheus metrics are enabled by default and are available at the `/metrics` endpoint. Simply unset the `METRICS` environment variable to disable.

## Configuration

|             |                                                                       |                      |                        |           |               |
|-------------|-----------------------------------------------------------------------|----------------------|------------------------|-----------|---------------|
| Name        | Description                                                           | Environment Variable | Command Line Argument  | Required | Default        |
| GIN MODE    | Runs web server in production or debug mode                           | `GIN_MODE`           | NONE                   | `false`  | `release`      |
| PORT        | Port for web server to listen on                                      | `PORT`               | NONE                   | `false`  | `8080`         |
| METRICS     | Enables prometheus metrics on `/metrics`(unset for false)             | `METRICS`            | `--metrics`            | `false`  | `true`         |
| RULES       | File containing user defined ruleset(default looks to `./rules.yaml`) | NONE                 | `--file`               | `true`   | `./rules.yaml` |
| TLS         | Start web server listening on HTTPS                                   | `TLS_ENABLED`        | `--tls`                | `false`  | `false`        |
| TLS CERT    | TLS Certificate file path                                             | `TLS_CERT`           | `--tls-cert`           | `false`  | None           |
| TLS KEY     | TLS key file path                                                     | `TLS_KEY`            | `--tls-key`            | `false ` | None           |
| TLS PORT    | TLS listening port                                                    | None                 | `--tls-port`           | `false`  | `8443`         |

## Development

### Build

```shell
make build
```

### Run

```shell
make run
```

Access via http://localhost:8080

### Test

```shell
make test
```

## Changelog

[Changelog.md](CHANGELOG.md)

## Troubleshooting

### Debug web server

Try setting the debug environment variable for the gin web server which is `GIN_MODE=debug`.

### Rule validation problems

Ensure there are no regex compilation errors by accessing the `/validate` endpoint.
