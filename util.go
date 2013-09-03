package dict

import (
    "math/rand"
    "time"
    "io"
    "os"
    "bufio"
    "strconv"
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

func ToInt64(v interface{}) (int64, bool) {
    switch v := v.(type) {
    case int:
        return int64(v), true
    case int32:
        return int64(v), true
    case int64:
        return v, true
    case float32:
        return int64(v), true
    case float64:
        return int64(v), true
    case string:
        result, err := strconv.ParseInt(v, 10, 64)
        if err != nil {
            return 0, false
        }
        return result, true
    }
    return 0, false
}

func ToInt32(v interface{}) (int32, bool) {
    switch v := v.(type) {
    case int:
        return int32(v), true
    case int32:
        return v, true
    case int64:
        return int32(v), true
    case float32:
        return int32(v), true
    case float64:
        return int32(v), true
    case string:
        result, err := strconv.ParseInt(v, 10, 32)
        if err != nil {
            return 0, false
        }
        return int32(result), true
    }

    return 0, false
}

func ToInt(v interface{}) (int, bool) {
    switch v := v.(type) {
    case int:
        return v, true
    case int32:
        return int(v), true
    case int64:
        return int(v), true
    case float32:
        return int(v), true
    case float64:
        return int(v), true
    case string:
        result, err := strconv.ParseInt(v, 10, 0)
        if err != nil {
            return 0, false
        }
        return int(result), true
    }

    return 0, false
}