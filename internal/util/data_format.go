package util

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
	"payment-gateway/internal/models"
	"strings"
)

const (
	contentTypeApplicationJson = "application/json"
	contentTypeTextXml         = "text/xml"
	contentTypeApplicationXml  = "application/xml"
)

// DecodeRequest decodes the incoming request based on content type
func DecodeRequest(r *http.Request, request *models.TransactionRequest) error {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case contentTypeApplicationJson:
		return json.NewDecoder(r.Body).Decode(request)
	case contentTypeTextXml:
		return xml.NewDecoder(r.Body).Decode(request)
	case contentTypeApplicationXml:
		return xml.NewDecoder(r.Body).Decode(request)
	default:
		return fmt.Errorf("unsupported content type")
	}
}

// EncodeResponse encode response
func EncodeResponse(w http.ResponseWriter, r *http.Request, response interface{}) error {
	acceptHeader := r.Header.Get("Accept")
	var contentType string

	// Split Accept header into parts and check each
	for _, part := range strings.Split(acceptHeader, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if part == "*/*" {
			contentType = contentTypeApplicationJson // Default if */* is explicitly listed
			break
		}
		mediaType, _, err := mime.ParseMediaType(part)
		if err != nil {
			continue // Skip invalid entries
		}
		switch mediaType {
		case contentTypeApplicationJson, contentTypeTextXml, contentTypeApplicationXml:
			contentType = mediaType
			break // Use the first supported type
		}
		if contentType != "" {
			break
		}
	}

	// Default to JSON if no valid Accept header or types
	if contentType == "" {
		contentType = contentTypeApplicationJson
	}

	w.Header().Set("Content-Type", contentType)

	switch contentType {
	case contentTypeApplicationJson:
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
	case contentTypeTextXml, contentTypeApplicationXml:
		if err := xml.NewEncoder(w).Encode(response); err != nil {
			return fmt.Errorf("error encoding XML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported Accept type: %s", contentType)
	}
	return nil
}
