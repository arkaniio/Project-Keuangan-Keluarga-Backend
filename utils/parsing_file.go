package utils

import "net/http"

func DetectContentType(buff []byte) string {

	content_type := http.DetectContentType(buff)
	if content_type != "jpg" && content_type != "jpeg" {
		return "Failed to detect the content type of the file!"
	}

	return ""

}
