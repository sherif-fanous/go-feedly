package mime

import (
	"bytes"
	"io"
	"mime/multipart"
)

func CreateMultipartMIMEAttachment(attachment io.Reader) (*bytes.Buffer, string, error) {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("cover", "coverImage")
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(part, attachment); err != nil {
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &body, writer.FormDataContentType(), nil
}
