package jack

var (
	DefaultLoader *Loader
)

func Load(name string) (*Client, error) {
	if DefaultLoader == nil {
		loader, err := NewLoader([]string{})
		if err != nil {
			panic(err)
		}

		DefaultLoader = loader
	}
	return DefaultLoader.Load(name)
}
