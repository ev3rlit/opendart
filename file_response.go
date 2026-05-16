package opendart

// FileResponse contains bytes returned by OpenDART file APIs.
type FileResponse struct {
	ContentType string
	Body        []byte
}
