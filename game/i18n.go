package game

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed locales/en.json
var enJSON []byte

//go:embed locales/it.json
var itJSON []byte

// Translator holds a flat map of translation keys to strings for one language.
type Translator struct {
	strings map[string]string
}

// newTranslator loads the translator for "en" or "it" (defaults to "en").
func newTranslator(lang string) (*Translator, error) {
	data := enJSON
	if lang == "it" {
		data = itJSON
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &Translator{strings: m}, nil
}

// T looks up key and formats it with optional args (like fmt.Sprintf).
// Returns "?key?" if the key is missing so missing translations are obvious.
func (t *Translator) T(key string, args ...any) string {
	s, ok := t.strings[key]
	if !ok {
		return "?" + key + "?"
	}
	if len(args) == 0 {
		return s
	}
	return fmt.Sprintf(s, args...)
}
