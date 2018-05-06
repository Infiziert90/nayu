package find_mime

import (
	"gopkg.in/h2non/filetype.v1"
	"gopkg.in/h2non/filetype.v1/types"
)

var AllowedMime = []string{
	"application/json",
	"application/zip",
	"audio/x-wav",
	"image/bmp",
	"image/jpeg",
	"image/png",
	"image/svg+xml",
	"image/tiff",
	"image/gif",
	"video/mpeg",
	"video/quicktime",
	"video/x-msvideo",
}

func Find(head []byte) (types.Type, error) {
	return filetype.Match(head)
}
