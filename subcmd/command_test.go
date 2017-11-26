package subcmd

import "testing"

func TestCommandName(t *testing.T) {
	cmd := Command{
		UsageLine: "subcommand [options]",
	}

	nm := "subcommand"

	if cmd.Name() != nm {
		t.Error("expected", nm, "got", cmd.Name())
	}
}

func TestCommandRunnable(t *testing.T) {
	cmd := Command{
		Run: func(cmd *Command, args []string) {
		},
	}

	if cmd.Runnable() != true {
		t.Error("expected command to be runnable")
	}

	cmd = Command{}

	if cmd.Runnable() != false {
		t.Error("expected command not to be runnable")
	}
}

func TestCommandRun(t *testing.T) {
	teststr := "not executed"
	cmd := Command{
		Run: func(cmd *Command, args []string) {
			teststr = "executed"
		},
	}
	var args []string

	cmd.Run(&cmd, args)

	if teststr != "executed" {
		t.Error("expected", "executed", "got", teststr)
	}
}
