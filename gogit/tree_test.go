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
	AssertEq(t, code.Type.String(), "blob")
	AssertEq(t, code.Path, "code.c")
	AssertEq(t, code.Hash, "37fc91b3ff6de8b20d7e643f0eb3cab29ae1fa68")
	AssertEq(t, code.Mode, "0644")

	AssertEq(t, tree.Entries[1].Type, BlobEntry)

	// test dir
	tdir := tree.Entries[2]
	AssertEq(t, tdir.Type.String(), "tree")
	AssertEq(t, tdir.Hash, "61918ebbb4dfa45f24a57b0f1d3ccf7a34961c3c")
	AssertEq(t, tdir.Mode, "0000")
	AssertEq(t, tdir.Path, "test")
}
func TestEntryType(t *testing.T) {
	AssertEq(t, BlobEntry.Val(), "10")
	AssertEq(t, TreeEntry.Val(), "4")
}

// $ git ls-tree 26684dd98e763fe8ae01f6cebd3b0180d5e7a904
// 100644 blob 37fc91b3ff6de8b20d7e643f0eb3cab29ae1fa68    code.c
// 100644 blob fac692d6a859684017135fbd4ab6f767cc185e7b    test.c
// 040000 tree 61918ebbb4dfa45f24a57b0f1d3ccf7a34961c3c    test
func TestSerializeTree(t *testing.T) {

	tree := &Tree{
		Entries: []Entry{
			{
				Path: "code.c",
				Hash: "37fc91b3ff6de8b20d7e643f0eb3cab29ae1fa68",
				Mode: "0644",
				Type: BlobEntry,
			},
			{
				Path: "test.c",
				Hash: "fac692d6a859684017135fbd4ab6f767cc185e7b",
				Mode: "0644",
				Type: BlobEntry,
			},
			{
				Path: "test",
				Hash: "61918ebbb4dfa45f24a57b0f1d3ccf7a34961c3c",
				Mode: "0000",
				Type: TreeEntry,
			},
		},
	}
	hash, err := WriteObject(tree, nil)
	if err != nil {
		t.Fatal(err)
	}
	AssertEq(t, hash, "26684dd98e763fe8ae01f6cebd3b0180d5e7a904")

}
