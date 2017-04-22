package drafthouse

import (
	"testing"
	"time"
)

type Foo struct {
	name string
}

func TestGet(t *testing.T) {
	cache := NewCache(2 * time.Second)

	data, exists := cache.Get("hello")
	if exists || data != "" {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("mykey", "test")
	data, exists = cache.Get("mykey")
	if !exists {
		t.Errorf("Expected cache to return data for `mykey`")
	}
	if data != "test" {
		t.Errorf("Expected cache to return `test` for `mykey`")
	}

	f := Foo{name: "bob"}
	cache.Set("myfoo", f)
	data, exists = cache.Get("myfoo")
	if !exists {
		t.Errorf("Expected cache to return data for `myfoo`")
	}
	if data.(Foo).name != "bob" {
		t.Errorf("Failed to get my struct back for `myfoo`")
	}

	time.Sleep(2 * time.Second)
	data, exists = cache.Get("myfoo")
	if exists {
		t.Errorf("Expected cache to expire data for `myfoo`")
	}

	data, exists = cache.Get("mykey")
	if exists {
		t.Errorf("Expected cache to expire data for `mykey`")
	}

}
