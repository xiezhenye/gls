package gls

import (
	"context"
	"reflect"
	"runtime/pprof"
	"unsafe"
)

const glsLabel = "$github.com/xiezhenye/gls/hack_label$"

//go:linkname runtime_getProfLabel runtime/pprof.runtime_getProfLabel
func runtime_getProfLabel() unsafe.Pointer

//go:linkname runtime_setProfLabel runtime/pprof.runtime_setProfLabel
func runtime_setProfLabel(labels unsafe.Pointer)

func getLabels() map[string]string {
	mp := runtime_getProfLabel()
	if mp == nil {
		m := make(map[string]string)
		mp = unsafe.Pointer(&m)
		runtime_setProfLabel(mp)
	}
	return *(*map[string]string)(mp)
}

//go:uintptrescapes
func pointerToHackStr(p uintptr) string {
	sh := &reflect.StringHeader {
		Data: p,
		Len: 0,
	}
	return *(*string)(unsafe.Pointer(sh))
}

func pointerFromHackStr(s string) uintptr {
	return (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
}

func AttachGls(ctx context.Context) (newCtx context.Context) {
	lbs := getLabels()
	if ctx == nil {
		ctx = context.Background()
	}
	if hs, ok := lbs[glsLabel]; ok {
		if _, ok = pprof.Label(ctx, glsLabel); ok {
			newCtx = ctx
		} else {
			newCtx = pprof.WithLabels(ctx, pprof.Labels(glsLabel, hs))
		}
		return
	}
	m := make(map[string]interface{})
	p := *(*uintptr)(unsafe.Pointer(&m))
	hs := pointerToHackStr(p)
	newCtx = pprof.WithLabels(ctx, pprof.Labels(glsLabel, hs))
	pprof.SetGoroutineLabels(newCtx)
	return
}

func Get(key string) (interface{}, bool) {
	m := GetGlsMap()
	ret, ok := m[key]
	return ret, ok
}

func Set(key string, value interface{}) {
	lbs := GetGlsMap()
	lbs[key] = value
}

func GetGlsMap() map[string]interface{} {
	lbs := getLabels()
	hs, ok := lbs[glsLabel]
	if !ok {
		m := make(map[string]interface{})
		p := *(*uintptr)(unsafe.Pointer(&m))
		hs := pointerToHackStr(p)
		lbs[glsLabel] = hs
		return m
	}
	p := pointerFromHackStr(hs)
	return *(*map[string]interface{})(unsafe.Pointer(&p))
}

func Go(f func()) {
	m := GetGlsMap()
	go func(m map[string]interface{}) {
		newM := GetGlsMap()
		for k, v := range m {
			newM[k] = v
		}
		f()
	}(m)
}

func GoWithContext(ctx context.Context, f func(ctx context.Context)) {
	m := GetGlsMap()
	go func(m map[string]interface{}) {
		newM := GetGlsMap()
		for k, v := range m {
			newM[k] = v
		}
		newCtx := AttachGls(ctx)
		f(newCtx)
	}(m)
}

