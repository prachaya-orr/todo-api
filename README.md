```go build\
    -ldflags "-X main.buildcommit=`git rev-parse --short Head` \
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
    -o app