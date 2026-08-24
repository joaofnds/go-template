package marshal

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

type Codec interface {
	Marshaler
	Unmarshaler
}
