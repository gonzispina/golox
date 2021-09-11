package lox

import "fmt"

func isDigit(v rune) bool {
	return v >= '0' && v <= '9'
}

func isAlpha(v rune) bool {
	return v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' || v == '_'
}

func isAlphanumeric(v rune) bool {
	return isAlpha(v) || isDigit(v)
}

func report(line int, where, message string) {
	fmt.Printf("[line %v] Error %s: %s", line, where, message)
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func cast2String(value interface{}, t *Token) (string, error) {
	v, ok := value.(string)
	if !ok {
		return "", InvalidDataTypeError(t, getDataType(value), str)
	}
	return v, nil
}

func cast2Float(value interface{}, t *Token) (float64, error) {
	v, ok := value.(float64)
	if !ok {
		return 0, InvalidDataTypeError(t, getDataType(value), number)
	}
	return v, nil
}

func cast2Bool(value interface{}, t *Token) (bool, error) {
	v, ok := value.(bool)
	if !ok {
		return false, InvalidDataTypeError(t, getDataType(value), boolean)
	}
	return v, nil
}

func both2String(first, second interface{}, t *Token) (v1 string, v2 string, err error) {
	v1, err = cast2String(first, t)
	if err != nil {
		return "", "", err
	}

	v2, err = cast2String(second, t)
	if err != nil {
		return "", "", err
	}

	return v1, v2, nil
}

func both2Float(first, second interface{}, t *Token) (v1 float64, v2 float64, err error) {
	v1, err = cast2Float(first, t)
	if err != nil {
		return 0, 0, err
	}

	v2, err = cast2Float(second, t)
	if err != nil {
		return 0, 0, err
	}

	return v1, v2, nil
}

func both2Bool(first, second interface{}, t *Token) (v1 bool, v2 bool, err error) {
	v1, err = cast2Bool(first, t)
	if err != nil {
		return false, false, err
	}

	v2, err = cast2Bool(second, t)
	if err != nil {
		return false, false, err
	}

	return v1, v2, nil
}

func addValues(first, second interface{}, t *Token) (interface{}, error) {
	dt := getDataType(first)
	if dt == number {
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return nil, err
		}
		return v1 + v2, nil
	} else if dt == str {
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return nil, err
		}
		return v1 + v2, nil
	}
	return nil, InvalidOperationError(t, dt, getDataType(second))
}

func isEqual(first, second interface{}, t *Token) (bool, error) {
	dt := getDataType(first)
	switch dt {
	case number:
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 == v2, nil
	case str:
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 == v2, nil
	case boolean:
		v1, v2, err := both2Bool(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 == v2, nil
	default:
		return false, InvalidOperationError(t, dt, getDataType(second))
	}
}

func notEqual(first, second interface{}, t *Token) (bool, error) {
	dt := getDataType(first)
	switch dt {
	case number:
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 != v2, nil
	case str:
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 != v2, nil
	case boolean:
		v1, v2, err := both2Bool(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 != v2, nil
	default:
		return false, InvalidOperationError(t, dt, getDataType(second))
	}
}

func greaterEqual(first, second interface{}, t *Token) (bool, error) {
	dt := getDataType(first)
	switch dt {
	case number:
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 >= v2, nil
	case str:
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 >= v2, nil
	default:
		return false, InvalidOperationError(t, dt, getDataType(second))
	}
}

func lesserEqual(first, second interface{}, t *Token) (bool, error) {
	dt := getDataType(first)
	switch dt {
	case number:
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 <= v2, nil
	case str:
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 <= v2, nil
	default:
		return false, InvalidOperationError(t, dt, getDataType(second))
	}
}

func greaterThan(first, second interface{}, t *Token) (bool, error) {
	dt := getDataType(first)
	switch dt {
	case number:
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 > v2, nil
	case str:
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 > v2, nil
	default:
		return false, InvalidOperationError(t, dt, getDataType(second))
	}
}

func lesserThan(first, second interface{}, t *Token) (bool, error) {
	dt := getDataType(first)
	switch dt {
	case number:
		v1, v2, err := both2Float(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 < v2, nil
	case str:
		v1, v2, err := both2String(first, second, t)
		if err != nil {
			return false, err
		}
		return v1 < v2, nil
	default:
		return false, InvalidOperationError(t, dt, getDataType(second))
	}
}
