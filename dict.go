package dict

import (
    "os"
    "encoding/json"
    "strings"
    "regexp"
)

/**
 * WordProp
 */

type WordProp struct {
    props map[string] interface {}
}

func NewWordProp(props map[string]interface{}) *WordProp{
    return &WordProp { props }
}

func NewWordPropJSON(data string) *WordProp {
    instance := new(WordProp)
    instance.FromJSON(data)
    return instance
}

func (this *WordProp) Prop(name string) (interface{}, bool) {
    prop, ok := this.props[name]

    return prop, ok
}

func (this *WordProp) MustProp(name string) (interface{}) {
    prop, ok := this.props[name]

    if !ok {
        panic("Prop not exist: " + name)
    }

    return prop
}

func (this *WordProp) PropInt(name string) (int, bool) {
    prop, ok := this.props[name]
    if !ok {
        return 0, false
    }
    propint, ok := prop.(int)
    return propint, ok
}

func (this *WordProp) MustPropInt(name string) (int) {
    prop, ok := this.PropInt(name)
    if !ok {
        panic ("Prop not exist: " + name)
    }

    return prop
}

func (this *WordProp) PropString(name string) (string, bool) {
    prop, ok := this.props[name]
    if !ok {
        return "", false
    }
    propstring, ok := prop.(string)

    return propstring, ok
}

func (this *WordProp) MustPropString(name string) (string) {
    prop, ok := this.PropString(name)
    if !ok {
        panic ("Prop not exist: " + name)
    }

    return prop
}

func (this *WordProp) PropBool(name string) (bool, bool) {
    prop, ok := this.props[name]
    if !ok {
        return false, false
    }

    propbool, ok := prop.(bool)

    return propbool, ok
}

func (this *WordProp) MustPropBool(name string) (bool) {
    prop, ok := this.PropBool(name)
    if !ok {
        panic ("Prop not exist: " + name)
    }

    return prop
}

func (this *WordProp) SetProp(name string, value interface{}) {
    this.props[name] = value
}

func (this *WordProp) ToJSON() (string, error) {
    data, err := json.Marshal(&this.props)

    return string(data), err
}

func (this *WordProp) FromJSON(data string) error {
    return json.Unmarshal([]byte(data), this.props)
}

/**
 * Dict
 */

func NewDict() *Dict {
    return &Dict {make(map[string] *WordProp)}
}

type Dict struct {
    words map[string] *WordProp
}

func (this *Dict) Load(file string) {
    WalkFileLines(file, func(line string) bool{
        parts := strings.Split(string(line), "\t")
        var prop *WordProp
        if len(parts) > 1 {
            prop = NewWordPropJSON(parts[1])
        } else {
            prop = NewWordProp(map[string]interface{}{})
        }

        this.words[parts[0]] = prop

        return true
    })
}

func (this *Dict) Export(file string) {
    fi, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0660)

    if err != nil {
        panic(err)
    }

    defer fi.Close()

    for k, v := range this.words {
        jsonString, err := v.ToJSON()
        if err != nil {
            panic(err)
        }
        line := k + "\t" + jsonString + "\n"
        if _, err := fi.WriteString(line); err != nil {
            panic(err)
        }

    }
}

func (this *Dict) Add(name string, prop *WordProp) {
    this.words[name] = prop
}

func (this *Dict) AddMap(name string, m map[string]interface{}) {
    this.words[name] = NewWordProp(m)
}

func (this *Dict) Clear() {
    this.words = map[string] *WordProp{}
}

func (this *Dict) Prop(word string) (*WordProp, bool) {
    v, ok := this.words[word]

    return v, ok
}

func (this *Dict) MustProp(word string) *WordProp {
    v, ok := this.words[word]
    if !ok {
        panic("Word not exist: " + word)
    }

    return v
}

func (this *Dict) Count() int {
    return len(this.words)
}

func (this *Dict) LookupAll(pattern string) []string {
    return this.Lookup(pattern, 0, 0)
}

func (this *Dict) Lookup(pattern string, offset int, limit int) []string {
    var compiledPattern = regexp.MustCompile("^" + pattern + "$")
    var matched []string
    found := 0
    this.Walk(func(word string, prop *WordProp) bool {
        if pattern == "" || compiledPattern.MatchString(word) {
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

func (this *Dict) LookupOne(pattern string) string {
    words := this.Lookup(pattern, 0, 1)
    if len(words) == 0 {
        return ""
    }

    return words[0]
}

func (this *Dict) Filter(f func(string, *WordProp) bool) []string {
    var words []string
    for k, v := range this.words {
        if f(k, v) {
            words = append(words, k)
        }
    }

    return words
}

func (this *Dict) Walk(f func(string, *WordProp) bool) {
    for k, v := range this.words {
        if !f(k, v) {
            break
        }
    }
}