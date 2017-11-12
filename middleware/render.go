package middleware

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ayerdudu.com/utils"

	"github.com/labstack/echo"
)

//RenderOpt 自定义 Template 参数
type RenderOpt struct {
	Directory string
	Suffix    string
	DevMode   bool
}

type tmplPath struct {
	name   string
	path   string
	suffix string
}

//Template *
type Template struct {
	templates *template.Template
	tmplPaths map[string]tmplPath
	suffix    string
	directory string
	devmode   bool
}

//MwRender Echo 自定义 Render
func MwRender(opts ...RenderOpt) *Template {
	t := &Template{}

	t.tmplPaths = make(map[string]tmplPath)

	var opt RenderOpt

	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "./templates"
	}

	if len(opt.Suffix) == 0 {
		opt.Suffix = ".html"
	}

	t.directory = opt.Directory
	t.suffix = opt.Suffix
	t.devmode = opt.DevMode

	templatePathWalk := func(p string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		} else if f.IsDir() {
			return nil
		} else if (f.Mode() & os.ModeSymlink) > 0 {
			return nil
		}

		if f.Size() > 0 {
			if strings.HasSuffix(p, opt.Suffix) {
				t.tmplPaths[strings.TrimSuffix(p, opt.Suffix)] = tmplPath{path: p, name: f.Name(), suffix: opt.Suffix}
			}
		}
		return err
	}

	err := filepath.Walk(strings.TrimRight(opt.Directory, "/"), func(p string, f os.FileInfo, err error) error {
		return templatePathWalk(p, f, err)
	})

	for k, v := range t.tmplPaths {
		if strings.EqualFold(strings.ToLower(v.suffix), strings.ToLower(opt.Suffix)) {
			tk := k[len(opt.Directory)+1 : len(k)]
			htmlStr, _ := ioutil.ReadFile(v.path)
			htmlTxt := string(htmlStr)
			if len(htmlTxt) != 0 {
				var tpl *template.Template
				if t.templates == nil {
					t.templates = template.New(tk)
				}
				if tk == t.templates.Name() {
					tpl = t.templates
				} else {
					tpl = t.templates.New(tk)
				}
				tpl = tpl.Delims("{%", "%}")
				_, err = tpl.Parse(htmlTxt)
				t.templates = template.Must(t.templates, err)
			}
		}
	}
	return t
}

// Render renders a template document
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.devmode {
		t, err := template.ParseFiles(t.directory + utils.PathSeparator + name + t.suffix)
		if nil != err {
			log.Panicf("template parse glob error: %v", err)
			return err
		}
		t = t.Delims("{%", "%}")
		return t.Execute(w, data)
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
