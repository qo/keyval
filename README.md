## What is it

Implementation of key-value storage described in [Cloud Native Go by Matthew A. Titmus](https://www.oreilly.com/library/view/cloud-native-go/9781492076322/)

## Features

- Concurrency-safe persisting storage implemented with built-in `map` and `sync.RWMutex`

- REST API implemented with [chi](https://github.com/go-chi/chi)

## How to use

### 1. [Install Git](https://git-scm.com/downloads)
### 2. Run `git clone https://github.com/qo/keyval.git`
### 3. [Install Go](https://go.dev/doc/install)
### 4. Run `go run .`
### 5. Access key-value storage via REST API

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

## Todos

### Common
- Create config file to store: API port, log file name, events and errors channels capacity
- Implement SQL data storage
- Implement Transport Layer Security

### File Logger
- Come up with a better solution to process `DELETE` lines in log
- Close log file
- Set maximum size for key and value

## Bugs

### Common
- Service could shut down while there are still events in the events channel

### File logger
- If keys or values contain multiple lines/spaces, bad things might happen (events are parsed with `\t`)
