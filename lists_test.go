// taygete - a game engine for a game.
// Copyright (c) 2026 Michael D Henderson.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package taygete

import (
	"reflect"
	"testing"
)

func TestNewList(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		l := NewList[int]()
		if l.Len() != 0 {
			t.Errorf("expected len 0, got %d", l.Len())
		}
		if l.Values() != nil {
			t.Errorf("expected nil values for empty list")
		}
	})

	t.Run("with values", func(t *testing.T) {
		l := NewList(1, 2, 3)
		if l.Len() != 3 {
			t.Errorf("expected len 3, got %d", l.Len())
		}
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})

	t.Run("strings", func(t *testing.T) {
		l := NewList("a", "b", "c")
		if l.Len() != 3 {
			t.Errorf("expected len 3, got %d", l.Len())
		}
		if !reflect.DeepEqual(l.Values(), []string{"a", "b", "c"}) {
			t.Errorf("expected [a,b,c], got %v", l.Values())
		}
	})

	t.Run("does not share memory", func(t *testing.T) {
		orig := []int{1, 2, 3}
		l := NewList(orig...)
		orig[0] = 999
		if l.Values()[0] == 999 {
			t.Error("NewList should not share memory with input")
		}
	})
}

func TestList_Append(t *testing.T) {
	var l IList
	l.Append(1)
	l.Append(2)
	l.Append(3)

	if l.Len() != 3 {
		t.Errorf("expected len 3, got %d", l.Len())
	}
	if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
		t.Errorf("expected [1,2,3], got %v", l.Values())
	}
}

func TestList_Prepend(t *testing.T) {
	var l IList
	l.Prepend(3)
	l.Prepend(2)
	l.Prepend(1)

	if l.Len() != 3 {
		t.Errorf("expected len 3, got %d", l.Len())
	}
	if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
		t.Errorf("expected [1,2,3], got %v", l.Values())
	}
}

func TestList_Delete(t *testing.T) {
	t.Run("delete first", func(t *testing.T) {
		l := NewList(1, 2, 3)
		l.Delete(0)
		if !reflect.DeepEqual(l.Values(), []int{2, 3}) {
			t.Errorf("expected [2,3], got %v", l.Values())
		}
	})

	t.Run("delete middle", func(t *testing.T) {
		l := NewList(1, 2, 3)
		l.Delete(1)
		if !reflect.DeepEqual(l.Values(), []int{1, 3}) {
			t.Errorf("expected [1,3], got %v", l.Values())
		}
	})

	t.Run("delete last", func(t *testing.T) {
		l := NewList(1, 2, 3)
		l.Delete(2)
		if !reflect.DeepEqual(l.Values(), []int{1, 2}) {
			t.Errorf("expected [1,2], got %v", l.Values())
		}
	})

	t.Run("panics on negative index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative index")
			}
		}()
		l := NewList(1, 2, 3)
		l.Delete(-1)
	})

	t.Run("panics on out of bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds")
			}
		}()
		l := NewList(1, 2, 3)
		l.Delete(3)
	})
}

func TestList_Clear(t *testing.T) {
	l := NewList(1, 2, 3)
	l.Clear()

	if l.Len() != 0 {
		t.Errorf("expected len 0 after clear, got %d", l.Len())
	}

	// Underlying capacity should be retained
	l.Append(4)
	if l.Len() != 1 {
		t.Errorf("expected len 1 after append, got %d", l.Len())
	}
}

func TestList_Reclaim(t *testing.T) {
	l := NewList(1, 2, 3)
	l.Reclaim()

	if l.Len() != 0 {
		t.Errorf("expected len 0 after reclaim, got %d", l.Len())
	}
	if l.Values() != nil {
		t.Error("expected nil values after reclaim")
	}
}

func TestList_Lookup(t *testing.T) {
	l := NewList(10, 20, 30)

	if idx := l.Lookup(20); idx != 1 {
		t.Errorf("expected index 1 for value 20, got %d", idx)
	}

	if idx := l.Lookup(40); idx != -1 {
		t.Errorf("expected index -1 for missing value, got %d", idx)
	}
}

func TestList_IndexOf(t *testing.T) {
	l := NewList(10, 20, 30)

	idx, found := l.IndexOf(20)
	if !found || idx != 1 {
		t.Errorf("expected (1, true) for value 20, got (%d, %v)", idx, found)
	}

	idx, found = l.IndexOf(40)
	if found || idx != -1 {
		t.Errorf("expected (-1, false) for missing value, got (%d, %v)", idx, found)
	}
}

func TestList_RemValue(t *testing.T) {
	t.Run("removes all occurrences", func(t *testing.T) {
		l := NewList(1, 2, 3, 2, 4, 2)
		l.RemValue(2)
		if !reflect.DeepEqual(l.Values(), []int{1, 3, 4}) {
			t.Errorf("expected [1,3,4], got %v", l.Values())
		}
	})

	t.Run("no match", func(t *testing.T) {
		l := NewList(1, 2, 3)
		l.RemValue(5)
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var l IList
		l.RemValue(1) // should not panic
		if l.Len() != 0 {
			t.Errorf("expected len 0, got %d", l.Len())
		}
	})
}

func TestList_RemValueUniq(t *testing.T) {
	t.Run("removes last occurrence from back", func(t *testing.T) {
		l := NewList(1, 2, 3, 2, 4)
		l.RemValueUniq(2)
		// Should remove the LAST 2 (at index 3), not the first
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3, 4}) {
			t.Errorf("expected [1,2,3,4], got %v", l.Values())
		}
	})

	t.Run("no match", func(t *testing.T) {
		l := NewList(1, 2, 3)
		l.RemValueUniq(5)
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var l IList
		l.RemValueUniq(1) // should not panic
		if l.Len() != 0 {
			t.Errorf("expected len 0, got %d", l.Len())
		}
	})
}

func TestList_Copy(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var l *IList
		cp := l.Copy()
		if cp != nil {
			t.Error("expected nil copy from nil receiver")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var l IList
		cp := l.Copy()
		if cp == nil {
			t.Fatal("expected non-nil copy")
		}
		if cp.Len() != 0 {
			t.Errorf("expected len 0, got %d", cp.Len())
		}
	})

	t.Run("with values", func(t *testing.T) {
		l := NewList(1, 2, 3)
		cp := l.Copy()

		if !reflect.DeepEqual(cp.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", cp.Values())
		}

		// Modify original, copy should be unchanged
		l.Append(4)
		if len(cp.Values()) != 3 {
			t.Error("copy should not share memory with original")
		}
	})
}

func TestList_Insert(t *testing.T) {
	t.Run("insert at beginning", func(t *testing.T) {
		l := NewList(2, 3)
		l.Insert(0, 1)
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})

	t.Run("insert at middle", func(t *testing.T) {
		l := NewList(1, 3)
		l.Insert(1, 2)
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})

	t.Run("insert at end", func(t *testing.T) {
		l := NewList(1, 2)
		l.Insert(2, 3)
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})

	t.Run("insert into empty", func(t *testing.T) {
		var l IList
		l.Insert(0, 1)
		if !reflect.DeepEqual(l.Values(), []int{1}) {
			t.Errorf("expected [1], got %v", l.Values())
		}
	})

	t.Run("panics on negative position", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative position")
			}
		}()
		l := NewList(1, 2, 3)
		l.Insert(-1, 0)
	})

	t.Run("panics on out of bounds position", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds position")
			}
		}()
		l := NewList(1, 2, 3)
		l.Insert(4, 0)
	})
}

func TestList_Scramble(t *testing.T) {
	l := NewList(1, 2, 3, 4, 5)
	original := make([]int, 5)
	copy(original, l.Values())

	l.Scramble()

	// Check all elements are preserved
	if l.Len() != 5 {
		t.Errorf("expected len 5 after scramble, got %d", l.Len())
	}

	// Check all original values still present
	for _, v := range original {
		if l.Lookup(v) == -1 {
			t.Errorf("value %d missing after scramble", v)
		}
	}
}

func TestSortList(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		l := NewList(3, 1, 4, 1, 5, 9, 2, 6)
		SortList(&l)
		if !reflect.DeepEqual(l.Values(), []int{1, 1, 2, 3, 4, 5, 6, 9}) {
			t.Errorf("expected sorted list, got %v", l.Values())
		}
	})

	t.Run("strings", func(t *testing.T) {
		l := NewList("cherry", "apple", "banana")
		SortList(&l)
		if !reflect.DeepEqual(l.Values(), []string{"apple", "banana", "cherry"}) {
			t.Errorf("expected sorted list, got %v", l.Values())
		}
	})

	t.Run("nil list", func(t *testing.T) {
		var l *List[int]
		SortList(l) // should not panic
	})

	t.Run("empty list", func(t *testing.T) {
		var l IList
		SortList(&l) // should not panic
		if l.Len() != 0 {
			t.Errorf("expected len 0, got %d", l.Len())
		}
	})
}

func TestSortListFunc(t *testing.T) {
	l := NewList(3, 1, 4, 1, 5)
	// Sort descending
	SortListFunc(&l, func(a, b int) int { return b - a })
	if !reflect.DeepEqual(l.Values(), []int{5, 4, 3, 1, 1}) {
		t.Errorf("expected descending order, got %v", l.Values())
	}
}

func TestList_Values(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var l *IList
		if l.Values() != nil {
			t.Error("expected nil from nil receiver")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var l IList
		if l.Values() != nil {
			t.Error("expected nil from empty list")
		}
	})

	t.Run("with values", func(t *testing.T) {
		l := NewList(1, 2, 3)
		if !reflect.DeepEqual(l.Values(), []int{1, 2, 3}) {
			t.Errorf("expected [1,2,3], got %v", l.Values())
		}
	})
}

func TestList_Len(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var l *IList
		if l.Len() != 0 {
			t.Error("expected 0 from nil receiver")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var l IList
		if l.Len() != 0 {
			t.Errorf("expected 0, got %d", l.Len())
		}
	})

	t.Run("with values", func(t *testing.T) {
		l := NewList(1, 2, 3)
		if l.Len() != 3 {
			t.Errorf("expected 3, got %d", l.Len())
		}
	})
}

// Test IList alias works correctly
func TestIListAlias(t *testing.T) {
	var l IList
	l.Append(1)
	l.Append(2)

	if l.Len() != 2 {
		t.Errorf("IList alias: expected len 2, got %d", l.Len())
	}
	if !reflect.DeepEqual(l.Values(), []int{1, 2}) {
		t.Errorf("IList alias: expected [1,2], got %v", l.Values())
	}
}

// Test IList as struct field works with zero value (regression for Task 26.12.3)
func TestIListStructFields(t *testing.T) {
	t.Run("entity_player.deliver_lore", func(t *testing.T) {
		var p entity_player
		// Zero value should be usable without initialization
		p.deliver_lore.Append(100)
		p.deliver_lore.Append(200)
		if p.deliver_lore.Len() != 2 {
			t.Errorf("expected len 2, got %d", p.deliver_lore.Len())
		}
		if p.deliver_lore.Lookup(100) == -1 {
			t.Error("expected to find 100 in deliver_lore")
		}
		p.deliver_lore.Clear()
		if p.deliver_lore.Len() != 0 {
			t.Error("expected len 0 after clear")
		}
	})

	t.Run("char_magic.pledged_to_us", func(t *testing.T) {
		var m char_magic
		m.pledged_to_us.Append(1001)
		m.pledged_to_us.Append(1002)
		m.pledged_to_us.RemValueUniq(1001)
		if m.pledged_to_us.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.pledged_to_us.Len())
		}
		if m.pledged_to_us.Lookup(1002) == -1 {
			t.Error("expected to find 1002")
		}
	})

	t.Run("entity_misc.garr_watch", func(t *testing.T) {
		var m entity_misc
		// Simulate garrison adding watched units
		m.garr_watch.Append(5001)
		m.garr_watch.Append(5002)
		m.garr_watch.Append(5003)
		if m.garr_watch.Len() != 3 {
			t.Errorf("expected len 3, got %d", m.garr_watch.Len())
		}
		// Lookup should find watched unit
		if idx := m.garr_watch.Lookup(5002); idx == -1 {
			t.Error("expected to find 5002 in garr_watch")
		}
	})

	t.Run("entity_misc.garr_host", func(t *testing.T) {
		var m entity_misc
		m.garr_host.Append(6001)
		m.garr_host.Append(6002)
		// Lookup for hostile check
		if idx := m.garr_host.Lookup(6001); idx == -1 {
			t.Error("expected to find 6001 in garr_host")
		}
		if idx := m.garr_host.Lookup(9999); idx != -1 {
			t.Error("expected not to find 9999 in garr_host")
		}
	})
}

// Test persistent IList struct fields work with zero value (regression for Task 26.12.4)
func TestIListPersistentStructFields(t *testing.T) {
	t.Run("entity_player.units", func(t *testing.T) {
		var p entity_player
		p.units.Append(1001)
		p.units.Append(1002)
		p.units.Append(1003)
		if p.units.Len() != 3 {
			t.Errorf("expected len 3, got %d", p.units.Len())
		}
		if p.units.Lookup(1002) == -1 {
			t.Error("expected to find 1002 in units")
		}
		p.units.RemValueUniq(1002)
		if p.units.Len() != 2 {
			t.Errorf("expected len 2 after remove, got %d", p.units.Len())
		}
	})

	t.Run("entity_player.unformed", func(t *testing.T) {
		var p entity_player
		p.unformed.Append(2001)
		p.unformed.Append(2002)
		if p.unformed.Len() != 2 {
			t.Errorf("expected len 2, got %d", p.unformed.Len())
		}
		p.unformed.Clear()
		if p.unformed.Len() != 0 {
			t.Error("expected len 0 after clear")
		}
	})

	t.Run("entity_subloc.teaches", func(t *testing.T) {
		var s entity_subloc
		s.teaches.Append(9001)
		s.teaches.Append(9002)
		s.teaches.Append(9003)
		if s.teaches.Len() != 3 {
			t.Errorf("expected len 3, got %d", s.teaches.Len())
		}
		if s.teaches.Lookup(9002) != 1 {
			t.Error("expected to find 9002 at index 1")
		}
	})

	t.Run("entity_subloc.near_cities", func(t *testing.T) {
		var s entity_subloc
		s.near_cities.Append(10001)
		s.near_cities.Append(10002)
		if s.near_cities.Len() != 2 {
			t.Errorf("expected len 2, got %d", s.near_cities.Len())
		}
		if s.near_cities.Lookup(10001) == -1 {
			t.Error("expected to find 10001 in near_cities")
		}
	})

	t.Run("item_magic.may_use", func(t *testing.T) {
		var im item_magic
		im.may_use.Append(8001)
		im.may_use.Append(8002)
		if im.may_use.Len() != 2 {
			t.Errorf("expected len 2, got %d", im.may_use.Len())
		}
		if im.may_use.Lookup(8001) == -1 {
			t.Error("expected to find 8001 in may_use")
		}
	})

	t.Run("item_magic.may_study", func(t *testing.T) {
		var im item_magic
		im.may_study.Append(7001)
		im.may_study.Append(7002)
		im.may_study.Append(7003)
		if im.may_study.Len() != 3 {
			t.Errorf("expected len 3, got %d", im.may_study.Len())
		}
		im.may_study.RemValue(7002)
		if im.may_study.Len() != 2 {
			t.Errorf("expected len 2 after remove, got %d", im.may_study.Len())
		}
	})

	t.Run("att_ent.neutral", func(t *testing.T) {
		var a att_ent
		a.neutral.Append(3001)
		a.neutral.Append(3002)
		if a.neutral.Len() != 2 {
			t.Errorf("expected len 2, got %d", a.neutral.Len())
		}
		if a.neutral.Lookup(3001) == -1 {
			t.Error("expected to find 3001 in neutral")
		}
	})

	t.Run("att_ent.hostile", func(t *testing.T) {
		var a att_ent
		a.hostile.Append(4001)
		a.hostile.Append(4002)
		a.hostile.Append(4003)
		if a.hostile.Len() != 3 {
			t.Errorf("expected len 3, got %d", a.hostile.Len())
		}
		a.hostile.RemValueUniq(4002)
		if a.hostile.Len() != 2 {
			t.Errorf("expected len 2 after remove, got %d", a.hostile.Len())
		}
	})

	t.Run("att_ent.defend", func(t *testing.T) {
		var a att_ent
		a.defend.Append(5001)
		if a.defend.Len() != 1 {
			t.Errorf("expected len 1, got %d", a.defend.Len())
		}
		if a.defend.Lookup(5001) != 0 {
			t.Error("expected to find 5001 at index 0")
		}
	})

	t.Run("admit.l", func(t *testing.T) {
		var adm admit
		adm.l.Append(6001)
		adm.l.Append(6002)
		adm.l.Append(6003)
		if adm.l.Len() != 3 {
			t.Errorf("expected len 3, got %d", adm.l.Len())
		}
		if adm.l.Lookup(6002) != 1 {
			t.Error("expected to find 6002 at index 1")
		}
		adm.l.Clear()
		if adm.l.Len() != 0 {
			t.Error("expected len 0 after clear")
		}
	})
}

// Test List[string] to verify generic behavior
func TestListString(t *testing.T) {
	l := NewList("hello", "world")
	l.Append("!")
	l.Prepend("say")

	expected := []string{"say", "hello", "world", "!"}
	if !reflect.DeepEqual(l.Values(), expected) {
		t.Errorf("expected %v, got %v", expected, l.Values())
	}

	l.RemValueUniq("world")
	expected = []string{"say", "hello", "!"}
	if !reflect.DeepEqual(l.Values(), expected) {
		t.Errorf("after remove: expected %v, got %v", expected, l.Values())
	}

	idx := l.Lookup("hello")
	if idx != 1 {
		t.Errorf("expected index 1 for 'hello', got %d", idx)
	}
}
