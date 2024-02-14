package request

const (
	KeyRequestID Key = "request_id"
	KeyTraceID   Key = "trace_id"
	KeySpanID    Key = "span_id"
	Bearer       Key = "Bearer"
	Token        Key = "auth_token"
)

type Key string

func (k Key) String() string {
	return string(k)
}
