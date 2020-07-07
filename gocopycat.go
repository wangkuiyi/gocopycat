package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	dir := flag.String("dir", "", "Go source directory to be parsed")
	pkg := flag.String("pkg", "ast", "Only list declarations in the package")
	flag.Parse()

	if e := listPackages(*dir, *pkg); e != nil {
		log.Fatal(e)
	}
}

func listPackages(dir, pkg string) error {
	fset := token.NewFileSet() // positions are relative to fset
	pkgs, e := parser.ParseDir(fset, "/tmp/go/src/go/ast/", nil, 0)
	if e != nil {
		return fmt.Errorf("Failed to parse directory %s: %v", dir, e)
	}

	for name, p := range pkgs {
		if name == pkg || pkg == "" {
			if e := listFiles(name, p); e != nil {
				return e
			}
		}
	}
	return nil
}

func listFiles(name string, pkg *ast.Package) error {
	for _, file := range pkg.Files {
		if e := listDeclarations(name, file); e != nil {
			return e
		}
	}
	return nil
}

func listDeclarations(name string, file *ast.File) error {
	for _, d := range file.Decls {
		var e error
		switch v := d.(type) {
		case *ast.GenDecl:
			e = listTypeDecl(name, v)
		case *ast.FuncDecl:
			e = listFuncDecl(name, v)
		}

		if e != nil {
			return e
		}
	}
	return nil
}

func listTypeDecl(name string, decl *ast.GenDecl) error {
	// It is a type declaration, other than import, const, or variable.
	if decl.Tok == token.TYPE {
		printComment(decl.Doc)
		if e := listTypeSpecs(name, decl.Specs); e != nil {
			return e
		}
	}
	return nil
}

func listFuncDecl(name string, decl *ast.FuncDecl) error {
	if token.IsExported(decl.Name.Name) {
		// Remove body and print signature.
		decl.Body = nil
		var sig bytes.Buffer
		fset := token.NewFileSet()
		format.Node(&sig, fset, decl)

		// Print a new body that
		fmt.Printf("%s\n", sig.String())
	}
	return nil
}

func listTypeSpecs(name string, specs []ast.Spec) error {
	for _, s := range specs {
		v := s.(*ast.TypeSpec)
		if token.IsExported(v.Name.Name) {
			printComment(v.Comment)
			fmt.Printf("type %s = %s.%s\n", v.Name, name, v.Name)
		}
	}
	return nil
}

func printComment(cmt *ast.CommentGroup) {
	if cmt != nil {
		fset := token.NewFileSet()
		var buf bytes.Buffer
		format.Node(&buf, fset, cmt)
		fmt.Println(buf.String())
	}
}
