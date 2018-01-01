package cmds

import "testing"

func TestMultiline(t *testing.T) {
	s := "ls|\\\n m\\\nv||echo \\\nhello\n"

	if cmd, _ := Find(s, 0); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}
	if cmd, _ := Find(s, 3); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 4); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 10); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 24); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
}

func TestComments(t *testing.T) {
	s := "# ls\nmv #; ls \\\ngrep\n"

	if cmd, err := Find(s, 0); err == nil {
		t.Error("Expected no command, got", cmd.Name)
	}
	if cmd, _ := Find(s, 5); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 8); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 16); cmd.Name != "grep" {
		t.Error("Expected grep, got", cmd.Name)
	}
}

func TestSeparators(t *testing.T) {
	s := "ls |mv|| grep|&chown &&pwd;\n"

	if cmd, _ := Find(s, 0); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}
	if cmd, _ := Find(s, 4); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, err := Find(s, 7); err == nil {
		t.Error("Expected no command, got", cmd.Name)
	}
	if cmd, _ := Find(s, 9); cmd.Name != "grep" {
		t.Error("Expected grep, got", cmd.Name)
	}
	if cmd, _ := Find(s, 13); cmd.Name != "grep" {
		t.Error("Expected grep, got", cmd.Name)
	}
	if cmd, err := Find(s, 14); err == nil {
		t.Error("Expected no command, got", cmd.Name)
	}
	if cmd, _ := Find(s, 16); cmd.Name != "chown" {
		t.Error("Expected chown, got", cmd.Name)
	}
	if cmd, _ := Find(s, 26); cmd.Name != "pwd" {
		t.Error("Expected pwd, got", cmd.Name)
	}
	if cmd, err := Find(s, 27); err == nil {
		t.Error("Expected no command, got", cmd.Name)
	}
}

func TestSubstitution(t *testing.T) {
	s := "echo `mv \\`ls \\\\\\`cd \\\\\\\\\\\\\\`gzip hello\\\\\\\\\\\\\\`\\\\\\`\\``\n"

	if cmd, _ := Find(s, 0); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 5); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 6); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 9); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, err := Find(s, 10); err == nil {
		t.Error("Expected no command, got", cmd.Name)
	}
	if cmd, _ := Find(s, 12); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}
	if cmd, _ := Find(s, 20); cmd.Name != "cd" {
		t.Error("Expected cd, got", cmd.Name)
	}
	if cmd, _ := Find(s, 36); cmd.Name != "gzip" {
		t.Error("Expected gzip, got", cmd.Name)
	}
	if cmd, _ := Find(s, 47); cmd.Name != "cd" {
		t.Error("Expected cd, got", cmd.Name)
	}
	if cmd, _ := Find(s, 54); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}

	s = "echo $(ls $(mv `grep`))\n"
	if cmd, _ := Find(s, 0); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 7); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}
	if cmd, _ := Find(s, 10); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}
	if cmd, _ := Find(s, 12); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 16); cmd.Name != "grep" {
		t.Error("Expected grep, got", cmd.Name)
	}
	if cmd, _ := Find(s, 21); cmd.Name != "mv" {
		t.Error("Expected mv, got", cmd.Name)
	}
	if cmd, _ := Find(s, 22); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}

	s = "echo ) \\` | ls\n"
	if cmd, _ := Find(s, 10); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 11); cmd.Name != "ls" {
		t.Error("Expected ls, got", cmd.Name)
	}
}

func TestQuotes(t *testing.T) {
	s := "echo 'hello `grep` $(pwd) \\ #'\n"

	if cmd, _ := Find(s, 0); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 13); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 21); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 28); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
	if cmd, _ := Find(s, 30); cmd.Name != "echo" {
		t.Error("Expected echo, got", cmd.Name)
	}
}
