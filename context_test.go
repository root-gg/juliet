package juliet

import (
	"testing"
)

func TestGet(t *testing.T) {
	ctx := NewContext()
	ctx.values["foo"] = "bar"
	ctx.values["faa"] = nil

	if val, ok := ctx.Get("foo"); ok {
		if val != "bar" {
			t.Fatalf("Invalid value for key foo. Expected 'bar' but was '%v'", val)
		}
	} else {
		t.Fatalf("Missing value for key foo")
	}

	if val, ok := ctx.Get("faa"); ok {
		if val != nil {
			t.Fatalf("Invalid value for key faa. Expected nil but was '%v'", val)
		}
	} else {
		t.Fatalf("Missing value for key faa")
	}

	if val, ok := ctx.Get("xyz"); ok {
		t.Fatalf("Invalid value '%v' for key xyz", val)
	}
}

func TestSet(t *testing.T) {
	ctx := NewContext()

	ctx.Set("foo", "bar")
	if val, ok := ctx.values["foo"]; ok {
		if val != "bar" {
			t.Fatalf("Invalid value for key foo. Expected 'bar' but was '%v'", val)
		}
	} else {
		t.Fatalf("Missing value for key foo")
	}
}

func TestDelete(t *testing.T) {
	ctx := NewContext()
	ctx.values["foo"] = "bar"

	ctx.Delete("foo")
	if val, ok := ctx.Get("foo"); ok {
		t.Fatalf("Invalid value '%v' for key foo", val)
	}

	ctx.Delete("faa")
}

func TestFork(t *testing.T) {
	ctx := NewContext()
	ctx.Set("foo", "bar")
	ctx.Set("plip", "plop")
	fork := ctx.Fork()
	fork.Set("foo", "baz")

	if val, ok := fork.Get("foo"); ok {
		if val != "baz" {
			t.Fatalf("Invalid value for key foo. Expected 'baz' but was '%v'", val)
		}
	} else {
		t.Fatalf("Missing value for key foo")
	}

	if val, ok := fork.Get("plip"); ok {
		if val != "plop" {
			t.Fatalf("Invalid value for key faa. Expected 'plop' but was '%v'", val)
		}
	} else {
		t.Fatalf("Missing value for key plop")
	}
}

func TestString(t *testing.T) {
	ctx := NewContext()
	ctx.Set("foo", "bar")
	ctx.Set("plip", "plop")
	fork := ctx.Fork()
	fork.Set("foo", "baz")
	fork2 := fork.Fork()
	fork2.Set("zig", "zag")
	if fork2.String() != "foo => bar\nplip => plop\nfoo => baz\nzig => zag\n" {
		t.Fatalf("Invalid context string representation")
	}
}
