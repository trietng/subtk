package languages

type SupportedLanguages interface {
	GetSupportedLanguages() map[string]struct{}
}