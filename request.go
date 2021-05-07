package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/jaztec/go-binance/model"
)

var (
	signatureRequired = make(map[string]struct{})
	signatureMut      = sync.Mutex{}
)

func requireSignature(path string) {
	signatureMut.Lock()
	defer signatureMut.Unlock()
	signatureRequired[path] = struct{}{}
}

func generateSignature(secret, path string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(path))

	return hex.EncodeToString(h.Sum(nil))
}

func requiresSignature(path string) bool {
	signatureMut.Lock()
	defer signatureMut.Unlock()
	_, ok := signatureRequired[path]
	return ok
}

func (a *api) client() *http.Client {
	c := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	return c
}

func (a *api) request(method string, path string, query Parameters) (*http.Request, error) {
	var sig string
	var qS string
	if query != nil {
		qS = query.Encode()
	}
	if requiresSignature(path) {
		sig = generateSignature(a.cfg.Secret, qS)
		qS += "&signature=" + sig
	}

	var body io.Reader
	switch method {
	case http.MethodGet:
		if qS != "" {
			path += "?" + qS
		}
	case http.MethodPost:
		body = strings.NewReader(qS)
	}

	fullURL := a.cfg.BaseURI + path

	_ = a.logger.Log("calling", fmt.Sprintf("%s %s", method, fullURL), "params", qS)

	r, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set(APIKeyHeaderName, a.cfg.Key)

	return r, nil
}

func (a *api) doRequest(method, path string, q Parameters) ([]byte, error) {
	if !a.checker.allowed {
		return nil, AtTimeout
	}

	req, err := a.request(method, path, q)
	if err != nil {
		return nil, err
	}

	res, err := a.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := a.checker.checkResponse(res); err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var resErr model.Error
		err = json.Unmarshal(resBody, &resErr)
		if err != nil {
			return nil, err
		}
		return nil, resErr
	}

	return resBody, nil
}
