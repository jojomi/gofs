package gofs

func (x File) Renderer() (*Renderer, error) {
	return NewRendererFromFile(x)
}

func (x File) MustRenderer() *Renderer {
	renderer, err := x.Renderer()
	if err != nil {
		panic(err)
	}
	return renderer
}
