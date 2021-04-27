package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var signatureRequired = make(map[string]struct{})

func requireSignature(path string) {
	signatureRequired[path] = struct{}{}
}

func generateSignature(secret, path string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(path))

	return hex.EncodeToString(h.Sum(nil))
}

func requiresSignature(path string) bool {
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
	}

	var body io.Reader
	switch method {
	case http.MethodGet:
		path += "?" + qS + "&signature=" + sig
	case http.MethodPost:
		body = strings.NewReader(qS + "&signature=" + sig)
	}

	fullUrl := baseApi + path

	log.Printf("Calling %s", fullUrl)

	r, err := http.NewRequest(method, fullUrl, body)
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

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var resErr model.BinanceError
		err = json.Unmarshal(resBody, &resErr)
		if err != nil {
			return nil, err
		}
		return nil, resErr
	}

	return resBody, nil
}
