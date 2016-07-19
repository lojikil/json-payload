package main

import (
    "fmt"
    "encoding/json"
    "net/url"
    "os"
    "strings"
    "strconv"
)

/* TODO:
 * - a proper parser, so that we can escape ':' in strings
 * - array support, including typed arrays: `test:[1, 2, 3, 4, 5]:int`
 * - `multi.level.objects=foo`
 */

func main() {
    res :=  make(map[string] interface{})
    args := os.Args[1:]

    for i := 0; i < len(args); i++ {
        vals := strings.Split(args[i], ":")
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
