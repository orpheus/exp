package interfaces

type Logger interface {
	Log(message string) error
}
