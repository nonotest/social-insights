package export

import (
	"github.com/athletifit/social-network-insights/models"
)

// Exporter is our main interface.
// Takes a document and is able to export it.
type Exporter interface {
	Export(document Document)
}

// Document represent the physical document we export.
type Document struct {
	name     string
	dataSets []models.DataSet
}

// NewDocument creates a new document.
func NewDocument(name string, dataSets []models.DataSet) Document {
	return Document{
		name:     name,
		dataSets: dataSets,
	}
}
