package codeowners

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var instance *Codeowners
var once sync.Once

// Codeowners - patterns/owners mappings for the given repo
type Codeowners struct {
	repoRoot string
	Patterns []Codeowner
}

// Codeowner - owners for a given pattern
type Codeowner struct {
	Pattern string
	re      *regexp.Regexp
	Owners  []string
	Section string // gitlab code owners file has section
}

func (c Codeowner) String() string {
	return fmt.Sprintf("%s\t%v", c.Pattern, strings.Join(c.Owners, ", "))
}

func dirExists(fsys fs.FS, path string) (bool, error) {
	fi, err := fs.Stat(fsys, path)
	if err == nil && fi.IsDir() {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// findCodeownersFile - find a CODEOWNERS file somewhere within or below
// the working directory (wd), and open it.
func findCodeownersFile(fsys fs.FS, wd string) (io.Reader, string, error) {
	dir := wd
	for {
		for _, p := range []string{".", "docs", ".github", ".gitlab"} {
			pth := path.Join(dir, p)
			exists, err := dirExists(fsys, pth)
			if err != nil {
				return nil, "", err
			}
			if exists {
				f := path.Join(pth, "CODEOWNERS")
				_, err := fs.Stat(fsys, f)
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						continue
					}
					return nil, "", err
				}
				r, err := fsys.Open(f)
				return r, dir, err
			}
		}
		odir := dir
		dir = path.Dir(odir)
		// if we can't go up any further...
		if odir == dir {
			break
		}
		// if we're heading above the volume name (relevant on Windows)...
		if len(dir) < len(filepath.VolumeName(odir)) {
			break
		}
	}
	return nil, "", nil
}

func NewSingleCodeOwners(path string) (*Codeowners, error) {
	var err error
	once.Do(func() {
		instance, err = FromFile(path)
	})
	return instance, err
}

// FromFile creates a Codeowners from the path to a local file. Consider using
// [FromFileWithFS] instead.
func FromFile(path string) (*Codeowners, error) {
	base := "/"
	if filepath.IsAbs(path) && filepath.VolumeName(path) != "" {
		base = path[:len(filepath.VolumeName(path))+1]
	}
	path = path[len(base):]

	return FromFileWithFS(os.DirFS(base), path)
}

// FromFileWithFS creates a Codeowners from the path to a file relative to the
// given filesystem.
func FromFileWithFS(fsys fs.FS, path string) (*Codeowners, error) {
	r, root, err := findCodeownersFile(fsys, path)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("no CODEOWNERS found in %s", path)
	}
	return FromReader(r, root)
}

// FromReader creates a Codeowners from a given Reader instance and root path.
func FromReader(r io.Reader, repoRoot string) (*Codeowners, error) {
	co := &Codeowners{
		repoRoot: repoRoot,
	}
	co.Patterns = parseCodeowners(r)
	return co, nil
}

// parseCodeowners parses a list of Codeowners from a Reader
func parseCodeowners(r io.Reader) []Codeowner {
	co := []Codeowner{}
	s := bufio.NewScanner(r)
	curSection := ""
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) > 0 && strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) > 0 && strings.HasPrefix(fields[0], "[") && strings.HasSuffix(fields[len(fields)-1], "]") {
			curSection = fields[0]
			curSection = strings.TrimSuffix(strings.TrimPrefix(curSection, "["), "]")
			continue
		}
		if len(fields) > 1 {
			fields = combineEscapedSpaces(fields)
			c := NewCodeowner(fields[0], fields[1:], curSection)
			co = append(co, c)
		}
	}
	return co
}

// if any of the elements ends with a \, it was an escaped space
// put it back together properly so it's not treated as separate fields
func combineEscapedSpaces(fields []string) []string {
	outFields := make([]string, 0)
	escape := `\`
	for i := 0; i < len(fields); i++ {
		outField := fields[i]
		for strings.HasSuffix(fields[i], escape) && i+1 < len(fields) {
			outField = strings.Join([]string{strings.TrimRight(outField, escape), fields[i+1]}, " ")
			i++
		}
		outFields = append(outFields, outField)
	}

	return outFields
}

// NewCodeowner -
func NewCodeowner(pattern string, owners []string, section string) Codeowner {
	re := getPattern(pattern)
	c := Codeowner{
		Pattern: pattern,
		re:      re,
		Owners:  owners,
		Section: section,
	}
	return c
}

// Owners - return the list of code owners for the given path
// (within the repo root)
func (c *Codeowners) Owners(path string) []string {
	if strings.HasPrefix(path, c.repoRoot) {
		path = strings.Replace(path, c.repoRoot, "", 1)
	}

	// Order is important; the last matching pattern takes the most precedence.
	for i := len(c.Patterns) - 1; i >= 0; i-- {
		p := c.Patterns[i]

		if p.re.MatchString(path) {
			return p.Owners
		}
	}

	return nil
}

// Section - return the section of code owners for the given path
// (within the repo root)
func (c *Codeowners) Section(path string) string {
	if strings.HasPrefix(path, c.repoRoot) {
		path = strings.Replace(path, c.repoRoot, "", 1)
	}

	// Order is important; the last matching pattern takes the most precedence.
	for i := len(c.Patterns) - 1; i >= 0; i-- {
		p := c.Patterns[i]

		if p.re.MatchString(path) {
			return p.Section
		}
	}

	return ""
}

// based on github.com/sabhiram/go-gitignore
// but modified so that 'dir/*' only matches files in 'dir/'
func getPattern(line string) *regexp.Regexp {
	// when # or ! is escaped with a \
	if regexp.MustCompile(`^(\\#|\\!)`).MatchString(line) {
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.ReplaceAll(line, "?", `\?`)

	line = strings.ReplaceAll(line, magicStar, "*")

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern
}
