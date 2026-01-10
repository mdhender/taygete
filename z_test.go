// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"log/slog"
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func TestIListAppendPrepend(t *testing.T) {
	var l []int

	IListAppend(&l, 1)
	IListAppend(&l, 2)
	IListAppend(&l, 3)

	if IListLen(l) != 3 {
		t.Errorf("expected len 3, got %d", IListLen(l))
	}
	if l[0] != 1 || l[1] != 2 || l[2] != 3 {
		t.Errorf("expected [1,2,3], got %v", l)
	}

	IListPrepend(&l, 0)
	if IListLen(l) != 4 {
		t.Errorf("expected len 4, got %d", IListLen(l))
	}
	if l[0] != 0 || l[1] != 1 || l[2] != 2 || l[3] != 3 {
		t.Errorf("expected [0,1,2,3], got %v", l)
	}
}

func TestIListDelete(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}

	IListDelete(&l, 2)
	if IListLen(l) != 4 {
		t.Errorf("expected len 4, got %d", IListLen(l))
	}
	if l[0] != 1 || l[1] != 2 || l[2] != 4 || l[3] != 5 {
		t.Errorf("expected [1,2,4,5], got %v", l)
	}

	IListDelete(&l, 0)
	if l[0] != 2 {
		t.Errorf("expected first element 2, got %d", l[0])
	}

	IListDelete(&l, IListLen(l)-1)
	if l[IListLen(l)-1] != 4 {
		t.Errorf("expected last element 4, got %d", l[IListLen(l)-1])
	}
}

func TestIListLookup(t *testing.T) {
	l := []int{10, 20, 30, 40, 50}

	if idx := IListLookup(l, 30); idx != 2 {
		t.Errorf("expected index 2 for 30, got %d", idx)
	}
	if idx := IListLookup(l, 10); idx != 0 {
		t.Errorf("expected index 0 for 10, got %d", idx)
	}
	if idx := IListLookup(l, 50); idx != 4 {
		t.Errorf("expected index 4 for 50, got %d", idx)
	}
	if idx := IListLookup(l, 999); idx != -1 {
		t.Errorf("expected -1 for missing value, got %d", idx)
	}
	if idx := IListLookup(nil, 10); idx != -1 {
		t.Errorf("expected -1 for nil list, got %d", idx)
	}
}

func TestIListRemValue(t *testing.T) {
	l := []int{1, 2, 3, 2, 4, 2, 5}

	IListRemValue(&l, 2)
	if IListLen(l) != 4 {
		t.Errorf("expected len 4 after removing all 2s, got %d", IListLen(l))
	}
	if IListLookup(l, 2) != -1 {
		t.Errorf("expected no 2s remaining, but found one")
	}

	l2 := []int{1, 2, 3, 2, 4}
	IListRemValueUniq(&l2, 2)
	if IListLen(l2) != 4 {
		t.Errorf("expected len 4 after removing first 2, got %d", IListLen(l2))
	}
	if IListLookup(l2, 2) == -1 {
		t.Errorf("expected one 2 remaining")
	}
}

func TestIListCopy(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	cp := IListCopy(l)

	if IListLen(cp) != IListLen(l) {
		t.Errorf("copy length mismatch: %d vs %d", IListLen(cp), IListLen(l))
	}
	for i := range l {
		if l[i] != cp[i] {
			t.Errorf("copy mismatch at %d: %d vs %d", i, l[i], cp[i])
		}
	}

	cp[0] = 999
	if l[0] == 999 {
		t.Errorf("modifying copy affected original")
	}

	if IListCopy(nil) != nil {
		t.Errorf("copy of nil should be nil")
	}
}

func TestIListScramble(t *testing.T) {
	teg = &Engine{
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(12_345, 67_890)),
	}

	l := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	original := IListCopy(l)

	IListScramble(l)

	if IListLen(l) != IListLen(original) {
		t.Errorf("scramble changed length")
	}

	for _, v := range original {
		if IListLookup(l, v) == -1 {
			t.Errorf("scramble lost value %d", v)
		}
	}
}

func TestIListInsert(t *testing.T) {
	l := []int{1, 2, 4, 5}

	IListInsert(&l, 2, 3)
	if IListLen(l) != 5 {
		t.Errorf("expected len 5, got %d", IListLen(l))
	}
	expected := []int{1, 2, 3, 4, 5}
	for i, v := range expected {
		if l[i] != v {
			t.Errorf("at %d: expected %d, got %d", i, v, l[i])
		}
	}
}

func TestIListClearReclaim(t *testing.T) {
	l := []int{1, 2, 3}

	IListClear(&l)
	if IListLen(l) != 0 {
		t.Errorf("expected len 0 after clear, got %d", IListLen(l))
	}

	l = []int{1, 2, 3}
	IListReclaim(&l)
	if l != nil {
		t.Errorf("expected nil after reclaim")
	}
}

func TestPListAppendPrepend(t *testing.T) {
	var l []any

	PListAppend(&l, "a")
	PListAppend(&l, "b")
	PListAppend(&l, "c")

	if PListLen(l) != 3 {
		t.Errorf("expected len 3, got %d", PListLen(l))
	}

	PListPrepend(&l, "z")
	if PListLen(l) != 4 {
		t.Errorf("expected len 4, got %d", PListLen(l))
	}
	if l[0] != "z" {
		t.Errorf("expected first element 'z', got %v", l[0])
	}
}

func TestPListDelete(t *testing.T) {
	l := []any{"a", "b", "c", "d"}

	PListDelete(&l, 1)
	if PListLen(l) != 3 {
		t.Errorf("expected len 3, got %d", PListLen(l))
	}
	if l[1] != "c" {
		t.Errorf("expected 'c' at index 1, got %v", l[1])
	}
}

func TestPListLookup(t *testing.T) {
	a, b, c := "a", "b", "c"
	l := []any{&a, &b, &c}

	if idx := PListLookup(l, &b); idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
	if idx := PListLookup(l, nil); idx != -1 {
		t.Errorf("expected -1 for nil, got %d", idx)
	}
	if idx := PListLookup(nil, &a); idx != -1 {
		t.Errorf("expected -1 for nil list, got %d", idx)
	}
}

func TestPListCopy(t *testing.T) {
	l := []any{1, "two", 3.0}
	cp := PListCopy(l)

	if PListLen(cp) != PListLen(l) {
		t.Errorf("copy length mismatch")
	}

	if PListCopy(nil) != nil {
		t.Errorf("copy of nil should be nil")
	}
}

func Test_i_strcmp(t *testing.T) {
	tests := []struct {
		s, t string
		want int
	}{
		{"abc", "abc", 0},
		{"ABC", "abc", 0},
		{"abc", "ABC", 0},
		{"abc", "abd", -1},
		{"abd", "abc", 1},
		{"", "", 0},
		{"a", "", 1},
		{"", "a", -1},
	}

	for _, tt := range tests {
		got := i_strcmp(tt.s, tt.t)
		if (tt.want == 0 && got != 0) || (tt.want < 0 && got >= 0) || (tt.want > 0 && got <= 0) {
			t.Errorf("i_strcmp(%q, %q) = %d, want sign of %d", tt.s, tt.t, got, tt.want)
		}
	}
}

func Test_i_strncmp(t *testing.T) {
	tests := []struct {
		s, t string
		n    int
		want int
	}{
		{"abcdef", "abcxxx", 3, 0},
		{"ABCDEF", "abcxxx", 3, 0},
		{"abcdef", "abdxxx", 3, -1},
		{"abcdef", "abcdef", 10, 0},
		{"abc", "abcdef", 3, 0},
		{"", "", 5, 0},
	}

	for _, tt := range tests {
		got := i_strncmp(tt.s, tt.t, tt.n)
		if (tt.want == 0 && got != 0) || (tt.want < 0 && got >= 0) || (tt.want > 0 && got <= 0) {
			t.Errorf("i_strncmp(%q, %q, %d) = %d, want sign of %d", tt.s, tt.t, tt.n, got, tt.want)
		}
	}
}

func Test_fuzzy_strcmp(t *testing.T) {
	tests := []struct {
		one, two string
		want     bool
	}{
		{"abcd", "abdc", true},
		{"test", "tset", true},
		{"helllo", "hello", true},
		{"hello", "helllo", true},
		{"hellox", "hello", true},
		{"hello", "hallo", true},
		{"abc", "abc", false},
		{"ab", "ba", false},
		{"abc", "xyz", false},
		{"abcdefgh", "abcdefgi", true},
	}

	for _, tt := range tests {
		got := fuzzy_strcmp(tt.one, tt.two)
		if got != tt.want {
			t.Errorf("fuzzy_strcmp(%q, %q) = %v, want %v", tt.one, tt.two, got, tt.want)
		}
	}
}

func Test_iswhite(t *testing.T) {
	if !iswhite(' ') {
		t.Error("space should be whitespace")
	}
	if !iswhite('\t') {
		t.Error("tab should be whitespace")
	}
	if iswhite('a') {
		t.Error("'a' should not be whitespace")
	}
	if iswhite('\n') {
		t.Error("newline should not be whitespace (per C definition)")
	}
}

func Test_int_comp(t *testing.T) {
	if int_comp(1, 2) >= 0 {
		t.Error("1 < 2 should return negative")
	}
	if int_comp(2, 1) <= 0 {
		t.Error("2 > 1 should return positive")
	}
	if int_comp(5, 5) != 0 {
		t.Error("5 == 5 should return 0")
	}
}

func Test_lcase(t *testing.T) {
	if lcase("HELLO") != "hello" {
		t.Errorf("lcase(HELLO) = %q, want hello", lcase("HELLO"))
	}
	if lcase("Hello World") != "hello world" {
		t.Errorf("lcase failed on mixed case")
	}
}
