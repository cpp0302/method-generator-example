package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

//go:generate go run . "../pkg"

const outputFileName = "method_gen.go"

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("usage: %v [package path]\n", filepath.Base(os.Args[0]))
		return
	}

	packagePath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := exec(packagePath); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func exec(packagePath string) error {

	structNames, err := getStructNames(packagePath)
	if err != nil {
		return err
	}

	fmt.Println("target")
	for _, v := range structNames {
		fmt.Printf("  %+v\n", v)
	}

	packageName := filepath.Base(packagePath)

	templData := struct {
		PackageName string
		StructNames []string
	}{
		PackageName: packageName,
		StructNames: structNames,
	}

	var w bytes.Buffer
	tmpl := template.Must(template.New("mytemplate").Parse(templStr))
	if err := tmpl.Execute(&w, templData); err != nil {
		return err
	}

	bytes, err := format.Source(w.Bytes())
	if err != nil {
		return err
	}

	bytes, err = imports.Process(outputFileName, bytes, nil)
	if err != nil {
		return err
	}

	outputFilePath := filepath.Join(packagePath, outputFileName)
	if err := ioutil.WriteFile(outputFilePath, bytes, 0644); err != nil {
		return err
	}

	fmt.Printf("wrote: %v", outputFilePath)

	return nil
}

// getStructNames パッケージ配下のソースコードを構文解析し、構造体名の一覧を取得する
func getStructNames(packagePath string) ([]string, error) {

	cfg := &packages.Config{
		// 型情報を取得するモードを指定する。必要であれば他のモードも追加して情報を取得可能
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}

	// パッケージ情報をロード
	pPkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, err
	}

	structNames := []string{}

	// パッケージ配下に存在する各要素の名前を取得
	for _, pPkg := range pPkgs {
		pkg := pPkg.Types

		for _, name := range pkg.Scope().Names() {

			// typeで定義しているかどうかチェック
			obj, ok := pkg.Scope().Lookup(name).(*types.TypeName)
			if !ok {
				continue
			}

			// structかどうかチェック
			if _, ok := obj.Type().Underlying().(*types.Struct); !ok {
				continue
			}

			structNames = append(structNames, obj.Name())
		}
	}

	return structNames, nil
}

const templStr = `// Code generated by gen/genmethods.go; DO NOT EDIT.

package {{ .PackageName }}

{{ range $structName := .StructNames }}
// PrintType 型情報を標準出力する
func (s {{ $structName }}) PrintType() {
	fmt.Printf("%T\n", s)
}
{{ end }}
`