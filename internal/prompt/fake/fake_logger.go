package fake

type Reader struct {
}

func (r Reader) ReadLine(string) (string, error) {
	return "", nil
}
