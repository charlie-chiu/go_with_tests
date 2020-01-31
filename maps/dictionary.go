package maps

type Dictionary map[string]string
type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

const (
	ErrNotFound   = DictionaryErr("could not find the word you were search for")
	ErrWordExists = DictionaryErr("cannot add word because it already exists")
)

func (d Dictionary) Search(word string) (value string, err error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}

	return definition, nil
}

func (d Dictionary) Add(word, definition string) error {
	_, exists := d[word]
	if exists {
		return ErrWordExists
	}
	d[word] = definition

	return nil
}
