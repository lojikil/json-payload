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

A more complex example is:

    $ ./json-payload 'test:"test this thing"' 'foo:"bar baz blah"' foo2 '"this is a key test":"test value"' foo3:bar3
    105
    {"foo":"bar baz blah","foo2":"","foo3":"bar3","test":"test this thing","this is a key test":"test value"}
    167
    %7B%22foo%22%3A%22bar+baz+blah%22%2C%22foo2%22%3A%22%22%2C%22foo3%22%3A%22bar3%22%2C%22test%22%3A%22test+this+thing%22%2C%22this+is+a+key+test%22%3A%22test+value%22%7D

# Ideas

- arrays: `foo:[1,2,3,4,5]`, or `foo:[1,2,3,4,5]:int`
- `multi.level.objects=test`
- JSON => CSRF payloads

# License

see `LICENSE` for details (ISC license).
