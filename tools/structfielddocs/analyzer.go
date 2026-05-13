// Package structfielddocs предоставляет статический анализатор, требующий
// doc-комментарий у каждого экспортируемого поля экспортируемой структуры.
//
// Что считается doc-комментарием:
//   - комментарий строкой выше поля (Field.Doc).
//
// Inline-комментарий справа от поля (Field.Comment) НЕ считается
// doc-комментарием — стиль документации единый, всегда строкой выше.
//
// Что игнорируется:
//   - сгенерированные файлы (// Code generated ... DO NOT EDIT.) — их форму
//     диктует кодген, требовать на них doc-комментарии бессмысленно;
//   - неэкспортируемые поля;
//   - поля неэкспортируемых структур;
//   - встроенные (embedded) поля — у них нет собственного имени;
//   - анонимные структуры (literal-выражения) — у них нет имени типа.
package structfielddocs

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer — точка подключения к go/analysis.
var Analyzer = &analysis.Analyzer{
	Name: "structfielddocs",
	Doc:  "checks that exported fields of exported structs have a doc comment",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if !ts.Name.IsExported() {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}
			checkStruct(pass, ts.Name.Name, st)
			return true
		})
	}
	//nolint:nilnil // go/analysis API: Run возвращает (Facts, error); анализатор не публикует Facts, ошибок не имеет.
	return nil, nil
}

func checkStruct(pass *analysis.Pass, typeName string, st *ast.StructType) {
	if st.Fields == nil {
		return
	}
	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		if fieldHasComment(field) {
			continue
		}
		for _, name := range field.Names {
			if !name.IsExported() {
				continue
			}
			pass.Reportf(
				name.Pos(),
				"exported field %s of exported struct %s must have a doc comment",
				name.Name,
				typeName,
			)
		}
	}
}

func fieldHasComment(field *ast.Field) bool {
	return field.Doc != nil && len(field.Doc.List) > 0
}
