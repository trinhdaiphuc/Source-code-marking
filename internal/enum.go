package internal

type OrderType int8

const (
	ASC  OrderType = 0
	DESC OrderType = 1
)

func (orderType OrderType) String() string {
	names := []string{
		"ASC",
		"DESC",
	}
	if orderType < ASC || orderType > DESC {
		return "Unknown"
	}
	return names[orderType]
}
