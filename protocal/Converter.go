package protocal

import (
	"strconv"
)

//
func CombineBlockBytes(items ...[]byte) []byte {
	bytes := make([]byte, 0, 8)
	for _, item := range items {
		bytes = append(bytes, item...)
	}

	return bytes
}

//Convert string to bytes
func ConvertToSimpleStringBytes(content string) []byte {
	bytes := []byte{'+'}

	bytes = append(bytes, []byte(content)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//Convert string to bytes
func ConvertToErrorBytes(content string) []byte {
	bytes := []byte{'-'}

	bytes = append(bytes, []byte(content)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//
func ConvertToErrorBytes_Error(content error) []byte {
	bytes := []byte{'-'}

	bytes = append(bytes, []byte(content.Error())...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//Convert Bulk String Bytes
func ConvertToBulkStringBytes(content *string) []byte {
	if content == nil { // can be null
		return []byte{'$', '-', '1', '\r', '\n'}
	}

	bytes := append([]byte("$"), strconv.Itoa(len(*content))...)
	bytes = append(bytes, '\r', '\n')

	bytes = append(bytes, []byte(*content)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//
func ConvertToBulkStringBytesEx(contents ...string) []byte {
	var bytes = make([]byte, 0, 8)
	for _, content := range contents {
		bytes = append(bytes, ConvertToBulkStringBytes(&content)...)
	}

	return bytes
}

//
func ConvertToBlockBytes(content []byte) []byte {
	if content == nil {
		return []byte{'$', '-', '1', '\r', '\n'}
	}
	bytes := append([]byte("!"), strconv.Itoa(len(content))...)
	bytes = append(bytes, '\r', '\n')

	bytes = append(bytes, content...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//
func ConvertToIntegerBytes(content int) []byte {
	str := strconv.Itoa(content)

	bytes := append([]byte{':'}, []byte(str)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//
func ConvertToInt64Bytes(content int64) []byte {
	str := strconv.FormatInt(content, 10)

	bytes := append([]byte{':'}, []byte(str)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//
func ConvertToArrayBytes_ForString(items ...string) []byte {
	bytes := []byte{'*'}
	if items == nil {
		bytes = append(bytes, '-', '1', '\r', '\n')
	} else {
		bytes = append(bytes, []byte(strconv.Itoa(len(items)))...)
		bytes = append(bytes, '\r', '\n')

		for _, c := range items {
			bytes = append(bytes, ConvertToBulkStringBytes(&c)...)
		}
	}

	return bytes
}

//
func ConvertToArrayBytes(items []IBlock) []byte {
	bytes := []byte{'*'}
	if items == nil {
		bytes = append(bytes, '-', '1', '\r', '\n')
	} else {
		bytes = append(bytes, []byte(strconv.Itoa(len(items)))...)
		bytes = append(bytes, '\r', '\n')

		for _, c := range items {
			bytes = append(bytes, c.GetBytes()...)
		}
	}

	return bytes
}
