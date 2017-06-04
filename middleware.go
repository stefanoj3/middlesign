package middlesign

import (
	"math"
	"net/http"
	"time"
)

// NewSignedRequestMiddleware creates the handler with the given configuration
func NewSignedRequestMiddleware(nextHandler http.Handler, errorHandler http.Handler, config MiddleSignConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isRequestValid(r, &config) {
			nextHandler.ServeHTTP(w, r)
		} else {
			errorHandler.ServeHTTP(w, r)
		}
	})
}

func isRequestValid(r *http.Request, config *MiddleSignConfig) bool {
	query := r.URL.Query()

	timestampAsString := query.Get(config.TimestampKey)

	if timestampAsString == "" { // if there is no timestamp we cannot validate the request
		return false
	}

	timestamp, err := time.Parse(config.TimestampFormat, timestampAsString)
	if err != nil { // if the timestamp is in wrong format we cannot validate the request
		return false
	}

	diff := math.Abs(float64(time.Now().Unix() - timestamp.Unix()))

	if diff > config.Threshold { // timestamp is too old, request is not valid
		return false
	}

	signature := query.Get(config.SignatureKey)
	if signature == "" { // if there is no signature we cannot validate the request
		return false
	}

	query.Del(config.SignatureKey)

	source := query.Encode()

	valid, _ := IsSignatureValid(source, config.Secret, signature)

	return valid
}
