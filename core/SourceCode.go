package core

import "fmt"

type SourceCode struct {
	language Language
	code     string
}

func NewSourceCode(language Language, code string) SourceCode {
	return SourceCode{
		language: language,
		code:     code,
	}
}

func (sc SourceCode) GetLanguage() Language {
	return sc.language
}

func (sc SourceCode) GetCode() string {
	return sc.code
}

func (sc *SourceCode) SetLanguage(language Language) {
	sc.language = language
}

func (sc *SourceCode) SetCode(code string) {
	sc.code = code
}

func (sc SourceCode) String() string {
	return fmt.Sprintf("{language: %s, code: %s}", sc.language, sc.code)
}