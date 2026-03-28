package httpstreamv1

type TransferRequest struct {
	Source   *HTTPRequest     `json:"source"`
	Target   *HTTPRequest     `json:"target"`
	Pipeline []*PipelineStage `json:"pipeline,omitempty"`
}

type HTTPRequest struct {
	Method        string            `json:"method"`
	URL           string            `json:"url"`
	Headers       map[string]string `json:"headers,omitempty"`
	ContentLength int64             `json:"content_length,omitempty"`
}

type PipelineStage struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config,omitempty"`
}

type TransferResponse struct {
	TransferID       string `json:"transfer_id"`
	BytesTransferred int64  `json:"bytes_transferred"`
	SourceStatusCode int32  `json:"source_status_code"`
	TargetStatusCode int32  `json:"target_status_code"`
}
