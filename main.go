package main

import (
	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/nullism/gotoweb/builder"
	"github.com/nullism/gotoweb/config"
	"github.com/nullism/gotoweb/fsys"
	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/newsite"
)

type Build struct {
	Foo string
}

type New struct {
	Name string `arg:""`
}

var CLI struct {
	Verbose bool  `short:"v" help:"Enable verbose output."`
	Build   Build `cmd:"build" help:"Build the project"`
	New     New   `cmd:"new" help:"Create a new project"`
}

type Page struct {
	Title string
	Body  string
}

func main() {
	ctx := kong.Parse(&CLI)
	if CLI.Verbose {
		logging.Configure(slog.LevelDebug)
	}
	files := &fsys.OsFileSystem{}
	switch ctx.Command() {
	case "new <name>":
		ns, err := newsite.New(CLI.New.Name, files)
		if err != nil {
			panic(err)
		}
		println("New command " + ns.Name)
	case "build":

		conf, err := config.SiteFromConfig(files)
		if err != nil {
			panic(err)
		}
		bldr, err := builder.New(conf, files)
		if err != nil {
			panic(err)
		}
		err = bldr.BuildAll()
		if err != nil {
			panic(err)
		}
	default:
		_ = ctx.PrintUsage(true)
	}

}
