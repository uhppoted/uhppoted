package conf

type Unmarshaler interface {
	UnmarshalConf([]byte) (interface{}, error)
}

func Unmarshal(bytes []byte, m interface{}) error {
	return nil
}
