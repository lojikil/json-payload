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
 * - a proper parser, so that we can escape ':' in strings
 * - array support, including typed arrays: `test:[1, 2, 3, 4, 5]:int`
 * - `multi.level.objects=foo`
 */

func parse(src string) (string, interface{}, bool) {
    /* we return the:
     * - name
     * - value
     * as a properly parsed structure
     */
    state := Start
    offset := 0
    svalue := make([]byte, 256)
    validx := 0
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
                    name = string(svalue[0:validx])
                    validx = -1
                    if idx + 1 >= len(src) {
                        state = End
                    } else {
                        state = ValTransfer
                    }
                } else {
                    svalue[validx] = src[idx]
                }
                validx++
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
                    retval = string(svalue[0:validx])
                    if idx + 1 >= len(src) {
                        state = End
                    } else {
                        state = TypeTransfer
                    }
                } else {
                    svalue[validx] = src[idx]
                }
                validx++
            case TypeTransfer:
                if idx + 1 >= len(src) {
                    state = End
                } else if src[idx] == ':' {
                    state = TypeDec
                } else {
                    state = Error
                }
            case Error:
                return name, nil, true
        }
    }
    return name, retval, !(state == End)
}

func main() {
    res :=  make(map[string] interface{})
    args := os.Args[1:]

    for i := 0; i < len(args); i++ {
        /*vals := strings.Split(args[i], ":")
        l := len(vals);
        if l == 1 {
            res[args[i]] = ""
        } else if l == 2 {
            res[vals[0]] = vals[1]
        } else if l == 3 {
            switch vals[2] {
                case "int":
                    j, err := strconv.Atoi(vals[1])

                    if err != nil {
                        break
                    }

                    res[vals[0]] = j
                case "float":
                    f, err := strconv.ParseFloat(vals[1], 64)

                    if err != nil {
                        break
                    }

                    res[vals[0]] = f
                default:
                    res[vals[0]] = vals[1];
            }
            fmt.Println("name with value & type")
        } else {
            fmt.Println("error");
        }*/
        name, value, err := parse(args[i])
        if !err {
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
