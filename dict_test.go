package dict

import (
    "testing"
)

func sampleDict() *Dict {
    d := NewDict()

    d.AddMap("Duck", map[string]interface{}{
        "Legs": 2,
        "Swim": true,
        "Fly": false,
    })

    // prop = NewDictWord()
    d.AddMap("Dog", map[string]interface{}{
        "Legs": 4,
        "Swim": true,
        "Fly": false,
    })

    d.AddMap("Snake", map[string]interface{}{
        "Legs": 0,
        "Swim": false,
        "Fly": false,
    })

    d.AddMap("Bird", map[string]interface{}{
        "Legs": 2,
        "Swim": false,
        "Fly": true,
    })

    d.AddMap("Lion", map[string]interface{}{
        "Legs": 4,
        "Swim": false,
        "Fly": false,
        "Color": "yellow",
    })

    return d
}

func TestDict(t *testing.T) {
    d := NewDict()
    
    if v := d.Count(); v != 0 {
        t.Errorf("Initial length is %d", v)
    }

    // load dict
    d.Load("/usr/share/dict/words")
    if v := d.Count(); v < 1000 {
        t.Errorf("System dict is too small: %d", v)
    }

    d.Clear()
    if v := d.Count(); v != 0 {
        t.Errorf("Cleared length is %d", v)
    }

    d = sampleDict()

    if v := d.Count(); v != 5 {
        t.Errorf("Wrong length %d != %d", v, 5)
    }

    // Export
    d.Export("/tmp/dict")

    d.Clear()

    // Load
    d.Load("/tmp/dict")
}

func TestDictLookup(t *testing.T) {
    d := NewDict()
    d.Load("/usr/share/dict/words")

    var word string
    var words []string
    var length int
    
    words = d.Lookup("", 0, 10)
    length = len(words)
    t.Log("Got: ", words)
    if length != 10 {
        t.Errorf("Lookup %d words instead of %d", length, 10)
    }

    if word = d.LookupOne(""); word == "" {
        t.Error("Lookup one word failed")
    }

    words = d.Lookup("hello", 0, 10)
    length = len(words)
    t.Log("Got: ", words)
    if length != 1 {
        t.Errorf("Unexcepted match")
    }

    if l := len(d.LookupAll("123nonexist")); l != 0 {
        t.Errorf("Unexcepted match")
    }

    w := d.LookupOne("zo{2}.*")
    t.Logf("Got: %s", w)

    if w[0:3] != "zoo" {
        t.Errorf("Unexcepted match: %s", w)
    }

    d = sampleDict()

    words = d.LookupAll(`.*o.*`)
    t.Log("Got", words)


    if words[0] != "Dog" || words[1] != "Lion" || len(words) != 2 {
        t.Errorf("Unexcepted match")
    }

    if l := len(d.Lookup(".*", 0, 0)); l != 5 {
        t.Errorf("Lookup limit error")
    }

    if l := len(d.Lookup(".*", 10, 10)); l != 0 {
        t.Errorf("Lookup limit error")
    }
}

func TestDictProp(t *testing.T) {
    d := sampleDict()
    if !d.MustGet("Duck").MustPropBool("Swim") {
        t.Errorf("Prop error")
    }
    if d.MustGet("Lion").MustPropString("Color") != "yellow" {
        t.Errorf("Prop error")
    }
    if d.MustGet("Bird").MustPropInt("Legs") != 2 {
        t.Errorf("Prop error")
    }
    if v, ok := d.Get("Human"); v != nil || ok != false {
        t.Errorf("Prop error")
    }
    if v, ok := d.MustGet("Snake").PropString("Age"); v != "" || ok != false {
        t.Errorf("Prop error")
    }
}

func TestDictFilter(t *testing.T) {
    d := sampleDict()

    list := d.Filter(func(prop *DictWord) bool {
        if p, _ := prop.PropInt("Legs"); p == 4 {
            return true
        }

        return false
    })

    if len(list) != 2 {
        t.Errorf("Filter result error")
    }

}

func TestDictWalk(t *testing.T) {
    d := sampleDict()

    length := 0
    d.Walk(func(prop *DictWord) bool {
        fly, _ := prop.PropBool("Fly")

        if prop.Word == "Lion" || fly {
            length ++
        }
        return true
    })

    if length != 2 {
        t.Errorf("Walk error")
    }
}

func TestDictSort(t *testing.T) {
    d := sampleDict()

    d.SortByWord()
    if d.LookupOne("") != "Bird" {
        t.Errorf("SortByWord error")
    }

    d.Sort(func(a *DictWord, b *DictWord) bool {
        return a.MustPropInt("Legs") < b.MustPropInt("Legs")
    })
    if d.LookupOne("") != "Snake" {
        t.Errorf("Sort error")
    }
}