package usecases

type Logger interface {
	Log(message string) error
}
