package ecnv

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func ResizeImage(w io.Writer, r io.Reader, i int, x, y uint, name, format string) error {
	if format == "" {
		format = filepath.Ext(name)
	}

	var img image.Image
	var err error
	switch strings.ToLower(format) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(r)
	case ".png":
		img, err = png.Decode(r)
	default:
		err = fmt.Errorf("unknown ext: %s", format)
	}
	if err != nil {
		return errors.Wrap(err, "failed to decode image from io.Reader")
	}

	img = resize.Resize(x, y, img, resize.Lanczos3)

	switch strings.ToLower(format) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(w, img, nil)
	case ".png":
		err = png.Encode(w, img)
	}
	return errors.Wrap(err, "failed to encode resized image")
}
