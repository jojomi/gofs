package gofs

import (
	"io"
	"text/template"
)

type Renderer struct {
	templateData []byte
	data         any
	funcs        template.FuncMap
}

func NewRendererFromFile(templateFile File) (*Renderer, error) {
	templateData, err := templateFile.Content()
	if err != nil {
		return nil, err
	}
	return NewRenderer(templateData), nil
}

func NewRenderer(templateData []byte) *Renderer {
	return &Renderer{
		templateData: templateData,
		funcs:        template.FuncMap{},
		data:         map[string]any{},
	}
}

func (x *Renderer) AddFuncs(funcs template.FuncMap) *Renderer {
	for name, f := range funcs {
		x.funcs[name] = f
	}
	return x
}

func (x *Renderer) WithData(data any) *Renderer {
	x.data = data
	return x
}

func (x *Renderer) RenderToFile(f File) error {
	w, err := f.Writer()
	if err != nil {
		return err
	}
	defer w.Close()
	return x.RenderTo(w)
}

func (x *Renderer) RenderTo(w io.Writer) error {
	tmpl := template.New("gofs-template")
	tmpl.Funcs(x.funcs)

	tmpl, err := tmpl.Parse(string(x.templateData))
	if err != nil {
		return err
	}

	return tmpl.Execute(w, x.data)
}
