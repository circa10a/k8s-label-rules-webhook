FROM golang:alpine
WORKDIR /go/src/app
COPY . .
ENV USER=go \
    UID=1000 \
    GID=1000 \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

RUN go build -ldflags="-s -w -X main.Version=$(grep -o '[0-9]\{1,\}\.[0-9]\{1,\}\.[0-9]\{1,\}' main.go)" \
    -o webhook && \
    addgroup --gid "$GID" "$USER" && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER" && \
    chown "$UID":"$GID" /go/src/app/webhook

FROM scratch
ENV GIN_MODE=release \
    METRICS=true
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /go/src/app/webhook /
USER 1000
ENTRYPOINT ["/webhook"]