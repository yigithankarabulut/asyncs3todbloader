package customerror

var (
	ErrBucketNotFound    = New("bucket not found", true)
	ErrObjectNotFound    = New("object not found", true)
	ErrNilObjectFields   = New("object fields are nil", true)
	ErrChannelClosed     = New("channel closed", true)
	ErrGetObjectFailed   = New("get s3 object failed", true)
	ErrFileScanFailed    = New("file scan failed", true)
	ErrIDExists          = New("id already exists", false)
	ErrETagExists        = New("etag already exists", true)
	ErrCreateIndexFailed = New("failed to create index", true)
	ErrCreateObjectInfo  = New("failed to create object info", true)
	ErrCreateProduct     = New("failed to create product", true)
)

type CustomError interface {
	Wrap(err error) CustomError
	Unwrap() error
	AddData(any) CustomError
	DestroyData() CustomError
	Error() string
}

type Error struct {
	Err      error
	Message  string
	Data     any `json:"-"`
	Loggable bool
}

func (e *Error) Wrap(err error) CustomError {
	e.Err = err
	return e
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) AddData(d any) CustomError {
	e.Data = d
	return e
}

func (e *Error) DestroyData() CustomError {
	e.Data = nil
	return e
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ", " + e.Message
	}
	return e.Message
}

func New(message string, l bool) CustomError {
	return &Error{
		Message:  message,
		Loggable: l,
	}
}
