package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Handle a serverless request
func Handle(req []byte) string {
	host := os.Getenv("gateway_hostname")
	if host == "" {
		host = "gateway.openfaas"
	}

	buf := bytes.NewBuffer(req)

	r, err := http.Post("http://"+host+":8080/function/sentimentanalysis", "application/text", buf)
	if err != nil {
		return err.Error()
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return fmt.Sprintf("Expected status 200, got %d", r.StatusCode)
	}

	type sentimentResponse struct {
		Polarity      float64 `json:"polarity,omitempty"`
		SentenceCount int     `json:"sentence_count,omitempty"`
		Subjectivity  float64 `json:"subjectivity,omitempty"`
	}

	var s sentimentResponse
	json.NewDecoder(r.Body).Decode(&s)

	if s.Polarity > 0.45 {
		return "That was probably positive"
	}
	return "That was probably negative"
}
