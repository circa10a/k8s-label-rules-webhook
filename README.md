# k8s-label-rules-webhook

Enforce standards for labels of resources being created in your k8s cluster

![mascot](https://i.imgur.com/4XDDuJK.png)

## Usage

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

```shell
test
```

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

Update the `rules.yaml` file used by you're deployed instance then send a `POST` request to `/reload` to reload the rules into memory without downtime.

### Rule validation

The regex supplied to each rule is compiled when the application starts and is then logged to indicate problems.

You can access the `/validate` endpoint via `GET` request to view any issues with the current ruleset that is loaded.

### Easily view loaded ruleset

Access the `rules` endpoint via `GET` request to see the current rules loaded.

### Prometheus Metrics

Prometheus metrics are enabled by default and are available at the `/metrics` endpoint. Simply unset the `METRICS` environment variable to disable.

## Configuration

|             |                                                                       |                      |                        |           |               |
|-------------|-----------------------------------------------------------------------|----------------------|------------------------|-----------|---------------|
| Name        | Description                                                           | Environment Variable | Command Line Argument  | Required | Default        |
| GIN MODE    | Runs web server in production or debug mode                           |`GIN_MODE`            | NONE                   | `false`  | `release`      |
| PORT        | Port for web server to listen on                                      | `PORT`               | NONE                   | `false`  | `8080`         |
| METRICS     | Enables prometheus metrics on `/metrics`(unset for false)             |`METRICS`             | `--metrics`            | `false`  | `true`         |
| SWAGGER-URL | Location of swagger `doc.json`                                        | NONE                 | `--swagger-url`        | `false`   | `http://localhost${PORT}/swagger/doc.json`|
| RULES       | File containing user defined ruleset(default looks to `./rules.yaml`) | NONE                 | `--file`               | `true`   | `./rules.yaml` |

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
