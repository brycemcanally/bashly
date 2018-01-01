/*
Package cmds implements functionality for finding commands in bash scripts.
*/
package cmds

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/bryce/bashly/util"
)

type state int

// Command holds information about a command.
type Command struct {
	Name    string
	Options []string
}

// Finder holds data for finding command names in a bash script.
type Finder struct {
	pos, offset  int
	line         []rune
	cmds, delims util.Stack
	state        state
}

const (
	command      state = iota // we are processing a command string
	options                   // we are processing an option
	throwaway                 // we are processing runes that can be thrown away
	separator                 // we are processing a command separator string
	substitution              // we are processing a substitution string
	quote                     // we are processing a single quoted string ('...')
)

var (
	regexSeparator    = regexp.MustCompile(`^(\|[|&]?|&&|;)`)
	regexSubstitution = regexp.MustCompile(`^(\$\(|\)|\\*` + "`" + `)`)
)

// Find returns the name of the command being worked on at offset.
func Find(script string, offset int) *Command {
	// Empty script
	if len(script) <= 1 {
		return nil
	}

	f := &Finder{offset: offset, delims: util.Stack{}, cmds: util.Stack{}, state: command}

	// Find the line in which the offset is located
	re := regexp.MustCompile(`.*\n`)
	matches := re.FindAllString(script, -1)
	re2 := regexp.MustCompile(`\\\n`)
	line := ""
	for _, match := range matches {
		line += match
		if !re2.MatchString(match) || strings.Contains(match, "#") {
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
		return nil
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
	return cmd.(*Command)

}

func (f *Finder) command() {
	cmd := &Command{}

	for f.state == command || f.state == options {
		if f.pos >= len(f.line) {
			f.state = throwaway
			break
		}

		line := f.line[f.pos:]

		switch {
		case isSpace(line):
			if f.state == command && len(cmd.Name) > 0 {
				cmd.Options = append(cmd.Options, "")
				f.state = options
			} else if f.state == options && cmd.Options[len(cmd.Options)-1] != "" {
				cmd.Options = append(cmd.Options, "")
			}
			f.pos++
		case isMultiline(line):
			f.pos += 2
		case isComment(line):
			f.state = throwaway
		case isSeparator(line):
			f.state = separator
		case isSubstitution(line):
			f.state = substitution
		default:
			if f.state == command {
				cmd.Name += string(line[0])
				f.pos++
			} else if f.state == options {
				if cmd.Options[len(cmd.Options)-1] == "" && line[0] != '-' {
					cmd.Options = cmd.Options[:len(cmd.Options)-1]
					f.state = throwaway
				} else {
					cmd.Options[len(cmd.Options)-1] += string(line[0])
					f.pos++
				}
			}
		}
	}

	f.cmds.Push(cmd)
}

func (f *Finder) throwaway() {
	for f.state == throwaway {
		if f.pos >= f.offset {
			break
		}

		line := f.line[f.pos:]

		switch {
		case isComment(line):
			f.cmds.Push(&Command{"", nil})
			f.pos = f.offset
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
		f.cmds.Push(&Command{"", nil})
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
	return regexSeparator.MatchString(string(line))
}

func isSubstitution(line []rune) bool {
	return regexSubstitution.MatchString(string(line))
}

func isQuote(line []rune) bool {
	return line[0] == '\''
}

func getSeparator(line []rune) string {
	return regexSeparator.FindString(string(line))
}

func getSubstitution(line []rune) string {
	return regexSubstitution.FindString(string(line))
}

func isClosingSubstitution(delims util.Stack, token string) bool {
	val, err := delims.Top()
	open := val.(string)
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

	val, err := delims.Top()
	open := val.(string)
	switch {
	case token == "`":
		return err != nil || !strings.Contains(open, "`")
	case token == "\\`":
		return open == "`"
	default:
		val1, val2, err := delims.Top2()
		open1 := val1.(string)
		open2 := val2.(string)
		if err != nil {
			return false
		}

		backslashes := (len(open1) - len(open2)) * 2
		return len(token) == len(open1)+backslashes
	}
}
