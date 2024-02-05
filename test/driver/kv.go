package driver

import (
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type KVDriver struct {
	url string
}

func NewKVDriver(baseURL string) *KVDriver {
	return &KVDriver{baseURL + "/kv"}
}

func (d *KVDriver) GetReq(key string) (*http.Response, error) {
	return req.Get(d.url+"/"+key, nil)
}

func (d *KVDriver) Get(key string) (string, error) {
	b, err := makeRequest(
		http.StatusOK,
		func() (*http.Response, error) { return d.GetReq(key) },
	)

	return string(b), err
}

func (d *KVDriver) MustGet(key string) string {
	return matchers.Must2(d.Get(key))
}

func (d *KVDriver) SetReq(key, value string) (*http.Response, error) {
	return req.Post(d.url+"/"+key+"/"+value, nil, nil)
}

func (d *KVDriver) Set(key, value string) error {
	return makeJSONRequest(params{
		status: http.StatusCreated,
		req:    func() (*http.Response, error) { return d.SetReq(key, value) },
	})
}

func (d *KVDriver) MustSet(key, value string) *http.Response {
	return matchers.Must2(d.SetReq(key, value))
}

func (d *KVDriver) DelReq(key string) (*http.Response, error) {
	return req.Delete(d.url+"/"+key, nil)
}

func (d *KVDriver) Del(key string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req:    func() (*http.Response, error) { return d.DelReq(key) },
	})
}

func (d *KVDriver) MustDel(key string) {
	matchers.Must(d.Del(key))
}
