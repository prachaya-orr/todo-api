```go build\
    -ldflags "-X main.buildcommit=`git rev-parse --short Head` \
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
    -o app

``` echo "GET http://localhost:8081/limitz" | vegeta attack -rate=100 -duration=1s | vegeta report