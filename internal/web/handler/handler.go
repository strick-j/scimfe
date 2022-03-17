package handler

import (
	"io"

	"github.com/strick-j/scimfe/internal/model"
	"github.com/strick-j/scimfe/internal/web"
)

// UnmarshalAndValidate unmarshals request from JSON in HTTP request and runs validation.
func UnmarshalAndValidate(r io.ReadCloser, dst interface{}) error {
	if err := web.UnmarshalJSON(r, dst); err != nil {
		return err
	}

	return model.Validate(dst)
}
