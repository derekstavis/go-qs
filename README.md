# go-qs

>A Go port of Rack's query string parser.

This package was written as I haven't found a good package that understands
[Rack/Rails](http://guides.rubyonrails.org/form_helpers.html#understanding-parameter-naming-conventions) query string [format](https://gist.github.com/dapplebeforedawn/3724090).

The good thing about this package is that it parses nested query strings into
`map[string]interface{}`, the same format as Go `json.Unmarshal`.

## Compatibility

`go-qs` is a literal port of [Rack's code](https://github.com/rack/rack/blob/rack-1.3/lib/rack/utils.rb#L114),
and the test suite [is also a port](https://github.com/derekstavis/go-qs/blob/master/qs/qs_test.go)
of [Rack tests](https://github.com/rack/rack/blob/rack-1.3/test/spec_utils.rb#L107).

## Usage

There's only one function available:

```go
package main

import "github.com/derekstavis/go-qs"

query, err := qs.Parse("foo=bar&names[]=foo&names[]=bar")

if err != nil {
  fmt.Printf("%#+v\n", query)
}
```

Should output this:

```
map[string]interface {}{"foo":"bar", "names":[]interface {}{"foo", "bar"}}
```

## License

MIT (c) 2016 Derek W. Stavis

