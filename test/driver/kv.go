package driver

import (
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type KVDriver struct {
	url     string
	headers req.Headers
}

func NewKVDriver(url string, headers req.Headers) *KVDriver {
	return &KVDriver{url: url + "/kv", headers: headers}
}

func (driver *KVDriver) GetReq(key string) (*http.Response, error) {
	return req.Get(driver.url+"/"+key, driver.headers)
}

func (driver *KVDriver) MustGetReq(key string) *http.Response {
	return matchers.Must2(driver.GetReq(key))
}

func (driver *KVDriver) Get(key string) (string, error) {
	bytes, err := makeRequest(
		http.StatusOK,
		func() (*http.Response, error) { return driver.GetReq(key) },
	)

	return string(bytes), err
}

func (driver *KVDriver) MustGet(key string) string {
	return matchers.Must2(driver.Get(key))
}

func (driver *KVDriver) SetReq(key, value string) (*http.Response, error) {
	return req.Post(driver.url+"/"+key+"/"+value, driver.headers, nil)
}

func (driver *KVDriver) Set(key, value string) error {
	return makeJSONRequest(params{
		status: http.StatusCreated,
		req:    func() (*http.Response, error) { return driver.SetReq(key, value) },
	})
}

func (driver *KVDriver) MustSet(key, value string) *http.Response {
	return matchers.Must2(driver.SetReq(key, value))
}

func (driver *KVDriver) DelReq(key string) (*http.Response, error) {
	return req.Delete(driver.url+"/"+key, driver.headers)
}

func (driver *KVDriver) Del(key string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req:    func() (*http.Response, error) { return driver.DelReq(key) },
	})
}

func (driver *KVDriver) MustDel(key string) {
	matchers.Must(driver.Del(key))
}
