package binhtml

import (
	"html/template"
	"path/filepath"
)

type BinTemplate struct {
	Asset    func(string) ([]byte, error)
	AssetDir func(string) ([]string, error)
	fm       map[string]interface{}
	tmpl     *template.Template
}

func New(
	a func(string) ([]byte, error),
	b func(string) ([]string, error),
) *BinTemplate {
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

func (t *BinTemplate) Load(directory string) error {
	files, err := t.AssetDir(directory)
	if err != nil {
		return err
	}

	for _, filePath := range files {
		// load subfolders recursively
		if filepath.Ext(filePath) == "" {
			err := t.Load(filepath.Join(directory, filePath))
			if err != nil {
				return err
			}
			continue
		}

		contents, err := t.Asset(filepath.Join(directory, filePath))
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

func (t *BinTemplate) MustLoad(directory string) *BinTemplate {
	if err := t.Load(directory); err != nil {
		panic("bindata templates loading failed: " + err.Error())
	}
	return t
}
