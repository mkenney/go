package cli

import (
	"log"
	"os"
)

func New1(arguments ...[]string) (*Cli, error) {
	return Parse(arguments...)
}

type Cli struct {
	cmd   string
	path  string
	Flags []*Flag
}

func (cli *Cli) Push(flag *Flag) {
	cli.Flags = append(cli.Flags, flag)
}

/*
Cli is an interface
*/
type CliIface interface {
	Cmd() string
	Path() string
	Usage()
	Workdir() string
}

/*
Parse parses the os.Args parameters by default. Optionally a []string
value may be passed.
*/
func Parse(arguments ...[]string) (*Cli, error) {
	return nil, nil
	//	args := os.Args
	//	if len(arguments) > 0 {
	//		args = arguments[0]
	//	}
	//
	//	// break the path and the command into separate tokens
	//	args = append([]string{path.Dir(args[0]), path.Base(args[0])}, args[1:]...)
	//
	//	cli := &Cli{}
	//	flag := &Flag{}
	//	var isLongFlag bool
	//	var isShortFlag bool
	//	var isArg bool
	//	var wasLongFlag bool
	//	var wasShortFlag bool
	//	var wasArg bool
	//	for idx, token := range args {
	//		wasLongFlag = isLongFlag
	//		wasShortFlag = isShortFlag
	//		wasArg = isArg
	//
	//		runes := []rune(token)
	//		var isArg bool
	//		if "---" == string(runes[0:3]) || "-" != string(runes[0:1]) {
	//			isArg = true
	//		}
	//
	//		var isLongFlag bool
	//		if "--" == string(runes[0:2]) && !isArg {
	//			isLongFlag = true
	//			runes = runes[2:]
	//		}
	//
	//		var isShortFlag bool
	//		if "-" == string(runes[0:]) && !isLongFlag {
	//			isShortFlag = true
	//			runes = runes[1:]
	//		}
	//
	//		if wasArg && isArg {
	//
	//		}
	//
	//		if isLongFlag || isShortFlag {
	//			parts := strings.Split(token, "=")
	//			token = parts[0]
	//			if 2 == len(parts) {
	//				flag.Value = parts[1]
	//			}
	//		}
	//
	//		//cli.Push()
	//
	//		fmt.Println(token)
	//	}
	//	return cli, nil
}

/*
Usage is a func
*/
type Usage func()

/*
Value is an interface
*/
type Value interface {
	String() string
	Set(string) error
}

type flags []Flag

func (f flags) Cmd() string {
	return f[1].String()
}
func (f flags) Path() string {
	return f[0].String()
}
func (f flags) Usage() {
	f.Usage()
}
func (f flags) Workdir() string {
	dir, err := os.Getwd()
	if nil != err {
		log.Fatal(err)
	}
	return dir
}
