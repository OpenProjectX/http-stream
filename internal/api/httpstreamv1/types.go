package httpstreamv1

type TransferRequest struct {
	Source   *HTTPRequest
	Target   *HTTPRequest
	Pipeline []*PipelineStage
}

type HTTPRequest struct {
	Method        string
	URL           string
	Headers       map[string]string
	ContentLength int64
	LocalPath     string
}

type PipelineStage struct {
	Name   string
	Config map[string]string
}

type TransferResponse struct {
	TransferID       string
	BytesTransferred int64
	SourceStatusCode int32
	TargetStatusCode int32
}
