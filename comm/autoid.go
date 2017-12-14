package comm

var globalID int

func AutoIncrementID() int {
	globalID++
	return globalID
}
