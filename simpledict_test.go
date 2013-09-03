package dict

import (
    "testing"
)

func TestSimpleDict(t *testing.T) {
    d := NewSimpleDict()
    
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

    d.AddWords("Hello", "World")
    if v := d.Count(); v != 2 {
        t.Errorf("Wrong length %d != %d", v, 2)
    }

    d.AddWordsList([]string{"Get", "Some", "Words"})
    if v := d.Count(); v != 5 {
        t.Errorf("Wrong length %d != %d", v, 5)
    }

    // Export
    d.Export("/tmp/simpledict")
}

func TestLookup(t *testing.T) {
    d := NewSimpleDict()
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

    d.Clear()
    d.AddWords("word1", "word2", "word5", "words", "word7")
    words = d.Lookup(`word\d+`, 1, 2)
    t.Log("Got: ", words)
    if words[0] != "word2" || words[1] != "word5" || len(words) != 2 {
        t.Errorf("Lookup limit error")
    }

    if l := len(d.Lookup("word.*", 0, 0)); l != 5 {
        t.Errorf("Lookup limit error")
    }

    if l := len(d.Lookup("word.*", 10, 10)); l != 0 {
        t.Errorf("Lookup limit error")
    }
}

func TestFilter(t *testing.T) {
    d := NewSimpleDict()
    d.AddWords("word1", "word2", "xxxx4", "word3")
    
    list := d.Filter(func(word string) bool {
        if word[0:4] == "word" {
            return true
        }
        return false
    })

    if len(list) != 3 {
        t.Errorf("Filter result error")
    }

}

func TestWalk(t *testing.T) {
    d := NewSimpleDict()
    d.AddWords("word1", "word2", "xxxx4", "word3")

    length := 0
    d.Walk(func(word string) bool {
        length ++
        return true
    })

    if length != 4 {
        t.Errorf("Walk error")
    }
}