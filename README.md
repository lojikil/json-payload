# Overview

`json-payload` is a simple CLI program, written in Go, that is meant to help with the process of creating JSON payloads.


    $ ./json-payload test:this test1
    26
    {"test":"this","test1":""}
    52
    %7B%22test%22%3A%22this%22%2C%22test1%22%3A%22%22%7D

The output is relatively simple:

- the first line is the length of the payload
- the second is the un-encoded payload
- the third is the length of the encoded payload
- the fourth is the percent-encoded payload

# Ideas

- arrays: `foo:[1,2,3,4,5]`, or `foo:[1,2,3,4,5]:int`
- `multi.level.objects=test`
- JSON => CSRF payloads

# License

see `LICENSE` for details (ISC license).
