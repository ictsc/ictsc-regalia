package pg_test

import "iter"

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func collectErrIter[V any](iter iter.Seq2[V, error]) ([]V, error) {
	slice := make([]V, 0)
	for v, err := range iter {
		if err != nil {
			return nil, err
		}
		slice = append(slice, v)
	}
	return slice, nil
}
