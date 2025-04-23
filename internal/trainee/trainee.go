package trainee

import "errors"

func Hello(name string, language string) (string, error) {
	if name == "" {
		name = "World"
	}

	prefix := ""

	switch language {
	case "english":
		prefix = "Hello"
	case "spanish":
		prefix = "Hola"
	case "german":
		prefix = "Hallo"
	default:
		return "", errors.New("need to provide a supported language")
	}

	return prefix + " " + name, nil
}

func FibRecursive(n uint) uint {
	if n <= 2 {
		return 1
	}
	return FibRecursive(n-1) + FibRecursive(n-2)
}

func FibIterative(position uint) uint {
	slc := make([]uint, position)
	slc[0] = 1
	slc[1] = 1

	if position <= 2 {
		return 1
	}

	var result, i uint
	for i = 2; i < position; i++ {
		result = slc[i-1] + slc[i-2]
		slc[i] = result
	}

	return result
}
