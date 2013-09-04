package dict

import (
    "os"
    "regexp"
    "github.com/ddliu/dict/util"
)

func NewSimpleDict() *SimpleDict {
    return &SimpleDict{}
}

type SimpleDict struct {
    words []string
}

func (this *SimpleDict) Load(dict string) {
    util.WalkFileLines(dict, func(line string) bool {
        this.words = append(this.words, string(line[:]))
        return true
    })
}

func (this *SimpleDict) Export(dict string) {
    fi, err := os.OpenFile(dict, os.O_CREATE|os.O_WRONLY, 0660)

    if err != nil {
        panic(err)
    }

    defer fi.Close()

    for i, length := 0, len(this.words); i < length; i++ {
        _, err := fi.WriteString(this.words[i] + "\n")
        if err != nil {
            panic(err)
        }
    }
}

func (this *SimpleDict) AddWords(words ...string) {
    this.AddWordsList(words)
}

func (this *SimpleDict) AddWordsList(words []string) {
    for _, w := range words {
        this.words = append(this.words, w)
    }
}

func (this *SimpleDict) Clear() {
    this.words = []string{}
}

func (this *SimpleDict) Count() int {
    return len(this.words)
}

func (this *SimpleDict) LookupAll(pattern string) []string {
    return this.Lookup(pattern, 0, 0)
}

func (this *SimpleDict) Lookup(pattern string, offset int, limit int) []string {
    if pattern == "" {
        if offset == 0 && limit == 0 {
            return this.words[:]
        }
        if offset >= len(this.words) {
            return nil
        }

        var end int
        if limit <= 0 {
            end = len(this.words)
        } else {
            end = offset + limit
        }

        return this.words[offset:end]
    }

    var compiledPattern = regexp.MustCompile("^" + pattern + "$")
    var matched []string
    found := 0
    this.Walk(func(word string) bool {
        if compiledPattern.MatchString(word) {
            found++

            if found < offset + 1 {
                return true
            }

            matched = append(matched, word)

            if limit > 0 && found >= offset + limit  {
                return false
            }            
        }

        return true
    })

    return matched
}

func (this *SimpleDict) LookupOne(pattern string) string {
    words := this.Lookup(pattern, 0, 1)
    if len(words) == 0 {
        return ""
    }

    return words[0]
}

func (this *SimpleDict) Filter(f func(string) bool) []string {
    var matched []string
    for i := 0; i < len(this.words); i++ {
        if f(this.words[i]) {
            matched = append(matched, this.words[i])
        }
    }

    return matched
}

func (this *SimpleDict) Walk(f func(string) bool) {
    for i := 0; i < len(this.words); i++ {
        if !f(this.words[i]) {
            break
        }
    }
}