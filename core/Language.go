package core

import "fmt"

type Language struct{
	name string
	version string
}

func NewLanguage(name string, version string) Language {
	return Language{
		name: name,
		version: version,
	}
}

func (l Language) GetName() string {
	return l.name
}

func (l Language) GetVersion() string {
	return l.version
}

func (l Language) String() string {
	return fmt.Sprintf("{name: %s, version: %s}", l.name, l.version)
}
