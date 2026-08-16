package main

import (
	"fmt"
	"os"
	"path/filepath"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	dir, err := os.MkdirTemp("", "ppt-parse-csv-file")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "data.csv")
	content := "# comment\nname;value;unit\ncpu;10;%\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		panic(err)
	}

	rows := pptutil.ParseCSVFile(path, ";", 3, true)

	for _, row := range rows {
		fmt.Println(row)
	}
}
