# Keisatsu

> *"Keisatsu desu"*

*Keisatsu* (警察) is a simple monitoring Go package which can help to notify you
by sending a webhook to the registered URL whenever there is a panic in your Go application 
or when you told it to send you a notification.  

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` will always pull the latest tagged release from the master branch.
```
go get github.com/hansenedrickh/keisatsu
``` 

### Usage

First, you need to initialize the Keisatsu by providing it the necessary information 
 
```go
k := keisatsu.New(appName, webhookURL, secretToken)
``` 

Then, you can put this script to handle a panic situation wherever you want

```go
defer k.WatchPanic()
```

Or if you want to handle expected error, then you can use this script

```go
k.Error(message)
```

It will send a `POST` http request to the given webhook url, with this json as a body
```json
{
  "app_name": <your_app_name>,
  "level": <error_level>,  
  "message": <error_message>,
  "stack_trace": <stack_trace>
}
```

The `secret token` will be put in the `X-Secret-Token` header which you can use for validating the request coming from Keisatsu.

## Credit

Made with :heart: by [Hansen Edrick Harianto](https://github.com/hansenedrickh)