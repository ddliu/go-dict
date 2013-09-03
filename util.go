package dict

import (
    "math/rand"
    "time"
    "io"
    "os"
    "bufio"
)

func randomSlice(slice []string, limit int) []string {
    length := len(slice)
    if limit >= length {
        return slice
    }

    result := []string{}

    found := map[int]bool {}

    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    for {
        k := r.Intn(length)
        if !found[k] {
            found[k] = true
            result = append(result, slice[k])
            limit--
            if limit <= 0 {
                break
            }
        }
    }

    return result
}

func WalkFileLines(file string, f func(line string) bool) {
    fi, err := os.Open(file)
    if err != nil {
        panic(err)
    }

    defer fi.Close()

    r := bufio.NewReader(fi)
    for  {
        line, _, err := r.ReadLine()
        if err == io.EOF {
            break
        }
        if len(line) > 0 {
            if !f(string(line)) {
                break
            }
        }
    }
}