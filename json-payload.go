package main

import (
    "fmt"
    "encoding/json"
    "net/url"
    "os"
    //"strings"
    //"strconv"
)

const (
    Start int = iota
    Name // a raw name like 'test'
    NameStr // a name inside string quotes '"test"'
    Value // a raw value like 'test'
    ValArr // an array of values, not nested
    ValStr // a quoted string value
    ValTransfer
    TypeTransfer
    TypeDec
    End
    Error
)


/* TODO:
 * - DONE: a proper parser, so that we can escape ':' in strings
 * - array support, including typed arrays: `test:[1, 2, 3, 4, 5]:int`
 * - `multi.level.objects=foo` need to figure out merging these too...
 */

func parse_quotedstring(src string, offset int) (string, int, bool) {
    svalue := make([]byte, 256)
    state := Start
    validx := 0
    for idx := offset; idx < len(src); idx, validx = idx + 1, validx + 1 {
        if src[idx] == '\\' {
            if src[idx + 1] == 'n' {
                svalue[validx] = '\n'
            } else if src[idx + 1] == 'r' {
                svalue[validx] = '\r'
            } else if src[idx + 1] == 't' {
                svalue[validx] = '\t'
            } else if src[idx + 1] == '\\' {
                svalue[validx] = '"'
            }
            idx++
        } else if src[idx] == '"' {
            state = End
        } else {
            svalue[validx] = src[idx]
        }

        if state == End {
            offset = idx
            break
        }
    }
    return string(svalue[0:validx]), offset, !(state == End)
}

func parse(src string) (string, interface{}, int) {
    /* we return the:
     * - name
     * - value
     * as a properly parsed structure
     */
    state := Start
    offset := 0
    tmperror := false
    var name string
    var retval interface{}
    for idx := 0; idx < len(src); idx++ {
        switch state {
            case Start:
                if  src[idx] == '"' {
                    state = NameStr
                } else {
                    state = Name
                }
            case Name:
                if src[idx] == ':' {
                    state = ValTransfer
                    name = src[0:idx]
                } else if idx + 1 >= len(src) {
                    retval = ""
                    name = src[0:idx+1]
                    state = End
                }
            case NameStr:
                /* can probably collapse the two bare-string and quoted-string
                 * states into one with a little book keeping.
                 */
                name, idx, tmperror = parse_quotedstring(src, idx)
                state = ValTransfer
            case Value:
                retval = string(src[offset:idx + 1])
                if src[idx] == ':' {
                    state = TypeTransfer
                } else if idx + 1 >= len(src) {
                    state = End
                }
            case ValTransfer:
                if src[idx] == '"' {
                    state = ValStr
                } else if src[idx] == '[' {
                    state = ValArr
                } else if src[idx] == ':' {
                    state = ValTransfer
                } else {
                    offset = idx
                    state = Value
                }
            case ValStr:
                retval, idx, tmperror = parse_quotedstring(src, idx)
                if tmperror {
                    fmt.Println("error on line 118")
                    state = Error
                } else {
                    state = TypeTransfer
                }
            case TypeTransfer:
                if idx + 1 >= len(src) {
                    state = End
                } else if src[idx] == ':' {
                    state = TypeDec
                } else {
                    state = Error
                }
            case Error:
                return name, nil, state
        }
    }

    if state == TypeTransfer {
        state = End
    }
    return name, retval, state
}

func main() {
    res :=  make(map[string] interface{})
    args := os.Args[1:]

    for i := 0; i < len(args); i++ {
        name, value, err := parse(args[i])

        if err == End {
            res[name] = value
        } 
    }

    b, err := json.Marshal(res)

    if err != nil {
        fmt.Println("json encoding error!")
    } else {
        output := string(b)
        eoutput := url.QueryEscape(output)
        fmt.Printf("%d\n%s\n%d\n%s\n", len(output), output, len(eoutput), eoutput)
    }
}
