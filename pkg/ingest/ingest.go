package ingest

// ModeType sets the IP address ingest mode
type ModeType string

// Mode defines the interface for IP address ingest mode processing
type Mode interface {
	Process() (*IPSet, error)
}
