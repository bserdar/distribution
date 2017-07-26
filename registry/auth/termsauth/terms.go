package rhauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/docker/distribution/context"
)

type restService struct {
	url *url.URL
}

var httpHeaders = http.Header{
	http.CanonicalHeaderKey("content-type"): []string{"application/json"},
	http.CanonicalHeaderKey("accept"):       []string{"application/json"},
}

func (s restService) get(ctx context.Context, path string, params url.Values) (*http.Response, []byte, error) {
	uri := *s.url
	uri.RawQuery = params.Encode()
	uri.Path = uri.Path + path
	log := context.GetLoggerWithFields(ctx, map[interface{}]interface{}{"method": "GET",
		"uri":    uri.String(),
		"params": params})

	log.Info("restsvc: calling")

	req := &http.Request{Method: "GET", URL: &uri, Header: httpHeaders}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithField("url", uri.String()).Error("RestServiceError", err.Error())
		t.Err(err)
		return &http.Response{}, nil, restServiceError{err, uri.String()}
	}

	context.GetLoggerWithFields(ctx, map[interface{}]interface{}{"status": resp.Status,
		"code": resp.StatusCode}).Infof("restsvc: completed")

	var errCode error
	switch r.StatusCode / 100 {
	case 2:
		break
	case 3:
		break
	case 4, 5:
		errCode = errors.New(fmt.Sprintf("service error: %v", r.StatusCode))
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return resp, []byte{}, errCode
	}
	return resp, b, errCode
}

type termsAckData struct {
	AcknowledgedTermsId int64
	Type                string
}

func findCustomerTermsAcknowledgements(ctx context.Context, termsSvc restService, long webCustomerId) ([]termsAckData, error) {
	resp, body, err := termsSvc.Get(ctx, fmt.Sprintf("ack/customerid=%d", webCustomerId), nil)
	if err != nil {
		return nil, err
	}
	var ackData []termsAckData
	if err = json.Unmarshal(body, &ackData); err != nil {
		return nil, err
	}
	return ackData
}
