package port

import "io"

// Standard status messages for all operations
var (
	ListSuccess   StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "list retrieved successfully", Success: true}
	FetchSuccess  StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "data retrieved successfully", Success: true}
	CreateSuccess StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 201, Message: "resource created successfully", Success: true}
	UpdateSuccess StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "resource updated successfully", Success: true}
	DeleteSuccess StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "resource deleted successfully", Success: true}
	CustomEnv     StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "This is environment specific", Success: true}
)

// OTP-related status constants
var (
	OTPSuccess     StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "OTP generated successfully", Success: true}
	OTPAuthSuccess StatusCodeAndMessage = StatusCodeAndMessage{StatusCode: 200, Message: "OTP authenticated successfully", Success: true}
)

// StatusCodeAndMessage is embedded in all response structs
// Provides consistent status code, success flag, and message
type StatusCodeAndMessage struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}

// Status returns HTTP status code (interface compliance)
func (s StatusCodeAndMessage) Status() int {
	return s.StatusCode
}

func (s StatusCodeAndMessage) ResponseType() string {
	return "standard"
}

func (s StatusCodeAndMessage) GetContentType() string {
	return "application/json"
}

func (s StatusCodeAndMessage) GetContentDisposition() string {
	return ""
}

func (s StatusCodeAndMessage) Object() []byte {
	return nil
}

// FileResponse for file downloads/uploads
type FileResponse struct {
	ContentDisposition string
	ContentType        string
	Data               []byte        // Memory-based payload
	Reader             io.ReadCloser // Optional streaming source
}

func (s FileResponse) GetContentType() string {
	return s.ContentType
}

func (s FileResponse) GetContentDisposition() string {
	return s.ContentDisposition
}

func (s FileResponse) ResponseType() string {
	return "file"
}

func (s FileResponse) Status() int {
	return 200
}

func (s FileResponse) Object() []byte {
	return s.Data
}

// Stream copies Reader to w if available; else writes Data
func (s FileResponse) Stream(w io.Writer) error {
	if s.Reader == nil {
		if len(s.Data) > 0 {
			_, err := w.Write(s.Data)
			return err
		}
		return nil
	}
	defer s.Reader.Close()
	_, err := io.Copy(w, s.Reader)
	return err
}

// MetaDataResponse provides pagination metadata
// Embed this in list response structs
type MetaDataResponse struct {
	Skip                 uint64 `json:"skip,default=0"`
	Limit                uint64 `json:"limit,default=10"`
	OrderBy              string `json:"order_by,omitempty"`
	SortType             string `json:"sort_type,omitempty"`
	TotalRecordsCount    int    `json:"total_records_count,omitempty"`
	ReturnedRecordsCount uint64 `json:"returned_records_count"`
}

// Helper function to create metadata response
func NewMetaDataResponse(skip, limit, total uint64) MetaDataResponse {
	return MetaDataResponse{
		Skip:                 skip,
		Limit:                limit,
		TotalRecordsCount:    int(total),
		ReturnedRecordsCount: limit,
	}
}
