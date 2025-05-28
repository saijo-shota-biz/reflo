package prompt

type Reader interface {
	ReadLine(prompt string) (string, error)
}
