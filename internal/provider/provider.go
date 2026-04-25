package provider

// Provider abstracts reading and writing a single macOS setting.
// Each instance is bound to one setting at construction time.
type Provider interface {
	Read() (string, bool, error)
	Write(value string) error
	Delete() error
}
