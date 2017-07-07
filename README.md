# Go Liquid Template Parser

 [![][travis-svg]][travis-url] [![][coveralls-svg]][coveralls-url] [![][go-report-card-svg]][go-report-card-url] [![][godoc-svg]][godoc-url] [![][license-svg]][license-url]

> “Any sufficiently complicated C or Fortran program contains an ad-hoc, informally-specified, bug-ridden, slow implementation of half of Common Lisp.” – Philip Greenspun

`liquid` ports [Shopify Liquid templates](https://shopify.github.io/liquid) to Go. It was developed for use in [gojekyll](https://github.com/osteele/gojekyll).

`liquid` provides a functional API for defining tags and filters. See examples [here](https://github.com/osteele/liquid/blob/master/filters/filters.go), [here](https://github.com/osteele/gojekyll/blob/master/filters/filters.go), and [here](https://github.com/osteele/gojekyll/blob/master/tags/tags.go).

<!-- TOC -->

- [Go Liquid Template Parser](#go-liquid-template-parser)
    - [Status](#status)
    - [Other Differences](#other-differences)
    - [Install](#install)
    - [Contributing](#contributing)
    - [References](#references)
    - [Attribution](#attribution)
    - [Other Implementations](#other-implementations)
        - [Go](#go)
        - [Other Languages](#other-languages)
    - [License](#license)

<!-- /TOC -->

## Status

This library is in early development. There's probably lots of corner cases, and the API (especially around defining tags) may still change.

- [ ] Basics
  - [x] Literals
    - [ ] String Escapes
  - [x] Variables
  - [x] Operators
  - [x] Arrays
  - [ ] Whitespace Control
- [ ] Tags
  - [x] Comment
  - [ ] Control Flow
    - [x] `if`/`else`/`elsif`
    - [x] `unless`
    - [x] `case`
      - [x] `when`
      - [ ] `when a or b`
      - [ ] `else`
  - [ ] Iteration
      - [x] modifiers (`limit`, `reversed`, `offset`)
      - [ ] range
      - [x] `break`, `continue`
      - [x] loop variables
      - [ ] `tablerow`
      - [ ] `cycle`
  - [x] Include
  - [x] Raw
  - [x] Variables
    - [x] Assign
    - [x] Capture
- [ ] Filters
  - [ ] `sort_natural`, `uniq`, `escape`, `truncatewords`, `url_decode`, `url_encode`
  - [x] everything else
- [x] Drops

## Other Differences

These will change:

* I haven't investigated the interaction of loop reversed, limit, and offset, and whether it matters which order they're specified, in Shopify Liquid.
* This implementation doesn't parse very many date formats. I have any idea for this; talk to me if you're interested in implementing it.
* Identifiers can include hyphens.
* `contains` isn't implemented on hashes.
* `contains` doesn't stringify its argument.
* The `nil` constant isn't implemented.

These might not:

* This parser accepts parentheses in places where Shopify Liquid does not.
* `else` and `elsif` work inside of `unless`

## Install

`go get -u github.com/osteele/goliquid`

`make install` install a command-line `liquid` program in your GO bin.
This is intended to make it easier to create test cases for bug reports.
Run `liquid --help` for help.

## Contributing

Bug reports, test cases, and code contributions are more than welcome.
Please refer to the [contribution guidelines](./CONTRIBUTING.md).

## References

* [Shopify.github.io/liquid](https://shopify.github.io/liquid)
* [Liquid for Designers](https://github.com/Shopify/liquid/wiki/Liquid-for-Designers)
* [Liquid for Programmers](https://github.com/Shopify/liquid/wiki/Liquid-for-Programmers)
* [Help.shopify.com](https://help.shopify.com/themes/liquid) goes into more detail, but includes features that aren't present in core Liquid as used by Jekyll.

## Attribution

| Package                                               | Author          | Description                                  | License            |
|-------------------------------------------------------|-----------------|----------------------------------------------|--------------------|
| [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml)   | Canonical       | YAML support (for printing parse trees)      | Apache License 2.0 |
| [jeffjen/datefmt](https://github.com/jeffjen/datefmt) | Jeffrey Jen     | Go bindings to GNU `strftime` and `strptime` | MIT                |
| [Ragel](http://www.colm.net/open-source/ragel/)       | Adrian Thurston | scanning expressions                         | MIT                |

Michael Hamrah's [Lexing with Ragel and Parsing with Yacc using Go](https://medium.com/@mhamrah/lexing-with-ragel-and-parsing-with-yacc-using-go-81e50475f88f) was essential to understanding `go yacc`.

The [original Liquid engine](https://shopify.github.io/liquid), of course, for the design and documentation of the Liquid template language. Many of the tag and filter test cases are taken directly from the Liquid documentation.

## Other Implementations

### Go

* [karlseguin/liquid](https://github.com/karlseguin/liquid) is a dormant implementation that inspired a lot of forks.
* [acstech/liquid](https://github.com/acstech/liquid) is a more active fork of Karl Seguin's implementation.
* [hownowstephen/go-liquid](https://github.com/hownowstephen/go-liquid)

### Other Languages

 See Shopify's [ports of Liquid to other environments](https://github.com/Shopify/liquid/wiki/Ports-of-Liquid-to-other-environments).

## License

MIT License

[coveralls-url]: https://coveralls.io/r/osteele/liquid?branch=master
[coveralls-svg]: https://img.shields.io/coveralls/osteele/liquid.svg?branch=master

[godoc-url]: https://godoc.org/github.com/osteele/liquid
[godoc-svg]: https://godoc.org/github.com/osteele/liquid?status.svg

[license-url]: https://github.com/osteele/liquid/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg

[go-report-card-url]: https://goreportcard.com/report/github.com/osteele/liquid
[go-report-card-svg]: https://goreportcard.com/badge/github.com/osteele/liquid

[travis-url]: https://travis-ci.org/osteele/liquid
[travis-svg]: https://img.shields.io/travis/osteele/liquid.svg?branch=master
