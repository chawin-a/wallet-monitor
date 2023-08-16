package utils

func Must[V any](t V, err error) V {
	if err != nil {
		panic(err)
	}
	return t
}

func MustOk[V any](t V, ok bool) V {
	if !ok {
		panic("should be ok")
	}
	return t
}
