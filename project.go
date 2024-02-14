package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"
)

const ext = ".template"

//go:embed templates/licenses templates/files/*
var fSys embed.FS

type Project struct {
	AppName   string
	Copyright string
	PkgName   string
	Version   string
	EnvPrefix string
	Folder    string
	License   string
	Docker    bool
	Sentry    bool
	Header    bool
	Database  Database
	ORM       ORM
	Router    Router

	Templates       map[string]string
	globalTemplates map[string]string

	skipFiles      map[string]struct{}
	packages       []string
	absolutePath   string
	executablePath string
}

func NewProject() Project {
	return Project{
		packages: []string{
			"get", // This allows easily getting the packages via exec
			"github.com/stretchr/testify",
			"github.com/spf13/cobra",
			"github.com/spf13/pflag",
			"github.com/spf13/viper",
		},
		Version:         getGoVersion(),
		globalTemplates: map[string]string{},
		executablePath:  getExecutableDirectory(),
		absolutePath:    getWorkingDirectory(),
	}
}

func (p *Project) Generate() error {
	if err := p.setup(); err != nil {
		return err
	}

	initApp, err := p.parseGoMod()
	if err != nil {
		return err
	}

	if err := makeFolder(p.Folder); err != nil {
		return err
	}

	if err := changeFolder(p.Folder); err != nil {
		return err
	}

	if initApp {
		if p.PkgName == "" {
			return fmt.Errorf("package name cannot be empty if go.mod does not exist yet")
		}

		log.Println("initializing go module")
		if err := exec.Command("go", "mod", "init", p.PkgName).Run(); err != nil {
			return err
		}
	}

	log.Println("getting packages")
	if err := exec.Command("go", p.packages...).Run(); err != nil {
		log.Println("WARN: unable to get packages")
		log.Printf("go %s\n", strings.Join(p.packages, " "))
	}

	if err := changeFolder(p.absolutePath); err != nil {
		return err
	}

	if err := p.makeFiles(); err != nil {
		return err
	}

	if err := changeFolder(p.Folder); err != nil {
		return err
	}

	log.Println("clean up")
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return err
	}
	log.Println("✓ tidy")

	if err := exec.Command("go", "fmt").Run(); err != nil {
		return err
	}
	log.Println("✓ format")

	return changeFolder(p.absolutePath)
}

func (p *Project) data() map[string]any {
	d := map[string]any{
		"AppName":      p.AppName,
		"Copyright":    p.Copyright,
		"Database":     p.Database,
		"Docker":       p.Docker,
		"Folder":       p.Folder,
		"Header":       p.Header,
		"License":      p.License,
		"ORM":          p.ORM,
		"PkgName":      p.PkgName,
		"Router":       p.Router,
		"Sentry":       p.Sentry,
		"Version":      p.Version,
		"EnvPrefix":    p.EnvPrefix + "_",
		"EnvPrefixVar": p.EnvPrefix,
	}

	return d
}

func (p *Project) makeFiles() error {
	log.Println("generating files from templates...")

	// Load templates
	templates, err := template.ParseFS(
		fSys,
		"templates/files/*"+ext,
		"templates/files/*/*"+ext,
		"templates/files/*/*/*"+ext,
		"templates/licenses/"+p.License+"/*"+ext,
	)
	if err != nil {
		return err
	}

	// Load global templates
	for name, content := range p.globalTemplates {
		templates, err = templates.Parse(fmt.Sprintf(`{{ define "%s" }}%s{{ end }}`, name, content))
		if err != nil {
			return err
		}
	}

	// Make templates
	root := "templates/files"
	if err := fs.WalkDir(fSys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == root {
			return err
		}

		// remove the "files" portion and the slash
		path = filepath.Clean(path[len(root)+1:])

		// change __application__ folder to user specified application folder
		file := p.replaceAppFolder(path)

		// Make directories as needed
		if d.IsDir() {
			return makeFolder(file)
		}

		file = file[:len(file)-len(ext)]
		if _, ok := p.skipFiles[file]; ok {
			return nil
		}

		var b bytes.Buffer
		if err := templates.ExecuteTemplate(&b, filepath.Base(path), p.data()); err != nil {
			return err
		}

		log.Println("making file:", file)
		return os.WriteFile(file, b.Bytes(), 0o640)
	}); err != nil {
		return err
	}

	// Make provided templates
	if len(p.Templates) > 0 {
		log.Println("creating user templates...")
	}
	for k, contents := range p.Templates {
		if contents != "" {
			templates, err = templates.Parse(contents)
			if err != nil {
				return err
			}

			if err := makeFolder(filepath.Dir(k)); err != nil {
				return nil
			}

			var b bytes.Buffer
			if err := templates.Execute(&b, p.data()); err != nil {
				return err
			}

			log.Println("making file:", k)
			return os.WriteFile(k, b.Bytes(), 0o640)
		}
	}

	return nil
}

func (p *Project) parseGoMod() (bool, error) {
	log.Println("checking go.mod...")
	needInit := true
	search := map[string]*string{
		"module ": &p.PkgName,
		"go ":     &p.Version,
	}

	return needInit, filepath.WalkDir(p.absolutePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Base(path) != "go.mod" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			for k, v := range search {
				if strings.HasPrefix(line, k) {
					*v = line[len(k):]
					delete(search, k)
					break
				}
			}

			if len(search) == 0 {
				break
			}
		}

		if len(search) != 0 {
			return fmt.Errorf("unable to parse go.mod: %s", path)
		}

		needInit = false
		return nil
	})
}

func (p *Project) replaceAppFolder(s string) string {
	app := "__application__"
	if p.Folder == "" {
		app += "/"
	}

	return strings.ReplaceAll(s, app, p.Folder)
}

func (p *Project) setup() error {
	if p.Copyright != "" {
		p.Copyright = fmt.Sprintf("Copyright © %d %s", time.Now().Year(), p.Copyright)
	}

	if err := p.setLicense(); err != nil {
		return err
	}

	if !p.Docker {
		p.skipFiles[".dockerignore"] = struct{}{}
		p.skipFiles["docker-compose.yml"] = struct{}{}
		p.skipFiles["Dockerfile"] = struct{}{}
	}

	if err := p.setDatabase(); err != nil {
		return err
	}

	if err := p.setORM(); err != nil {
		return err
	}

	if err := p.setRouter(); err != nil {
		return err
	}

	for k, f := range p.Templates {
		k = filepath.Clean(k)
		path, err := filepath.Abs(k)
		if err != nil {
			return fmt.Errorf("unable to parse template path: %w", err)
		}
		if !strings.HasPrefix(path, p.absolutePath) {
			return fmt.Errorf("templates should remain in project folder: %s", k)
		}

		n := p.replaceAppFolder(k)
		delete(p.Templates, k)

		p.Templates[n] = f
	}

	return nil
}

func (p *Project) setDatabase() error {
	var err error
	p.Database, err = findDatabase(p.Database.Name)
	if err != nil {
		return err
	}

	p.globalTemplates["Docker DB Env"] = p.Database.DockerEnv

	return nil
}

func (p *Project) setLicense() error {
	var err error
	p.License, err = findLicense(p.License)
	if err != nil {
		return err
	}

	if p.License == "" {
		p.skipFiles["LICENSE"] = struct{}{}
		p.globalTemplates["header.template"] = ""
	}

	return nil
}

func (p *Project) setORM() error {
	var err error
	p.ORM, err = findORM(p.ORM.Name)
	if err != nil {
		return err
	}

	p.packages = append(p.packages, p.ORM.Package)
	if driver := p.ORM.DBDriver[p.Database.Name]; driver != "" {
		p.ORM.Driver += p.ORM.DBDriver[p.Database.Name]
		p.packages = append(p.packages, p.ORM.Driver)
		p.globalTemplates["ORM Init"] = p.ORM.Init
	}

	return nil
}

func (p *Project) setRouter() error {
	var err error
	p.Router, err = findRouter(p.Router.Name)
	if err != nil {
		return err
	}

	p.packages = append(p.packages, p.Router.Package)
	return nil
}

func changeFolder(name string) error {
	if name == "" || name == "." {
		return nil
	}

	cur := getWorkingDirectory()

	dir, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	if cur == dir {
		return nil
	}

	log.Println("cd to folder:", name)
	return os.Chdir(name)
}

func getExecutableDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.Dir(ex)
}

func getGoVersion() string {
	v := runtime.Version()
	if !strings.HasPrefix(v, "go") {
		log.Fatalf("unknown go version: %s\n", v)
	}

	i := len(v)
	if strings.Count(v, ".") > 1 {
		i = strings.LastIndex(v, ".")
	}
	return v[2:i]
}

func getWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	return wd
}

func makeFolder(name string) error {
	if name == "" || name == "." {
		return nil
	}

	err := os.MkdirAll(name, 0o0755)
	if os.IsExist(err) {
		return nil
	}

	log.Println("making folder:", name)
	return err
}
