package main

import (
	"github.com/alecthomas/kong"
	"github.com/nullism/gotoweb/builder"
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
	Build Build `cmd:"build" description:"Build the project"`
	New   New   `cmd:"new" description:"Create a new project"`
}

type Page struct {
	Title string
	Body  string
}

func main() {
	ctx := kong.Parse(&CLI)
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
		bldr.Build()

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
		ctx.PrintUsage(true)
	}

}
