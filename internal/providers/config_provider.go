package providers

type Interface interface {
	GetValue(key string) (string, error)
	GetValueTree(prefix string) (map[string]string, error)
	//SetValue(key string) (string, error)
}
