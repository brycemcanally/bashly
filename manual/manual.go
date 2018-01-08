package manual

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/bryce/bashly/cmds"
	"github.com/youtube/vitess/go/cache"
)

// Page holds data associated with a manual page
type Page []byte

// Size gets the size of the manual page as used in the cache.
func (p Page) Size() int {
	return 1
}

var pageCache = cache.NewLRUCache(10)

// Get returns the manual page for a given command.
func Get(command *cmds.Command, width int) (Page, error) {
	key := command.Name + string(width)
	if val, ok := pageCache.Get(key); ok {
		return val.(Page), nil
	}

	man := exec.Command("/bin/bash", "-c", "export MANWIDTH="+strconv.Itoa(width)+"; man "+command.Name)
	bytes, err := man.Output()
	if err != nil {
		return nil, errors.New("No manual page found")
	}

	page := Page(bytes)
	pageCache.Set(key, page)
	return page, nil
}

// GetOptions returns the sections of the manual page for a given command
// that have the description for the current options.
func GetOptions(command *cmds.Command, width int) (Page, error) {
	page, err := Get(command, width)
	if err != nil {
		return nil, err
	}
	optionsPage := []byte{}

	for _, opt := range command.Options {
		// Empty option (hanging - or --)
		if len(opt) <= 1 ||
			opt[1] == '-' && len(opt) <= 2 {
			continue
		}

		// Handle long option
		if opt[:2] == "--" {
			re, _ := regexp.Compile(`\n([ ]{7}([^ ].*?)?` + string(opt) + `.*?(\n|[ ]{8,}.*?\n)+)`)
			matches := re.FindAllSubmatch(page, -1)
			if len(matches) == 1 {
				optionsPage = append(optionsPage, matches[0][1]...)
			}
		} else {
			// Handle short option
			for i := 1; i < len(opt); i++ {
				re, _ := regexp.Compile(`\n([ ]{7}-` + string(opt[i]) + `.*?(\n|[ ]{8,}.*?\n)+)`)
				match := re.FindSubmatch(page)
				if match != nil {
					optionsPage = append(optionsPage, match[1]...)
				}
			}
		}
	}

	return optionsPage, nil
}
