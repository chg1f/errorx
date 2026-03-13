package i18n

import (
	"encoding/json"
	"io/fs"
	"log/slog"

	"github.com/BurntSushi/toml"
	"github.com/chg1f/errorx/v2"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// bundle stores the shared go-i18n bundle used by the package-level loader.
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

// localize resolves a message by message key using the configured go-i18n bundle.
func localize(message string, attrs []slog.Attr, langs ...language.Tag) (string, bool) {
	unique := make(map[string]struct{}, len(langs))
	tags := make([]string, 0, len(langs))
	for i := range langs {
		s := langs[i].String()
		if _, ok := unique[s]; ok {
			continue
		}
		unique[s] = struct{}{}
		tags = append(tags, s)
	}
	localizer := goi18n.NewLocalizer(bundle, tags...)
	text, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    message,
		TemplateData: attrsToMap(attrs),
	})
	if err != nil {
		return "", false
	}
	return text, true
}

// attrsToMap converts structured slog attributes into go-i18n template data.
func attrsToMap(attrs []slog.Attr) map[string]any {
	data := make(map[string]any, len(attrs))
	for _, attr := range attrs {
		attr.Value = attr.Value.Resolve()
		if attr.Key == "" {
			continue
		}
		data[attr.Key] = valueToAny(attr.Value)
	}
	return data
}

// valueToAny normalizes slog values into plain Go values for templates.
func valueToAny(value slog.Value) any {
	switch value.Kind() {
	case slog.KindBool:
		return value.Bool()
	case slog.KindDuration:
		return value.Duration()
	case slog.KindFloat64:
		return value.Float64()
	case slog.KindInt64:
		return value.Int64()
	case slog.KindString:
		return value.String()
	case slog.KindTime:
		return value.Time()
	case slog.KindUint64:
		return value.Uint64()
	case slog.KindGroup:
		group := value.Group()
		return attrsToMap(group)
	case slog.KindLogValuer:
		return valueToAny(value.Resolve())
	default:
		return value.Any()
	}
}

func init() {
	errorx.Localize = localize
}
