package req

type Headers map[string]string

func (headers Headers) Set(key, value string) {
	headers[key] = value
}

func MergeHeaders(headers ...Headers) Headers {
	merged := Headers{}
	for _, h := range headers {
		for k, v := range h {
			merged.Set(k, v)
		}
	}
	return merged
}
