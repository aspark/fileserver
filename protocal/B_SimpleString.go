package protocal

import (
	"errors"
	"io"
)

//SimpleStringBlock struct
type SimpleStringBlock struct {
	Content string
}

//GetContent implement IBlock
func (blk SimpleStringBlock) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk SimpleStringBlock) GetBytes() []byte {
	return ConvertToSimpleStringBytes(blk.Content)
}

//
func (blk SimpleStringBlock) RawString() string {
	return blk.Content
}

//Parse SimpleStringBlock
func (blk SimpleStringBlock) Parse(bytes []byte, rd io.Reader) (BlockParseResult, error) {
	var result = BlockParseResult{}
	// var block Block

	if string(bytes[0]) != "+" {
		return result, errors.New("invalid format: " + string(bytes[0]))
	}

	line, bytes, err := readLine(bytes, rd)
	if err != nil {
		return result, nil
	}

	result.Value = SimpleStringBlock{Content: string(line[1:])}
	result.remainBytes = bytes

	return result, nil
}
