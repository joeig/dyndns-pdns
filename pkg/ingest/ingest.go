package ingest

// IngestModeType sets the IP address ingest mode
type IngestModeType string

// IngestMode defines the interface for IP address ingest mode processing
type IngestMode interface {
	Process() (*IPSet, error)
}
