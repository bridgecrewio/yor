package codeowners

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var (
	sample = `# comment
*	@everyone

   foobar/  someone@else.com

docs/**	@org/docteam @joe`
	sample2 = `* @hairyhenderson`
	sample3 = `baz/* @baz @qux`
	sample4 = `[test]
*   @everyone
[test2]
*/foo @everyoneelse`

	// based on https://help.github.com/en/github/creating-cloning-and-archiving-repositories/about-code-owners#codeowners-syntax
	// with a few unimportant modifications
	fullSample = `# This is a comment.
# Each line is a file pattern followed by one or more owners.

# These owners will be the default owners for everything in
# the repo. Unless a later match takes precedence,
# @global-owner1 and @global-owner2 will be requested for
# review when someone opens a pull request.
*       @global-owner1 @global-owner2

# Order is important; the last matching pattern takes the most
# precedence. When someone opens a pull request that only
# modifies JS files, only @js-owner and not the global
# owner(s) will be requested for a review.
*.js	@js-owner

# You can also use email addresses if you prefer. They'll be
# used to look up users just like we do for commit author
# emails.
*.go docs@example.com

# In this example, @doctocat owns any files in the build/logs
# directory at the root of the repository and any of its
# subdirectories.
/build/logs/ @doctocat

# The 'docs/*' pattern will match files like
# 'docs/getting-started.md' but not further nested files like
# 'docs/build-app/troubleshooting.md'.
docs/*  docs@example.com

# In this example, @octocat owns any file in an apps directory
# anywhere in your repository.
apps/ @octocat

# In this example, @doctocat owns any file in the '/docs'
# directory in the root of your repository.
/docs/ @doctocat

  foobar/ @fooowner

\#foo/ @hashowner

docs/*.md @mdowner

# this example tests an escaped space in the path
space/test\ space/ @spaceowner
`

	codeowners []Codeowner
)

func TestParseCodeowners(t *testing.T) {
	t.Parallel()
	r := bytes.NewBufferString(sample)
	c := parseCodeowners(r)
	expected := []Codeowner{
		co("*", []string{"@everyone"}, ""),
		co("foobar/", []string{"someone@else.com"}, ""),
		co("docs/**", []string{"@org/docteam", "@joe"}, ""),
	}
	assert.Equal(t, expected, c)
}

func TestParseCodeownersSections(t *testing.T) {
	t.Parallel()
	r := bytes.NewBufferString(sample4)
	c := parseCodeowners(r)
	expected := []Codeowner{
		co("*", []string{"@everyone"}, "test"),
		co("*/foo", []string{"@everyoneelse"}, "test2"),
	}
	assert.Equal(t, expected, c)
}

func BenchmarkParseCodeowners(b *testing.B) {
	r := bytes.NewBufferString(sample)
	var c []Codeowner

	for n := 0; n < b.N; n++ {
		c = parseCodeowners(r)
	}

	codeowners = c
}

func TestFindCodeownersFile(t *testing.T) {
	fsys := fstest.MapFS{
		"src/.github/CODEOWNERS":      &fstest.MapFile{Data: []byte(sample)},
		"src/foo/CODEOWNERS":          &fstest.MapFile{Data: []byte(sample2)},
		"src/foo/qux/docs/CODEOWNERS": &fstest.MapFile{Data: []byte(sample3)},
	}

	r, root, err := findCodeownersFile(fsys, "src")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src", root)

	b, _ := io.ReadAll(r)
	assert.Equal(t, sample, string(b))

	r, root, err = findCodeownersFile(fsys, "src/foo/bar")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src/foo", root)

	b, _ = io.ReadAll(r)
	assert.Equal(t, sample2, string(b))

	r, root, err = findCodeownersFile(fsys, "src/foo/qux/quux")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src/foo/qux", root)

	b, _ = io.ReadAll(r)
	assert.Equal(t, sample3, string(b))

	r, _, err = findCodeownersFile(fsys, ".")
	require.NoError(t, err)
	assert.Nil(t, r)
}

func co(pattern string, owners []string, section string) Codeowner {
	c := NewCodeowner(pattern, owners, section)
	return c
}

func TestFullParseCodeowners(t *testing.T) {
	t.Parallel()

	c := parseCodeowners(strings.NewReader(fullSample))
	codeowners := &Codeowners{
		repoRoot: "/build",
		Patterns: c,
	}

	// these tests were ported from https://github.com/softprops/codeowners
	data := []struct {
		path   string
		owners []string
	}{
		{"#foo/bar.go", []string{"@hashowner"}},
		{"foobar/baz.go", []string{"@fooowner"}},
		{"/docs/README.md", []string{"@mdowner"}},
		// XXX: uncertain about this one
		{"blah/docs/README.md", []string{"docs@example.com"}},
		{"foo.txt", []string{"@global-owner1", "@global-owner2"}},
		{"foo/bar.txt", []string{"@global-owner1", "@global-owner2"}},
		{"foo.js", []string{"@js-owner"}},
		{"foo/bar.js", []string{"@js-owner"}},
		{"foo.go", []string{"docs@example.com"}},
		{"foo/bar.go", []string{"docs@example.com"}},
		// relative to root
		{"build/logs/foo.go", []string{"@doctocat"}},
		{"build/logs/foo/bar.go", []string{"@doctocat"}},
		// not relative to root
		{"foo/build/logs/foo.go", []string{"docs@example.com"}},
		// docs anywhere
		{"foo/docs/foo.js", []string{"docs@example.com"}},
		{"foo/bar/docs/foo.js", []string{"docs@example.com"}},
		// but not nested
		{"foo/bar/docs/foo/foo.js", []string{"@js-owner"}},
		{"foo/apps/foo.js", []string{"@octocat"}},
		{"docs/foo.js", []string{"@doctocat"}},
		{"/docs/foo.js", []string{"@doctocat"}},
		{"/space/test space/doc1.txt", []string{"@spaceowner"}},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%q==%#v", d.path, d.owners), func(t *testing.T) {
			assert.EqualValues(t, d.owners, codeowners.Owners(d.path))
		})
	}
}

func TestOwners(t *testing.T) {
	foo := []string{"@foo"}
	bar := []string{"@bar"}
	baz := []string{"@baz"}
	data := []struct {
		patterns []Codeowner
		path     string
		expected []string
	}{
		{[]Codeowner{co("a/*", foo, "")}, "c/b", nil},
		{[]Codeowner{co("**", foo, "")}, "a/b", foo},
		{[]Codeowner{co("**", foo, ""), co("a/b/*", bar, "")}, "a/b/c", bar},
		{[]Codeowner{co("**", foo, ""), co("a/b/*", bar, ""), co("a/b/c", baz, "")}, "a/b/c", baz},
		{[]Codeowner{co("**", foo, ""), co("a/*/c", bar, ""), co("a/b/*", baz, "")}, "a/b/c", baz},
		{[]Codeowner{co("**", foo, ""), co("a/b/*", bar, ""), co("a/b/", baz, "")}, "a/b/bar", baz},
		{[]Codeowner{co("**", foo, ""), co("a/b/*", bar, ""), co("a/b/", baz, "")}, "/someroot/a/b/bar", baz},
		{[]Codeowner{
			co("*", foo, ""),
			co("/a/*", bar, ""),
			co("/b/**", baz, "")}, "/a/aa/file", foo},
		{[]Codeowner{
			co("*", foo, ""),
			co("/a/**", bar, "")}, "/a/bb/file", bar},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%s==%s", d.path, d.expected), func(t *testing.T) {
			c := &Codeowners{Patterns: d.patterns, repoRoot: "/someroot"}
			owners := c.Owners(d.path)
			assert.Equal(t, d.expected, owners)
		})
	}
}

func TestCombineEscapedSpaces(t *testing.T) {
	data := []struct {
		fields   []string
		expected []string
	}{
		{[]string{"docs/", "@owner"}, []string{"docs/", "@owner"}},
		{[]string{"docs/bob/**", "@owner"}, []string{"docs/bob/**", "@owner"}},
		{[]string{"docs/bob\\", "test/", "@owner"}, []string{"docs/bob test/", "@owner"}},
		{[]string{"docs/bob\\", "test/sub/final\\", "space/", "@owner"}, []string{"docs/bob test/sub/final space/", "@owner"}},
		{[]string{"docs/bob\\", "test/another\\", "test/**", "@owner"}, []string{"docs/bob test/another test/**", "@owner"}},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%s==%s", d.fields, d.expected), func(t *testing.T) {
			assert.Equal(t, d.expected, combineEscapedSpaces(d.fields))
		})
	}
}

func cwd() string {
	_, filename, _, _ := runtime.Caller(0)
	cwd := path.Dir(filename)
	return cwd
}

func ExampleFromFile() {
	tpath, _ := filepath.Abs(filepath.Dir(filepath.Dir(cwd())))
	tpath += "/tests"
	c, _ := FromFile(tpath)
	fmt.Println(c.Patterns[0])
	// Output:
	// *	bridgecrewio
}

func ExampleFromFileWithFS() {
	// open filesystem rooted at current working directory
	fsys := os.DirFS(cwd())

	c, _ := FromFileWithFS(fsys, ".")
	fmt.Println(c.Patterns[0])
	// Output:
	// *	bridgecrewio
}

func ExampleFromReader() {
	reader := strings.NewReader(sample2)
	c, _ := FromReader(reader, "")
	fmt.Println(c.Patterns[0])
	// Output:
	// *	@hairyhenderson
}

func ExampleCodeowners_Owners() {
	c, _ := FromFile(cwd())
	owners := c.Owners("README.md")
	for i, o := range owners {
		fmt.Printf("Owner #%d is %s\n", i, o)
	}
	// Output:
	// Owner #0 is @hairyhenderson
}
