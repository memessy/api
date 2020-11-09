package memessy

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
)

type Recognizer struct {
	FileField string
	Url       string
}

type memessyResponse struct {
	Text string `json:"text"`
}

func (r *Recognizer) Recognize(data []byte) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(r.FileField, "file")
	if err != nil {
		log.Error().Err(err).Msg("caught error while creating form")
		return "", err
	}
	_, err = part.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("caught error while writing to form")
		return "", err
	}
	writer.Close()
	req, err := http.NewRequest("POST", r.Url, body)
	if err != nil {
		log.Error().Err(err).Msg("caught error while creating new request")
		return "", err
	}
	req.Header.Set("content-type", writer.FormDataContentType())
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("caught error while sending request")
		return "", err
	}
	responseSchema := memessyResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseSchema)
	if err != nil {
		log.Error().Err(err).Msg("caught error while decoding response body")
		return "", err
	}
	return responseSchema.Text, nil
}
