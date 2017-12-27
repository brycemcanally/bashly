package manual

import (
	"errors"
	"os/exec"
	"strconv"

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
func Get(command string, width int) (Page, error) {
	key := command + string(width)
	if val, ok := pageCache.Get(key); ok {
		return val.(Page), nil
	}

	man := exec.Command("/bin/bash", "-c", "export MANWIDTH="+strconv.Itoa(width)+"; man "+command)
	bytes, err := man.Output()
	if err != nil {
		return nil, errors.New("No manual page found")
	}

	page := Page(bytes)
	pageCache.Set(key, page)
	return page, nil
}
