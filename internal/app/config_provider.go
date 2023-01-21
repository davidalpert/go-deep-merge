package app

type ConfigProvider interface {
	GetValue(key string) (string, error)
}
