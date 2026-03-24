package configs

import _ "embed"

// LanguageExtensionsJSON contains the language->extensions mapping used by --languages.
//
//go:embed language_extensions.json
var LanguageExtensionsJSON []byte
