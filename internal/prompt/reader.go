package prompt

//go:generate mockgen -source=reader.go -destination=../../mock/prompt/reader_mock.go -package=mock_prompt
type Reader interface {
	ReadLine(prompt string) (string, error)
	ReadCommand(prompt string) error
}
