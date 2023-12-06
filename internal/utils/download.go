package utils

import (
	"io"
	"net/http"
)

// Download lets user download the file
func Download(w http.ResponseWriter, file io.Reader) error {
	// Set HTTP headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="export_products.csv"`)

	// copy data from pipe reader to HTTP response writer
	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	return nil
}
