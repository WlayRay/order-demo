package convertor

import "sync"

var (
	itemConvertor     *ItemConvertor
	itemConvertorOnce sync.Once

	orderConvertor     *OrderConvertor
	orderConvertorOnce sync.Once
)

func GetItemConvertor() *ItemConvertor {
	itemConvertorOnce.Do(func() {
		itemConvertor = &ItemConvertor{}
	})
	return itemConvertor
}

func GetOrderConvertor() *OrderConvertor {
	orderConvertorOnce.Do(func() {
		orderConvertor = &OrderConvertor{}
	})
	return orderConvertor
}
