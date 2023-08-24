# Yeti

is a helper for maplestory

## Deps

```
go run main.go db init --uri "yeti.db"
go run main.go db migrate --uri "yeti.db"
```

## Features

- IED calculator

## Test

```
curl -X POST -H 'Content-Type: application/json' -d '{"group_id": 0, "message_id": 0, "message":"/mg $IGN", "sender": {"nickname":"nick", "user_id": 0}}' "localhost:25700"  
```