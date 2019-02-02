package db

import "lmm/api/testing"

func TestMasks(tt *testing.T) {
	t := testing.NewTester(tt)

	t.Is(Masks(0), "")
	t.Is(Masks(1), "(?)")
	t.Is(Masks(2), "(?,?)")
	t.Is(Masks(3), "(?,?,?)")
	t.Is(Masks(4), "(?,?,?,?)")
	t.Is(Masks(5), "(?,?,?,?,?)")
}

func TestFieldMasks(tt *testing.T) {
	t := testing.NewTester(tt)

	t.Is(FieldMasks("dummy", 0), "")
	t.Is(FieldMasks("field", 1), "(field,?)")
	t.Is(FieldMasks("id", 2), "(id,?,?)")
	t.Is(FieldMasks("created_at", 3), "(created_at,?,?,?)")
	t.Is(FieldMasks("post_at", 4), "(post_at,?,?,?,?)")
	t.Is(FieldMasks("score", 5), "(score,?,?,?,?,?)")
}
