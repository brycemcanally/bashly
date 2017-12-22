package cmds

import "testing"

func TestMultiline(t *testing.T) {
	s := "ls|\\\n m\\\nv||echo \\\nhello\n"
	if cmd := Find(s, 0); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}
	if cmd := Find(s, 3); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 4); cmd != "mv" {
		t.Error("Expected empty string, got", cmd)
	}
	if cmd := Find(s, 10); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 24); cmd != "echo" {
		t.Error("Expected echo, got", cmd)
	}
}

func TestSeparators(t *testing.T) {
	s := "ls |mv|| grep|&chown &&pwd;\n"

	if cmd := Find(s, 0); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}
	if cmd := Find(s, 4); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 7); cmd != "" {
		t.Error("Expected empty string, got", cmd)
	}
	if cmd := Find(s, 9); cmd != "grep" {
		t.Error("Expected grep, got", cmd)
	}
	if cmd := Find(s, 13); cmd != "grep" {
		t.Error("Expected grep, got", cmd)
	}
	if cmd := Find(s, 14); cmd != "" {
		t.Error("Expected empty string, got", cmd)
	}
	if cmd := Find(s, 16); cmd != "chown" {
		t.Error("Expected chown, got", cmd)
	}
	if cmd := Find(s, 26); cmd != "pwd" {
		t.Error("Expected pwd, got", cmd)
	}
	if cmd := Find(s, 27); cmd != "" {
		t.Error("Expected empty string, got", cmd)
	}
}

func TestSubstitution(t *testing.T) {
	s := "echo `mv \\`ls \\\\\\`cd \\\\\\\\\\\\\\`gzip hello\\\\\\\\\\\\\\`\\\\\\`\\``\n"

	if cmd := Find(s, 0); cmd != "echo" {
		t.Error("Expected echo, got", cmd)
	}
	if cmd := Find(s, 5); cmd != "echo" {
		t.Error("Expected echo, got", cmd)
	}
	if cmd := Find(s, 6); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 9); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 10); cmd != "" {
		t.Error("Expected empty string, got", cmd)
	}
	if cmd := Find(s, 12); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}
	if cmd := Find(s, 20); cmd != "cd" {
		t.Error("Expected cd, got", cmd)
	}
	if cmd := Find(s, 36); cmd != "gzip" {
		t.Error("Expected gzip, got", cmd)
	}
	if cmd := Find(s, 47); cmd != "cd" {
		t.Error("Expected cd, got", cmd)
	}
	if cmd := Find(s, 54); cmd != "echo" {
		t.Error("Expected echo, got", cmd)
	}

	s = "echo $(ls $(mv `grep`))\n"
	if cmd := Find(s, 0); cmd != "echo" {
		t.Error("Expected echo, got", cmd)
	}
	if cmd := Find(s, 7); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}
	if cmd := Find(s, 10); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}
	if cmd := Find(s, 12); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 16); cmd != "grep" {
		t.Error("Expected grep, got", cmd)
	}
	if cmd := Find(s, 21); cmd != "mv" {
		t.Error("Expected mv, got", cmd)
	}
	if cmd := Find(s, 22); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}

	s = "echo ) \\` | ls\n"
	if cmd := Find(s, 10); cmd != "echo" {
		t.Error("Expected echo, got", cmd)
	}
	if cmd := Find(s, 11); cmd != "ls" {
		t.Error("Expected ls, got", cmd)
	}

}
