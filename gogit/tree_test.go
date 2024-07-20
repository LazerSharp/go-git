package gogit

import "testing"

func TestTreeType(t *testing.T) {
	var c GitObject = &Tree{}
	if c.Type() != "tree" {
		t.Fail()
	}
}

func TestDeserializeTree(t *testing.T) {
	repo, err := NewGitRepository("./testdata", false)
	if err != nil {
		t.Fail()
	}
	obj, err := ReadObject(repo, "26684dd98e763fe8ae01f6cebd3b0180d5e7a904")
	if err != nil {
		t.Fail()
	}
	AssertEq(t, obj.Type(), "tree")
	tree := obj.(*Tree)
	AssertEq(t, len(tree.Entries), 3)
	// code.c
	code := tree.Entries[0]
	AssertEq(t, code.Type, "blob")
	AssertEq(t, code.Path, "code.c")
	AssertEq(t, code.Hash, "37fc91b3ff6de8b20d7e643f0eb3cab29ae1fa68")
	AssertEq(t, code.Mode, "0644")

	AssertEq(t, tree.Entries[1].Type, "blob")

	// test dir
	tdir := tree.Entries[2]
	AssertEq(t, tdir.Type, "tree")
	AssertEq(t, tdir.Hash, "61918ebbb4dfa45f24a57b0f1d3ccf7a34961c3c")
	AssertEq(t, tdir.Mode, "0000")
	AssertEq(t, tdir.Path, "test")
}
