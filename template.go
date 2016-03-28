package binhtml

import (
	"html/template"
	"path/filepath"
)

type AssetFunc func(string) ([]byte, error)
type AssetDirFunc func(string) ([]string, error)

type BinTemplate struct {
	Asset    AssetFunc
	AssetDir AssetDirFunc
	fm       map[string]interface{}
	tmpl     *template.Template
}

func New(a AssetFunc, b AssetDirFunc) *BinTemplate {
	return &BinTemplate{
		Asset:    a,
		AssetDir: b,
		fm:       make(map[string]interface{}),
		tmpl:     template.New(""),
	}
}

func (t *BinTemplate) Funcs(fm map[string]interface{}) *BinTemplate {
	t.fm = fm
	t.tmpl = t.tmpl.Funcs(t.fm)
	return t
}

func (t *BinTemplate) Template() *template.Template {
	if t.tmpl == nil {
		t.tmpl = &template.Template{}
	}
	return t.tmpl
}

func (t *BinTemplate) LoadDirectory(directory string) error {

	files, err := t.AssetDir(directory)
	if err != nil {
		return err
	}

	for _, filePath := range files {
		contents, err := t.Asset(directory + "/" + filePath)
		if err != nil {
			return err
		}

		name := filepath.Base(filePath)

		if name != t.tmpl.Name() {
			t.tmpl = t.tmpl.New(name)
		}

		if _, err = t.tmpl.Parse(string(contents)); err != nil {
			return err
		}
	}

	return nil
}

func (t *BinTemplate) MustLoadDirectory(directory string) {
	if err := t.LoadDirectory(directory); err != nil {
		panic(err)
	}
}
