PROJECT="circa10/k8s-label-rules-webook"
docker build -t "$PROJECT" . && \

docker run --rm -v $(pwd)/sample-rules.yaml:/rules.yaml \
    -p 8080:8080 \
    "$PROJECT" --file rules.yaml