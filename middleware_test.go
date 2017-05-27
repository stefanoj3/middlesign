package middlesign_test

import (
	"testing"
	"net/http"
	"github.com/stefanoj3/middlesign"
	"net/http/httptest"
	"time"
)

var middlesignConfig *middlesign.MiddleSignConfig

func init() {
	middlesignConfig = middlesign.DefaultConfig("my_awesome_secret")
}

func TestMiddlewareShouldAllowOnlyValidRequests(t *testing.T) {
	testCases := map[string]struct {
		request *http.Request
		expectedHttpCode int
	}{
		"unsigned request":{
			request: getRequestWithoutSignature(t),
			expectedHttpCode: http.StatusUnauthorized,
		},
		"wrong signature":{
			request: getRequestWithWrongSignature(t),
			expectedHttpCode: http.StatusUnauthorized,
		},
		"too old timestamp":{
			request: getRequestWithTooOldTimestamp(t),
			expectedHttpCode: http.StatusUnauthorized,
		},
		"timestamp in the future":{
			request: getRequestWithTimestampInTheFuture(t),
			expectedHttpCode: http.StatusUnauthorized,
		},
		"valid request":{
			request: getValidRequest(t),
			expectedHttpCode: http.StatusAccepted,
		},
		"valid request with multiple parameters":{
			request: getValidRequestWithMultipleParameters(t),
			expectedHttpCode: http.StatusAccepted,
		},
	}

	requestIsValidHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}

	reuqestIsInvalidHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}

	middleware := middlesign.NewSignedRequestMiddleware(
		http.HandlerFunc(requestIsValidHandler),
		http.HandlerFunc(reuqestIsInvalidHandler),
		middlesignConfig,
	)



	for scenario, test := range testCases {
		resp := httptest.NewRecorder()
		middleware.ServeHTTP(resp, test.request)

		if test.expectedHttpCode != resp.Code {
			t.Errorf(
				"Expected %d http code, got %d [scenario: %s]",
				test.expectedHttpCode,
				resp,
				scenario,
			)
		}
	}
}

func getRequestWithoutSignature(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:176547/", nil)
	if err != nil {
		t.Fatal("Unable to continue with tests, cannot getRequestWithoutSignature")
	}

	q := req.URL.Query()
	q.Add(middlesignConfig.TimestampKey, time.Now().Format(middlesignConfig.TimestampFormat))

	req.URL.RawQuery = q.Encode()

	return req
}

func getRequestWithWrongSignature(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:176547/", nil)
	if err != nil {
		t.Fatal("Unable to continue with tests, cannot getRequestWithoutSignature")
	}

	q := req.URL.Query()
	q.Add(middlesignConfig.TimestampKey, time.Now().Format(middlesignConfig.TimestampFormat))
	q.Add(middlesignConfig.SignatureKey, "some wrong signature")

	req.URL.RawQuery = q.Encode()

	return req
}

func getRequestWithTooOldTimestamp(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:176547/", nil)
	if err != nil {
		t.Fatal("Unable to continue with tests, cannot getRequestWithoutSignature")
	}

	q := req.URL.Query()
	q.Add(middlesignConfig.TimestampKey, time.Now().Add(time.Hour * -10).Format(middlesignConfig.TimestampFormat))

	signature := middlesign.SignString(q.Encode(), middlesignConfig.Secret)
	q.Add(middlesignConfig.SignatureKey, signature)

	req.URL.RawQuery = q.Encode()

	return req
}

func getRequestWithTimestampInTheFuture(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:176547/", nil)
	if err != nil {
		t.Fatal("Unable to continue with tests, cannot getRequestWithoutSignature")
	}

	q := req.URL.Query()
	q.Add(middlesignConfig.TimestampKey, time.Now().Add(time.Hour * 10).Format(middlesignConfig.TimestampFormat))

	signature := middlesign.SignString(q.Encode(), middlesignConfig.Secret)
	q.Add(middlesignConfig.SignatureKey, signature)

	req.URL.RawQuery = q.Encode()

	return req
}

func getValidRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:176547/", nil)
	if err != nil {
		t.Fatal("Unable to continue with tests, cannot getRequestWithoutSignature")
	}

	q := req.URL.Query()
	q.Add(middlesignConfig.TimestampKey, time.Now().Format(middlesignConfig.TimestampFormat))

	signature := middlesign.SignString(q.Encode(), middlesignConfig.Secret)
	q.Add(middlesignConfig.SignatureKey, signature)

	req.URL.RawQuery = q.Encode()

	return req
}

func getValidRequestWithMultipleParameters(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:176547/", nil)
	if err != nil {
		t.Fatal("Unable to continue with tests, cannot getRequestWithoutSignature")
	}

	q := req.URL.Query()
	q.Add(middlesignConfig.TimestampKey, time.Now().Format(middlesignConfig.TimestampFormat))

	q.Add("param_name1", "some_value")
	q.Add("param_name2", "2f4r2f32r23r3")

	signature := middlesign.SignString(q.Encode(), middlesignConfig.Secret)
	q.Add(middlesignConfig.SignatureKey, signature)


	req.URL.RawQuery = q.Encode()

	return req
}
