package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

const scriptA = `

println("in script A")

a = 1
b = 2

foo = fn(a) {
	println("in func foo:", a, b)
}

export b, foo
`

const scriptB = `

include "a.ql"

println("in script B")
foo(a)
`

const scriptC = `

import "a"
import "a" as g

b = 3

set(g, "b", 4)
println("in script C:", a.b, g.b)
a.foo(b)
`

func TestInclude(t *testing.T) {

	lang, _ := qlang.New(qlang.InsertSemis)

	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA), nil
	})

	err := lang.SafeExec([]byte(scriptB), "b.ql")
	if err != nil {
		t.Fatal(err)
	}
}

func TestImport(t *testing.T) {

	lang, _ := qlang.New(qlang.InsertSemis)

	qlang.SetFindEntry(func(file string, libs []string) (string, error) {
		return file, nil
	})

	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA), nil
	})

	err := lang.SafeExec([]byte(scriptC), "c.ql")
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------

const scriptA1 = `

defer fn() {
	x; x = 2
}()

x = 1
export x
`

const scriptB1 = `

import "a"

println("a.x:", a.x)
`

func TestModuleDefer(t *testing.T) {

	lang, _ := qlang.New(qlang.InsertSemis)

	qlang.SetFindEntry(func(file string, libs []string) (string, error) {
		return file, nil
	})

	qlang.SetReadFile(func(file string) ([]byte, error) {
		return []byte(scriptA1), nil
	})

	err := lang.SafeExec([]byte(scriptB1), "b.ql")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := lang.Var("a"); !ok || v.(map[string]interface{})["x"] != 2 {
		t.Fatal("x != 2, x =", v)
	}
}

// -----------------------------------------------------------------------------

