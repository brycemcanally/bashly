package cmds

// Command holds information about a command.
type Command struct {
	name      string
	options   []string
	arguments []interface{}
	offset    []int
}
