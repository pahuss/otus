package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	app2 "github.com/pahuss/otus/app"
	"os"
)

func main() {
	arguments := os.Args[1:]
	command := arguments[0]

	commands := map[string]string{
		"import": "import",
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic("Load environment error")
	}

	if command, ok := commands[command]; !ok {
		panic(fmt.Sprintf("Unknown command %s", command))
	}

	app := app2.NewApp(context.Background())
	app.InitDb(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), "localhost:3306", "tcp")

	if command == "import" {
		runImport(app)
	}

}

func runImport(app *app2.App) {
	app2.CopyRemoteCsvFile(
		"https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.csv",
		"people.csv",
	)
	app2.ReadCsvPeopleToDatabase("people.csv", app)
}
