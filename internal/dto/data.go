package dto

const INCORRECT_DATA_TYPE = "Incorrect data type "

type Data struct {
	Value any
}

func ParceData[T any](data Data) T {
	switch v := data.Value.(type) {
	case T:
		return v
	}

	panic(INCORRECT_DATA_TYPE)
}
