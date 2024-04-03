package constant

var (
	ErrBucketNotFound  = "bucket not found"
	ErrObjectNotFound  = "object not found"
	ErrNilObjectFields = "object fields are nil"
	ErrChannelClosed   = "channel closed"
	ErrGetObjectFailed = "get s3 object failed"
	ErrFileScanFailed  = "file scan failed"
)

var (
	ErrUnknown           = "unknown error"
	ErrIDExists          = "id already exists"
	ErrETagExists        = "etag already exists"
	ErrCreateIndexFailed = "failed to create index"
	ErrCreateObjectInfo  = "failed to create object file. already exists"
	ErrCreateProduct     = "failed to create product"
)
