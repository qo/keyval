## What is it

Implementation of key-value storage described in [Cloud Native Go by Matthew A. Titmus](https://www.oreilly.com/library/view/cloud-native-go/9781492076322/)

## Features

- Concurrency-safe storage implemented with built-in `map` and `sync.RWMutex`

- REST API implemented with [chi](https://github.com/go-chi/chi)

## How to use

### 1. [Install Go](https://go.dev/doc/install)
### 2. Run `go run .`
### 3. Access key-value storage via REST API

#### Route `localhost:8090/v1/key/{key}`:

##### `PUT`
Put value for provided `key` (don't forget to pass data)

##### `GET`
Get value for provided `key`

##### `DELETE`
Delete value for provided `key`

### Examples with [curl](https://curl.se/docs/manpage.html):

`curl -X PUT -d "1" http://localhost:8090/v1/key/a` - put value `1` for `a` key

`curl -X GET http://localhost:8090/v1/key/a` - get value for `a` key

`curl -X DELETE http://localhost:8090/v1/key/a` - delete value for `a` key

### Notes

- Empty keys/values are not supported
- Make sure nothing is running on `localhost:8090`

## Todo

- Persisting Resource State
- Implementing Transport Layer Security
