package output

import (
	"fmt"

	"l22.io/viesquery/internal/vies"
)

// Formatter defines the interface for output formatting
type Formatter interface {
	Format(result *vies.CheckVatResult) (string, error)
	FormatError(err error) (string, error)
}

// Manager manages different output formatters
type Manager struct {
	formatters map[string]Formatter
}

// NewManager creates a new formatter manager with default formatters
func NewManager() *Manager {
	return &Manager{
		formatters: map[string]Formatter{
			"plain": NewPlainFormatter(),
			"json":  NewJSONFormatter(),
		},
	}
}

// GetFormatter returns a formatter by name
func (m *Manager) GetFormatter(format string) (Formatter, error) {
	formatter, exists := m.formatters[format]
	if !exists {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	return formatter, nil
}

// RegisterFormatter registers a new formatter
func (m *Manager) RegisterFormatter(name string, formatter Formatter) {
	m.formatters[name] = formatter
}

// GetSupportedFormats returns a list of supported format names
func (m *Manager) GetSupportedFormats() []string {
	formats := make([]string, 0, len(m.formatters))
	for name := range m.formatters {
		formats = append(formats, name)
	}
	return formats
}
