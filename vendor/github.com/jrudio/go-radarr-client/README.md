A Go client for Radarr
===

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/jrudio/go-radarr-client)

`go get -u github.com/jrudio/go-radarr-client`



### Cli

You can tinker with this library using the command-line over [here](./cmd)

### Usage

```Go
    client, err := radarr.New("http://192.168.1.12:7878", "abc123")

    results, err := client.Search("Den of Thieves")
```