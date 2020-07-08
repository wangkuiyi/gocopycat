package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	src := flag.String("from", "", "Go source directory to be copied")
	pkg := flag.String("pkg", "", "Only list declarations in the package")
	dst := flag.String("to", "", "Where to write output files")
	flag.Parse()

	if *dst != "" {
		os.MkdirAll(*dst, 0755)
	}

	if e := copyDir(*src, *pkg, *dst); e != nil {
		log.Fatal(e)
	}
}

// copyDir parse Go source files in directory src and generate files with the
// same names in directory dst.
//
// A directory might contain more than one packages, for example,
// https://golang.org/pkg/go/ast contains ast and ast_test.  Without the
// parameter pkg, copyDir copycat all files and all packages in src; otherwise,
// it copycats only files implement the package pkg.
func copyDir(src, pkg, dst string) error {
	log.Printf("Copycat %s ...", src)
	fset := token.NewFileSet() // positions are relative to fset
	pkgs, e := parser.ParseDir(fset, src, nil, 0)
	if e != nil {
		return fmt.Errorf("Failed to parse directory %s: %v", src, e)
	}

	for pn, p := range pkgs {
		if pn == pkg || pkg == "" {
			pn, e := fullPkgName(src)
			if e != nil {
				return e
			}
			if e := copyPackage(pn, p, dst); e != nil {
				return e
			}
		}
	}
	return nil
}

func getGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func fullPkgName(dir string) (string, error) {
	gosrc, e := filepath.Abs(path.Join(getGoPath(), "src/"))
	if e != nil {
		return "", e
	}

	dir, e = filepath.Abs(dir)
	if e != nil {
		return "", e
	}

	if !strings.HasPrefix(dir, gosrc) {
		return "", fmt.Errorf("We relies on GOPATH %s to derive the full package name of %s; however, GOPATH is not a prefix of the source dir. Please confirm if you have Go source tree in $GOPAHT/src", gosrc, dir)
	}
	return strings.TrimPrefix(strings.TrimPrefix(dir, gosrc), "/"), nil
}

func shortPkgName(full string) string {
	return path.Base(full)
}

func copyPackage(pn string, pkg *ast.Package, dst string) error {
	for fn, file := range pkg.Files {
		e := copyFile(pn, path.Join(dst, path.Base(fn)), file)
		if e != nil {
			return e
		}
	}
	return nil
}

func copyFile(pn, fn string, file *ast.File) error {
	o, e := os.Create(fn)
	if e != nil {
		return fmt.Errorf("Cannot create output file %s: %v", fn, e)
	}
	defer o.Close()

	fmt.Fprintf(o, "package %s\n", pn)

	for _, d := range file.Decls {
		var e error
		switch v := d.(type) {
		case *ast.GenDecl:
			e = copyType(pn, v, o)
		case *ast.FuncDecl:
			e = copyFunc(pn, v, o)
		}

		if e != nil {
			return e
		}
	}
	return nil
}

func copyType(pn string, decl *ast.GenDecl, o io.Writer) error {
	// It is a type declaration, other than import, const, or variable.
	if decl.Tok == token.TYPE {
		if e := copyTypeSpecs(pn, decl.Specs, o); e != nil {
			return e
		}
	}
	return nil
}

func copyTypeSpecs(name string, specs []ast.Spec, o io.Writer) error {
	for _, s := range specs {
		v := s.(*ast.TypeSpec)
		if token.IsExported(v.Name.Name) {
			fmt.Fprintf(o, "type %s=%s.%s\n", v.Name, name, v.Name)
		}
	}
	return nil
}

// copyFunc replaces the body of the function.  For example, suppose that in
// package yi, there is a function
//
// func Foo(a int) error {
//   the body
// }
//
// copyFunc replaces the body but keeps the signature.
//
// func Foo(a int) error {
//    yi.Foo(a)
// }
//
func copyFunc(name string, decl *ast.FuncDecl, o io.Writer) error {
	// Only prints exported function, not methods, because methods have been
	// copied by listTypeDecl using the `type=` syntax.
	if token.IsExported(decl.Name.Name) && decl.Recv == nil {
		// Remove body and print signature.
		decl.Body = rewriteBody(name, decl)
		fset := token.NewFileSet()
		format.Node(o, fset, decl)
		fmt.Fprintln(o)
	}
	return nil
}

func rewriteBody(name string, decl *ast.FuncDecl) *ast.BlockStmt {
	return &ast.BlockStmt{
		Lbrace: token.NoPos,
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: name},
						Sel: decl.Name,
					},
					Lparen:   token.NoPos,
					Args:     args(decl),
					Ellipsis: token.NoPos,
					Rparen:   token.NoPos,
				},
			},
		},
		Rbrace: token.NoPos,
	}
}

// args returns a []ast.Expr where each element is a *ast.Ident naming a
// parameter of the function declaration decl.
func args(decl *ast.FuncDecl) []ast.Expr {
	r := make([]ast.Expr, 0)
	for _, l := range decl.Type.Params.List {
		for _, n := range l.Names {
			r = append(r, n)
		}
	}
	return r
}
