package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type FlagType int

const (
	Bool FlagType = iota
	Float
	Int
	JSON
	String
	Time
)

/*
New returns a pointer to a struct representing a CLI application.
*/
func New(usage Usage) *Cmd {
	args := os.Args
	args = append([]string{path.Dir(args[0]), path.Base(args[0])}, args[1:]...)
	name := args[1]
	path := args[0]
	args = args[2:]

	return &Cmd{
		Args:   make(map[string]*Arg),
		Called: true,
		Cmds:   make(map[string]*Cmd),
		Flags:  make(map[string]*Flag),
		Name:   name,
		Path:   path,
		Usage:  usage,
	}
}

/*
PushCmds adds command definitions to the application.
*/
func (cmd *Cmd) PushCmds(cmds ...*Cmd) {
	for _, c := range cmds {
		cmd.Cmds[c.Name] = c
	}
}

/*
PushFlags adds flag definitions to the application.
*/
func (cmd *Cmd) PushFlags(flags ...*Flag) {
	for _, flag := range flags {
		cmd.Flags[flag.Name] = flag
	}
}

/*
Execute parses the CLI arguments and compares them to the defined
commands, args, and flags.
*/
func (cmd *Cmd) Execute() error {

	// break the path and the command into separate tokens
	args := os.Args
	args = append([]string{path.Dir(args[0]), path.Base(args[0])}, args[1:]...)
	cmd.Path = args[0]
	cmd.Name = args[1]
	args = args[2:]

	type token struct {
		typ   string
		name  string
		value string
	}
	shift := func(args []string, cmd *Cmd) (*token, []string) {
		var tkn *token
		arg := args[0]
		args = args[1:]
		runes := []rune(arg)
		if "---" != string(runes[0:3]) && ("-" == string(runes[0:1]) || "--" == string(runes[0:2])) {
			parts := strings.Split(strings.Trim(arg, "-"), "=")
			flag, ok := cmd.Flags[parts[0]]
			if !ok {
				// no matching flag definition
				return nil, args
			}

			tkn.typ = "flag"
			tkn.name = parts[0]

			if 2 >= len(parts) {
				tkn.value = strings.Join(parts[1:], "=")
			} else {
				switch flag.Default.(type) {
				case bool:
					tkn.value = "1"
				default:
					if len(args) > 0 {
						tkn.value = args[0]
						args = args[1:]
					}
				}
			}
		} else if _, ok := cmd.Cmds[arg]; ok {
			tkn.name = arg
			tkn.typ = "cmd"

		} else if _, ok := cmd.Args[arg]; ok {
			tkn.name = arg
			tkn.typ = "arg"

		} else {
			return nil, args
		}
		return tkn, args
	}

	var tkn *token
	path := cmd
	for {
		tkn, args = shift(args, path)
		if nil == tkn {
			// error, unknown token
			continue
		}
		if "flag" == tkn.typ {
			if _, ok := path.Flags[tkn.name]; ok {
				path.Cmds[tkn.name].Called = true
				path.Flags[tkn.name].Set(tkn.value)
			}
		} else if "cmd" == tkn.typ {
			if _, ok := path.Cmds[tkn.name]; ok {
				path.Cmds[tkn.name].Called = true
				path = path.Cmds[tkn.name]
			}
		} else if "arg" == tkn.typ {
			if _, ok := path.Args[tkn.name]; ok {
				path.Args[tkn.name].Called = true
			}
		} else {
			return fmt.Errorf("Unknown type '%v'", tkn.typ)
		}

		if 0 == len(args) {
			break
		}
	}

	//for k, v := range args {
	//
	//}

	return nil
}

type Cmd struct {
	Args   map[string]*Arg
	Called bool
	Cmds   map[string]*Cmd
	Flags  map[string]*Flag
	Name   string
	Path   string
	Usage  Usage
}

/*
String implements stringer
*/
func (cmd *Cmd) String() string {
	str := ""
	if "" != cmd.Path {
		str = fmt.Sprintf("%s/%s", cmd.Path, cmd.Name)
	} else {
		str = cmd.Name
	}
	for _, cmd := range cmd.Cmds {
		str = fmt.Sprintf("%s %s", str, cmd)
	}
	for _, arg := range cmd.Args {
		str = fmt.Sprintf("%s %s", str, arg)
	}
	for _, flag := range cmd.Flags {
		str = fmt.Sprintf("%s %s", str, flag)
	}
	return str
}

type Arg struct {
	Called bool
	Name   string
	Flags  []*Flag
}

func (arg Arg) String() string {
	str := arg.Name
	for _, flag := range arg.Flags {
		str = fmt.Sprintf("%s %s", str, flag)
	}
	return str
}

type Flag struct {
	Called   bool
	Default  interface{}
	Name     string
	Required bool
	Value    interface{}
}

func (f Flag) String() string {
	str := "--" + f.Name

	val := f.Default
	if nil != f.Value {
		val = f.Value
	}

	switch f.Default.(type) {
	case bool:
		str = fmt.Sprintf("%s=%v", str, val.(bool))
	case []byte:
		str = fmt.Sprintf("%s=%s", str, string(val.([]byte)))
	case float32:
		str = fmt.Sprintf("%s=%f", str, val.(float32))
	case float64:
		str = fmt.Sprintf("%s=%f", str, val.(float64))
	case int:
		str = fmt.Sprintf("%s=%d", str, val.(int))
	case int8:
		str = fmt.Sprintf("%s=%d", str, val.(int8))
	case int16:
		str = fmt.Sprintf("%s=%d", str, val.(int16))
	case int32:
		str = fmt.Sprintf("%s=%d", str, val.(int32))
	case int64:
		str = fmt.Sprintf("%s=%d", str, val.(int64))
	case string:
		str = fmt.Sprintf("%s=%s", str, val.(string))
	}
	return str
}

func (f Flag) Set(val string) error {
	var err error

	switch f.Default.(type) {
	case bool:
		switch val {
		case "":
			f.Value = false
		case "false":
			f.Value = false
		case "0":
			f.Value = false
		case "1":
			f.Value = true
		case "true":
			f.Value = true
		default:
			return fmt.Errorf("invalid value for boolean argument: '%s'", val)
		}
	case float32:
		f.Value, err = strconv.ParseFloat(val, 10)
		if nil != err {
			return err
		}
	case int:
		f.Value, err = strconv.ParseInt(val, 10, 64)
		if nil != err {
			return err
		}
	case []byte:
		err = json.Unmarshal([]byte(val), &f.Value)
		if nil != err {
			return err
		}
	case string:
		f.Value = val
	default:
		f.Value = val
	}
	return nil
}
