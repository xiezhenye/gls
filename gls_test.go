package gls

import (
	"context"
	"runtime/pprof"
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {
	Set("my:key1", "my value1")
	if func0() != "my value1" {
		t.Error("can not get gls value")
	}
}

func TestSimpleGo(t *testing.T) {
	Set("my:key2", "my value2")
	if func1() != "my value2" {
		t.Error("can not get gls value")
	}
}

func TestContext(t *testing.T) {
	newCtx := AttachGls(context.Background())
	if _, ok := pprof.Label(newCtx, glsLabel); !ok {
		t.Error("missing glsLabel")
	}
	valueCtx := context.WithValue(newCtx, "test", "test")
	if _, ok := pprof.Label(valueCtx, glsLabel); !ok {
		t.Error("missing glsLabel")
	}
	TestSimple(t)
}

func TestContextGo(t *testing.T) {
	ctx := context.Background()
	Set("my:key3", "my value3")
	s := ""
	w := sync.WaitGroup{}
	w.Add(1)
	GoWithContext(ctx, func(newCtx context.Context) {
		if _, ok := pprof.Label(newCtx, glsLabel); !ok {
			t.Error("missing glsLabel")
		}
		if v, ok := Get("my:key3"); ok {
			s = v.(string)
		}
		w.Done()
	})
	w.Wait()
	if s != "my value3" {
		t.Error("can not get gls value")
	}
}

func func0() string {
	v, ok := Get("my:key1")
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
		if v, ok := Get("my:key2"); ok {
			s = v.(string)
		}
		w.Done()
	})
	w.Wait()
	return s
}


