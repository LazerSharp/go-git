package gogit

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Authoring struct {
	Name      string
	Email     string
	TimeStamp int64
	Zone      string
}

func (a *Authoring) String() string {
	return fmt.Sprintf("%s <%s> %d %s", a.Name, a.Email, a.TimeStamp, a.Zone)
}

func parseAuthoring(prefix string, line string) (*Authoring, error) {
	//fmt.Println("-->", prefix, "|", line)
	regx := prefix + ` (?P<Name>[^\<]+) \<(?P<Email>\w+@\w+\.\w+)\> (?P<TimeStamp>\d+) (?P<TZ>[+-]\d{4})$`
	m := matchGroup(line, regx)
	//fmt.Println(m)
	result := &Authoring{}
	result.Name = m["Name"]
	result.Email = m["Email"]
	timeStamp, err := strconv.ParseInt(m["TimeStamp"], 10, 64)
	if err != nil {
		return nil, err
	}
	result.TimeStamp = timeStamp
	result.Zone = m["TZ"]

	return result, nil
}

type Commit struct {
	Tree     *string
	Parent   *string
	Author   *Authoring
	Commiter *Authoring
	Comment  *string
}

func NewEmptyCommit() *Commit {
	return &Commit{}
}

func NewCommit(r io.Reader) (*Commit, error) {
	c := NewEmptyCommit()
	err := c.DeSerialize(r)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Commit) Type() string {
	return "commit"
}

func (c *Commit) Serialize(w io.Writer) error {
	fmt.Fprintf(w, "tree %s\n", *c.Tree)
	if c.Parent != nil {
		fmt.Fprintf(w, "parent %s\n", *c.Parent)
	}
	fmt.Fprintf(w, "author %s\n", c.Author)
	fmt.Fprintf(w, "committer %s\n", c.Commiter)
	fmt.Fprintln(w)
	fmt.Fprint(w, *c.Comment)
	return nil
}

func parseHah(line string) (*string, error) {
	splt := strings.Split(line, " ")
	if len(splt) != 2 {
		return nil, fmt.Errorf("error parsing commit tree")
	}
	return &splt[1], nil
}

func (c *Commit) parseParent(line string) error {
	hash, err := parseHah(line)
	if err != nil {
		return err
	}
	c.Parent = hash
	return nil
}
func (c *Commit) parseTree(line string) error {
	hash, err := parseHah(line)
	if err != nil {
		return err
	}
	c.Tree = hash
	return nil
}

func matchGroup(line string, regx string) map[string]string {
	m := map[string]string{}
	rx := regexp.MustCompile(regx)
	names := rx.SubexpNames()
	result := rx.FindStringSubmatch(line)
	for k, v := range result {
		m[names[k]] = v
	}
	return m

}

func (c *Commit) parseCommiter(line string) error {

	a, err := parseAuthoring("committer", line)
	if err != nil {
		return err
	}
	c.Commiter = a
	return nil
}

func (c *Commit) parseAuthor(line string) error {

	a, err := parseAuthoring("author", line)
	if err != nil {
		return err
	}
	c.Author = a
	return nil
}

func (c *Commit) DeSerialize(r io.Reader) error {
	reader := bufio.NewReader(r)
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		// done if it is a blank line
		if len(l) == 0 {
			break
		}
		line := string(l)
		//fmt.Println(line)
		switch {
		case strings.HasPrefix(line, "parent "):
			err = c.parseParent(line)
		case strings.HasPrefix(line, "tree "):
			err = c.parseTree(line)
		case strings.HasPrefix(line, "author "):
			err = c.parseAuthor(line)
		case strings.HasPrefix(line, "committer "):
			err = c.parseCommiter(line)
		}
		if err != nil {
			return err
		}
	}

	cmnt, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	v := string(cmnt)
	c.Comment = &v

	return nil
}
