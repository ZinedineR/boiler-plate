package rand

import gonanoid "github.com/matoous/go-nanoid/v2"

func Generate(rawAlphabet string, size int, prefix ...string) (string, error) {
	r, e := gonanoid.Generate(rawAlphabet, size)
	if e != nil {
		return "", e
	}
	if len(prefix) > 0 {
		return prefix[0] + r, nil
	}
	return r, nil
}

func GenerateNumeric(size int, prefix ...string) (string, error) {
	r, e := gonanoid.Generate("123456789", size)
	if e != nil {
		return "", e
	}
	if len(prefix) > 0 {
		return prefix[0] + r, nil
	}
	return r, nil
}
