package gogit

import (
	"strings"
	"testing"
)

func TestCommitType(t *testing.T) {
	var c GitObject = NewEmptyCommit()
	if c.Type() != "commit" {
		t.Fail()
	}
}

var commitContent = `tree 7c550fd1db2ce0398f84b9c803bb2a2aa94b2dfe
author LazerSharp <email2barun@gmail.com> 1721064235 +0530
committer Barun Halder <email2barun@gmail.com> 1721064235 +0530

Initial Commit!`

func TestDeserialize(t *testing.T) {
	c := NewEmptyCommit()
	err := c.DeSerialize(strings.NewReader(commitContent))
	if err != nil {
		t.Fail()
	}
	AssertEq(t, *c.Tree, "7c550fd1db2ce0398f84b9c803bb2a2aa94b2dfe")
	AssertEq(t, c.Author.Name, "LazerSharp")
	AssertEq(t, c.Author.Email, "email2barun@gmail.com")
	AssertEq(t, c.Author.TimeStamp, int64(1721064235))
	AssertEq(t, c.Commiter.TimeStamp, int64(1721064235))
	AssertEq(t, c.Commiter.Zone, "+0530")
	AssertEq(t, c.Commiter.Name, "Barun Halder")
}
