package sqlx

import (
	"testing"
)

func TestInNormalSlice(t *testing.T) {
	query, args, err := In("SELECT * FROM t WHERE id IN (?)", []int{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	expected := "SELECT * FROM t WHERE id IN (?,?,?)"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d", len(args))
	}
}

func TestInEmptySlice(t *testing.T) {
	query, args, err := In("SELECT * FROM t WHERE id IN (?)", []int{})
	if err != nil {
		t.Fatal(err)
	}
	expected := "SELECT * FROM t WHERE id IN (NULL)"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestInMultiplePlaceholdersOneEmpty(t *testing.T) {
	query, args, err := In("SELECT * FROM t WHERE id IN (?) AND name IN (?)", []int{1, 2}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	expected := "SELECT * FROM t WHERE id IN (?,?) AND name IN (NULL)"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

func TestInNonSliceArg(t *testing.T) {
	query, args, err := In("SELECT * FROM t WHERE id = ?", 42)
	if err != nil {
		t.Fatal(err)
	}
	expected := "SELECT * FROM t WHERE id = ?"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
	if len(args) != 1 || args[0] != 42 {
		t.Errorf("expected [42], got %v", args)
	}
}

func TestRebind(t *testing.T) {
	query := Rebind("SELECT * FROM t WHERE id = ? AND name = ?")
	expected := "SELECT * FROM t WHERE id = $1 AND name = $2"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
}

func TestRebindNoPlaceholders(t *testing.T) {
	query := Rebind("SELECT * FROM t")
	expected := "SELECT * FROM t"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
}
