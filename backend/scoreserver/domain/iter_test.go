package domain_test

import "iter"

// テスト用に便利なイテレータ操作を提供する
// どこかに切り出したほうがいいかもしれない

func asErrIter[V any](iter iter.Seq[V]) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		for v := range iter {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func singleErrIter[V any](v V, err error) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		if !yield(v, err) {
			return
		}
	}
}
