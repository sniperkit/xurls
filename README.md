# xurls

[![GoDoc](https://godoc.org/github.com/sniperkit/xurls?status.svg)](https://godoc.org/github.com/sniperkit/xurls)
[![Travis](https://travis-ci.org/mvdan/xurls.svg?branch=master)](https://travis-ci.org/mvdan/xurls)

Extract urls from text using regular expressions.

	go get -u github.com/sniperkit/xurls/...

```go
import "github.com/sniperkit/xurls/pkg"

func main() {
	xurls.Relaxed().FindString("Do gophers live in golang.org?")
	// "golang.org"
	xurls.Strict().FindAllString("foo.com is http://foo.com/.", -1)
	// []string{"http://foo.com/"}
}
```

Note that the funcs compile regexes, so avoid calling them repeatedly.

#### cmd/xurls

	go get -u github.com/sniperkit/xurls/cmd/xurls

```shell
$ echo "Do gophers live in http://golang.org?" | xurls
http://golang.org
```
