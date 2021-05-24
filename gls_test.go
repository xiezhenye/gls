package gls

import (
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {
	Set("my:key", "my value")
	if func0() != "my value" {
		t.Error("can not get gls value")
	}
}

func TestSimpleGo(t *testing.T) {
	Set("my:key", "my value")
	if func1() != "my value" {
		t.Error("can not get gls value")
	}
}

func func0() string {
	v, ok := Get("my:key")
	if !ok {
		return ""
	}
	return v.(string)
}

func func1() string {
	w := sync.WaitGroup{}
	w.Add(1)
	s := ""
	Go(func() {
		if v, ok := Get("my:key"); ok {
			s = v.(string)
		}
		w.Done()
	})
	w.Wait()
	return s
}


