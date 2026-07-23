package sqlx

import (
	"errors"
	"testing"
)

func TestInWithEmptySlice(t *testing.T) {
	query := "SELECT * FROM users WHERE id IN (?) AND status IN (?)"
	ids := []int{1, 2, 3}
	statuses := []string{} // Empty slice

	_, _, err := In(query, ids, statuses)
	if err == nil {
		t.Error("Expected an error when passing an empty slice to In, but got nil")
	}
	if !errors.Is(err, ErrEmptySlice) {
		t.Errorf("Expected ErrEmptySlice, got '%s'", err)
	}
}

func TestInWithMultipleSlices(t *testing.T) {
	query := "SELECT * FROM users WHERE id IN (?) AND status IN (?)"
	ids := []int{1, 2, 3}
	statuses := []string{"active", "pending"}

	result, args, err := In(query, ids, statuses)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectedQuery := "SELECT * FROM users WHERE id IN (?,?,?) AND status IN (?,?)"
	if result != expectedQuery {
		t.Errorf("expected '%s', got '%s'", expectedQuery, result)
	}
	if len(args) != 5 {
		t.Errorf("expected 5 expanded args, got %d", len(args))
	}
}

func TestInWithMixedArgs(t *testing.T) {
	query := "SELECT * FROM users WHERE name = ? AND id IN (?) AND status = ?"
	name := "alice"
	ids := []int{1, 2}
	status := "active"

	result, args, err := In(query, name, ids, status)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectedQuery := "SELECT * FROM users WHERE name = ? AND id IN (?,?) AND status = ?"
	if result != expectedQuery {
		t.Errorf("expected '%s', got '%s'", expectedQuery, result)
	}
	if len(args) != 4 {
		t.Errorf("expected 4 expanded args, got %d", len(args))
	}
}

func TestInWithSingleSlice(t *testing.T) {
	query := "SELECT * FROM users WHERE id IN (?)"
	ids := []int{10, 20, 30, 40}

	result, args, err := In(query, ids)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectedQuery := "SELECT * FROM users WHERE id IN (?,?,?,?)"
	if result != expectedQuery {
		t.Errorf("expected '%s', got '%s'", expectedQuery, result)
	}
	if len(args) != 4 {
		t.Errorf("expected 4 expanded args, got %d", len(args))
	}
}

func TestInWithEmptySliceFirst(t *testing.T) {
	query := "SELECT * FROM users WHERE id IN (?) AND status IN (?)"
	ids := []int{} // Empty
	statuses := []string{"active"}

	_, _, err := In(query, ids, statuses)
	if err == nil {
		t.Error("Expected error for empty first slice")
	}
	if !errors.Is(err, ErrEmptySlice) {
		t.Errorf("Expected ErrEmptySlice, got '%s'", err)
	}
}

func TestInWithOnlyScalars(t *testing.T) {
	query := "SELECT * FROM users WHERE name = ? AND age = ?"
	result, args, err := In(query, "alice", 30)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if result != query {
		t.Errorf("expected unchanged query, got '%s'", result)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

func TestInWithArray(t *testing.T) {
	query := "SELECT * FROM users WHERE id IN (?)"
	ids := [3]int{1, 2, 3}

	result, args, err := In(query, ids)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectedQuery := "SELECT * FROM users WHERE id IN (?,?,?)"
	if result != expectedQuery {
		t.Errorf("expected '%s', got '%s'", expectedQuery, result)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d", len(args))
	}
}

func TestInWithEmptyArray(t *testing.T) {
	query := "SELECT * FROM users WHERE id IN (?)"
	ids := [0]int{}

	_, _, err := In(query, ids)
	if err == nil {
		t.Error("Expected error for empty array")
	}
	if !errors.Is(err, ErrEmptySlice) {
		t.Errorf("Expected ErrEmptySlice, got '%s'", err)
	}
}
