```
{
"message":" {\"chipUserId\": \"61e079edc893200e127f5401\"}"
}
```

rm queue.zip; GOARCH=amd64 GOOS=linux go build ./cmd/queue; zip queue.zip queue; rm queue
