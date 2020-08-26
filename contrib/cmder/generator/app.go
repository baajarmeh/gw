package generator

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Application struct {
	cmdEngine *cobra.Command
	out       string
	pkg       string
	genFunc   bool
	genName   bool
}

func (app *Application) BindVars() {
	app.cmdEngine.PersistentFlags().StringVar(&app.pkg, "pkg", ".", "Specifies pkg name")
	app.cmdEngine.PersistentFlags().StringVar(&app.out, "out", ".", "Specifies output file")
	app.cmdEngine.PersistentFlags().BoolVar(&app.genFunc, "genFunc", true, "Specifies gen func")
	app.cmdEngine.PersistentFlags().BoolVar(&app.genName, "genName", true, "Specifies gen name")
}

func Run() {
	var app = &Application{}
	app.cmdEngine = &cobra.Command{
		Use:   "gwgen",
		Short: "Gw framework code generator",
		Long:  "Gw framework code generator",
		Run: func(cmd *cobra.Command, args []string) {
			//b,_ := ioutil.ReadFile("router.go")
			//fset := token.NewFileSet()
			//f, err := parser.ParseFile(fset, "", string(b), parser.ParseComments)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//
			//// Print the imports from the file's AST.
			//for _, s := range f.Comments {
			//	fmt.Println(s.List[0].Text)
			//}
		},
	}
	app.BindVars()
	err := app.cmdEngine.Execute()
	if err != nil {
		fmt.Printf("gwgen fail, err: %v", err)
	}
}
