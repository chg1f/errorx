package i18n

import (
	"encoding/json"
	"io/fs"

	"github.com/BurntSushi/toml"
	"github.com/chg1f/errorx/v2"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var bundle *goi18n.Bundle

func init() {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)
}

// LoadFiles walks the file system and loads all files into the default bundle.
func LoadFiles(fsys fs.FS) error {
	if fsys == nil {
		return fs.ErrInvalid
	}
	return fs.WalkDir(fsys, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		bs, readErr := fs.ReadFile(fsys, name)
		if readErr != nil {
			return readErr
		}
		_, parseErr := bundle.ParseMessageFileBytes(bs, name)
		return parseErr
	})
}

// localize resolves a message by locale key using the configured go-i18n bundle.
func localize(lang language.Tag, locale string, values map[string]any) (string, bool) {
	localizer := goi18n.NewLocalizer(bundle, lang.String())
	message, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    locale,
		TemplateData: values,
	})
	if err != nil {
		return "", false
	}
	return message, true
}

func init() {
	errorx.Localize = localize
}
