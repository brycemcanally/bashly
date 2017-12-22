/*
Package cmds implements functionality for finding commands in bash scripts.
*/
package cmds

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/bryce/bashly/cmds/util"
)

type state int

// Finder holds data for finding command names in a bash script.
type Finder struct {
	pos, offset  int
	line         []rune
	cmds, delims util.Stack
	state        state
}

const (
	command      state = iota // we are processing a command string
	throwaway                 // we are processing runes that can be thrown away
	separator                 // we are processing a command separator string
	substitution              // we are processing a substitution string
	quote                     // we are processing a single quoted string ('...')
)

// Find returns the name of the command being worked on at offset.
func Find(script string, offset int) string {
	// Empty script
	if len(script) <= 1 {
		return ""
	}

	f := &Finder{offset: offset, delims: util.Stack{}, cmds: util.Stack{}, state: command}

	// Find the line in which the offset is located
	re := regexp.MustCompile(`.*\n`)
	matches := re.FindAllString(script, -1)
	re2 := regexp.MustCompile(`[^#]*\\\n`)
	line := ""
	for _, match := range matches {
		line += match
		if !re2.MatchString(match) {
			if f.offset >= len(line) {
				f.offset -= len(line)
				line = ""
			} else {
				line = line[:len(line)-1]
				break
			}
		}
	}

	// Empty line
	if len(line) <= 0 {
		return ""
	}

	f.line = []rune(line)
	for f.pos < f.offset || f.state == command {

		switch f.state {
		case command:
			f.command()
		case throwaway:
			f.throwaway()
		case separator:
			f.separator()
		case substitution:
			f.substitution()
		case quote:
			f.quote()
		}
	}

	cmd, _ := f.cmds.Top()
	return cmd

}

// FindAll returns all of the commands in the script. (Stub)
func FindAll(script string) []string {
	return []string{}
}

func (f *Finder) command() {
	token := ""

done:
	for {
		if f.pos >= len(f.line) {
			f.state = throwaway
			break done
		}

		line := f.line[f.pos:]

		switch {
		case isSpace(line):
			if len(token) > 0 {
				f.state = throwaway
				break done
			} else {
				f.pos++
			}
		case isMultiline(line):
			f.pos += 2
		case isComment(line):
			f.state = throwaway
			break done
		case isSeparator(line):
			f.state = separator
			break done
		case isSubstitution(line):
			f.state = substitution
			break done
		default:
			token += string(line[0])
			f.pos++
		}
	}

	f.cmds.Push(token)
}

func (f *Finder) throwaway() {
	line := f.line[f.pos:]

	switch {
	case isComment(line):
		f.cmds.Push("")
		for f.pos < f.offset {
			f.pos++
		}
	case isSeparator(line):
		f.state = separator
	case isSubstitution(line):
		f.state = substitution
	case isQuote(line):
		f.state = quote
	default:
		f.pos++
	}
}

func (f *Finder) separator() {
	token := getSeparator(f.line[f.pos:])
	f.cmds.Pop()
	f.state = command
	f.consumeSeparatingToken(token)
}

func (f *Finder) substitution() {
	token := getSubstitution(f.line[f.pos:])

	switch {
	// Case for closing substitution
	case isClosingSubstitution(f.delims, token):
		f.cmds.Pop()
		f.delims.Pop()
		f.state = throwaway
		f.consumeSeparatingToken(token)
	// Case for opening parenthesis
	case isOpeningSubstitution(f.delims, token):
		f.delims.Push(token)
		f.state = command
		f.consumeSeparatingToken(token)
	default:
		f.state = throwaway
		f.pos += len(token)
	}
}

func (f *Finder) quote() {
	f.pos++
	if f.pos >= f.offset {
		f.state = throwaway
		return
	}

	r := f.line[f.pos]
	if r == '\'' {
		f.pos++
		f.state = throwaway
	}
}

func (f *Finder) consumeSeparatingToken(token string) {
	f.pos += len(token)
	// Offset is in the middle of a separating token between commands
	if f.pos > f.offset {
		f.cmds.Push("")
		f.state = throwaway
	}
}

func isSpace(line []rune) bool {
	return unicode.IsSpace(line[0])
}

func isMultiline(line []rune) bool {
	if len(line) <= 1 {
		return false
	}

	return string(line[:2]) == "\\\n"
}

func isComment(line []rune) bool {
	return line[0] == '#'
}

func isSeparator(line []rune) bool {
	re := regexp.MustCompile(`^(\|[|&]?|&&|;)`)
	return re.MatchString(string(line))
}

func isSubstitution(line []rune) bool {
	re := regexp.MustCompile(`^(\$\(|\)|\\*` + "`" + `)`)
	return re.MatchString(string(line))
}

func isQuote(line []rune) bool {
	return line[0] == '\''
}

func getSeparator(line []rune) string {
	re := regexp.MustCompile(`^(\|[|&]?|&&|;)`)
	return re.FindString(string(line))
}

func getSubstitution(line []rune) string {
	re := regexp.MustCompile(`^(\$\(|\)|\\*` + "`" + `)`)
	return re.FindString(string(line))
}

func isClosingSubstitution(delims util.Stack, token string) bool {
	open, err := delims.Top()
	if err != nil {
		return false
	}

	if open == "$(" {
		return token == ")"
	}

	return open == token
}

func isOpeningSubstitution(delims util.Stack, token string) bool {
	// Case for opening parenthesis
	if token == "$(" {
		return true
	}

	// Case for opening backticks
	if matched, _ := regexp.MatchString(`\\*`+"`", token); !matched {
		return false
	}

	open, err := delims.Top()
	switch {
	case token == "`":
		return err != nil || !strings.Contains(open, "`")
	case token == "\\`":
		return open == "`"
	default:
		opens, err := delims.Top2()
		if err != nil {
			return false
		}

		backslashes := (len(opens[1]) - len(opens[0])) * 2
		return len(token) == len(opens[1])+backslashes
	}
}