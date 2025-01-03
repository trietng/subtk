package languages

const FallbackLanguage = "en"

type SupportedLanguages interface {}

func GetSupportedLanguages(table string) SupportedLanguages {
	return nil
}