# gatherer_go

A rewrite of [gatherer_cli](https://github.com/Dorthu/gatherer_cli) in golang

### Why?

This is my first project using Go, I wanted to do something easy.

### How to Use

```bash
go get github.com/dorthu/gatherer_go
go install github.com/dorthu/gatherer_go
gatherer_go --help
```

This assumes `~/go/bin` is on your `PATH`.

### Incomplete Features

Right now without a `-b` this only outputs the search URL - in the future this
will query the page itself and scrape the results into something you can see in
your terminal (probably piped to less).
