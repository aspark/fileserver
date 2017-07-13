package protocal

import (
	"errors"
	"io"
)

//ErrorBlock struct
type ErrorBlock struct {
	Content string
}

//GetContent implement IBlock
func (blk ErrorBlock) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk ErrorBlock) GetBytes() []byte {
	return ConvertToErrorBytes(blk.Content)
}

//
func (blk ErrorBlock) RawString() string {
	return blk.Content
}

//
func (blk ErrorBlock) Parse(bytes []byte, rd io.Reader) (BlockParseResult, error) {
	var result = BlockParseResult{}

	if string(bytes[0]) != "-" {
		return result, errors.New("invalid format: " + string(bytes[0]))
	}

	line, bytes, err := readLine(bytes, rd)
	if err != nil {
		return result, nil
	}

	result.Value = ErrorBlock{Content: string(line[1:])}
	result.remainBytes = bytes

	return result, nil
}
