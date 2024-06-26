package main

import (
	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/nullism/gotoweb/builder"
	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/models"
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
	switch ctx.Command() {
	case "new <name>":
		ns, err := newsite.New(CLI.New.Name)
		if err != nil {
			panic(err)
		}
		println("New command " + ns.Name)
	case "build":
		conf, err := models.SiteFromConfig()
		if err != nil {
			panic(err)
		}
		bldr, err := builder.New(conf)
		if err != nil {
			panic(err)
		}
		err = bldr.BuildAll()
		if err != nil {
			panic(err)
		}

		// tpl, err := template.ParseFiles("themes/default/index.html", "themes/default/templates/thumb.html")
		// if err != nil {
		// 	panic(err)
		// }
		// err = tpl.Execute(os.Stdout, &Page{Title: "Hello", Body: "World"})
		// if err != nil {
		// 	panic(err)
		// }
		// println("Building project")
	default:
		_ = ctx.PrintUsage(true)
	}

}
