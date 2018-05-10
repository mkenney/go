package cli

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	app := New(func() { fmt.Println("") })
	cmds := []*Cmd{
		&Cmd{
			Name: "cmd-name",
			Flags: map[string]*Flag{
				"cmd-flag-name": {
					Name:     "cmd-flag-name",
					Default:  "cmd-flag-value",
					Required: false,
				},
			},
		},
	}
	flags := []*Flag{{
		Name:     "app-flag-name",
		Default:  "app-flag-value",
		Required: false,
	}}
	app.PushCmds(cmds...)
	app.PushFlags(flags...)
	fmt.Printf("\n\n%s\n\n", app)
	Parse()
	t.Errorf("no")
}
