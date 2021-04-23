package binance

import (
	"encoding/json"
	"gitlab.jaztec.info/checkers/checkers/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (a *api) client() *http.Client {
	c := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	return c
}

func (a *api) request(method string, path string, query url.Values, body io.Reader) (*http.Request, error) {
	if query != nil {
		qS := query.Encode()
		path += "?" + qS
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

func (a *api) doRequest(method, path string, q url.Values, body io.Reader) ([]byte, error) {
	if !a.checker.allowed {
		return nil, BinanceAtTimeout
	}

	req, err := a.request(method, path, q, body)
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
