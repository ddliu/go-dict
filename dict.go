package dict

import (
    "os"
    "encoding/json"
    "strings"
    "regexp"
    "fmt"
    // "reflect"
)

/**
 * DictWord
 */

type DictWord struct {
    Word string
    props map[string] interface {}
}

func NewDictWord(word string, props map[string]interface{}) *DictWord{
    return &DictWord {word, props}
}

func NewDictWordJSON(word string, data string) *DictWord {
    instance := new(DictWord)
    instance.Word = word
    instance.FromJSON(data)
    return instance
}

func (this *DictWord) Prop(name string) (interface{}, bool) {
    prop, ok := this.props[name]

    return prop, ok
}

func (this *DictWord) MustProp(name string) (interface{}) {
    prop, ok := this.props[name]

    if !ok {
        panic("Prop not exist: " + name)
    }

    return prop
}

func (this *DictWord) PropInt(name string) (int, bool) {
    prop, ok := this.props[name]
    if !ok {
        return 0, false
    }

    return ToInt(prop)
}

func (this *DictWord) MustPropInt(name string) (int) {
    prop, ok := this.PropInt(name)
    if !ok {
        panic (fmt.Sprintf("Prop not exist: %s", name))
    }

    return prop
}

func (this *DictWord) PropInt32(name string) (int32, bool) {
    prop, ok := this.props[name]
    if !ok {
        return 0, false
    }

    return ToInt32(prop)
}

func (this *DictWord) MustPropInt32(name string) (int32) {
    prop, ok := this.PropInt32(name)
    if !ok {
        panic (fmt.Sprintf("Prop not exist: %s", name))
    }

    return prop
}

func (this *DictWord) PropInt64(name string) (int64, bool) {
    prop, ok := this.props[name]
    if !ok {
        return 0, false
    }

    return ToInt64(prop)
}

func (this *DictWord) MustPropInt64(name string) (int64) {
    prop, ok := this.PropInt64(name)
    if !ok {
        panic (fmt.Sprintf("Prop not exist: %s", name))
    }

    return prop
}

func (this *DictWord) PropString(name string) (string, bool) {
    prop, ok := this.props[name]
    if !ok {
        return "", false
    }
    propstring, ok := prop.(string)

    return propstring, ok
}

func (this *DictWord) MustPropString(name string) (string) {
    prop, ok := this.PropString(name)
    if !ok {
        panic ("Prop not exist: " + name)
    }

    return prop
}

func (this *DictWord) PropBool(name string) (bool, bool) {
    prop, ok := this.props[name]
    if !ok {
        return false, false
    }

    propbool, ok := prop.(bool)

    return propbool, ok
}

func (this *DictWord) MustPropBool(name string) (bool) {
    prop, ok := this.PropBool(name)
    if !ok {
        panic ("Prop not exist: " + name)
    }

    return prop
}

func (this *DictWord) SetProp(name string, value interface{}) {
    this.props[name] = value
}

func (this *DictWord) ToJSON() (string, error) {
    data, err := json.Marshal(&this.props)

    return string(data), err
}

func (this *DictWord) FromJSON(data string) error {
    return json.Unmarshal([]byte(data), &this.props)
}

/**
 * Dict
 */

func NewDict() *Dict {
    return &Dict {[]*DictWord{}, map[string]int{}}
}

type Dict struct {
    words []*DictWord
    indexes map[string]int
}

func (this *Dict) Load(file string) {
    WalkFileLines(file, func(line string) bool{
        parts := strings.Split(line, "\t")
        var prop *DictWord
        if len(parts) > 1 {
            prop = NewDictWordJSON(parts[0], parts[1])
        } else {
            prop = NewDictWord(parts[0], map[string]interface{}{})
        }

        this.Add(prop)

        return true
    })
}

func (this *Dict) Export(file string) {
    fi, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)

    if err != nil {
        panic(err)
    }

    defer fi.Close()

    for _, v := range this.words {
        jsonString, err := v.ToJSON()
        if err != nil {
            panic(err)
        }
        line := v.Word + "\t" + jsonString + "\n"
        if _, err := fi.WriteString(line); err != nil {
            panic(err)
        }

    }
}

func (this *Dict) Add(prop *DictWord) {
    this.words = append(this.words, prop)
    this.indexes[prop.Word] = len(this.words) - 1
}

func (this *Dict) AddMap(word string, m map[string]interface{}) {
    this.Add(NewDictWord(word, m))
}

func (this *Dict) Clear() {
    this.words = []*DictWord{}
    this.indexes = map[string]int{}
}

func (this *Dict) Get(word string) (*DictWord, bool) {
    index, ok := this.indexes[word]

    if !ok {
        return nil, false
    }

    return this.words[index], true
}

func (this *Dict) MustGet(word string) *DictWord {
    prop, ok := this.Get(word)

    if !ok {
        panic("Word not exist: " + word)
    }

    return prop
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
    this.Walk(func(dw *DictWord) bool {
        if pattern == "" || compiledPattern.MatchString(dw.Word) {
            found++

            if found < offset + 1 {
                return true
            }

            matched = append(matched, dw.Word)

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

func (this *Dict) Filter(f func(*DictWord) bool) []string {
    var words []string
    for _, v := range this.words {
        if f(v) {
            words = append(words, v.Word)
        }
    }

    return words
}

func (this *Dict) Walk(f func(*DictWord) bool) {
    for _, v := range this.words {
        if !f(v) {
            break
        }
    }
}