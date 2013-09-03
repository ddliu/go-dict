package dict

import (
    "testing"
)

func sampleDict() *Dict {
    d := NewDict()
    var prop *WordProp

    prop = NewWordProp(map[string]interface{}{
        "Legs": 2,
        "Swim": true,
        "Fly": false,
    })
    d.Add("Duck", prop)

    // prop = NewWordProp()
    d.AddMap("Dog", map[string]interface{}{
        "Legs": 4,
        "Swim": true,
        "Fly": false,
    })

    prop = NewWordProp(map[string]interface{}{
        "Legs": 0,
        "Swim": false,
        "Fly": false,
    })
    d.Add("Snake", prop)

    prop = NewWordProp(map[string]interface{}{
        "Legs": 2,
        "Swim": false,
        "Fly": true,
    })
    d.Add("Bird", prop)

    prop = NewWordProp(map[string]interface{}{
        "Legs": 4,
        "Swim": false,
        "Fly": false,
        "Color": "yellow",
    })
    d.Add("Lion", prop)

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


    if /*words[0] != "Line" || words[1] != "Dog" || it's not in order :( */len(words) != 2 {
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
    if !d.MustProp("Duck").MustPropBool("Swim") {
        t.Errorf("Prop error")
    }
    if d.MustProp("Lion").MustPropString("Color") != "yellow" {
        t.Errorf("Prop error")
    }
    if d.MustProp("Bird").MustPropInt("Legs") != 2 {
        t.Errorf("Prop error")
    }
    if v, ok := d.Prop("Human"); v != nil || ok != false {
        t.Errorf("Prop error")
    }
    if v, ok := d.MustProp("Snake").PropString("Age"); v != "" || ok != false {
        t.Errorf("Prop error")
    }
}

func TestDictFilter(t *testing.T) {
    d := sampleDict()

    list := d.Filter(func(word string, prop *WordProp) bool {
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
    d.Walk(func(word string, prop *WordProp) bool {
        fly, _ := prop.PropBool("Fly")

        if word == "Lion" || fly {
            length ++
        }
        return true
    })

    if length != 2 {
        t.Errorf("Walk error")
    }
}