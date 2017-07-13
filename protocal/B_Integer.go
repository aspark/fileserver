package protocal

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

//IntegerBlock struct
type IntegerBlock struct {
	Content int
}

//Integer64Block struct
type Integer64Block struct {
	Content int64
}

//
func (blk IntegerBlock) String() string {
	return fmt.Sprintf("%d", blk.Content)
}

//GetContent implement IBlock
func (blk IntegerBlock) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk IntegerBlock) GetBytes() []byte {
	return ConvertToIntegerBytes(blk.Content)
}

//
func (blk IntegerBlock) RawString() string {
	return fmt.Sprintf("%d", blk.Content)
}

//
func (blk Integer64Block) String() string {
	return fmt.Sprintf("%d", blk.Content)
}

//GetContent implement IBlock
func (blk Integer64Block) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk Integer64Block) GetBytes() []byte {
	return ConvertToInt64Bytes(blk.Content)
}

//
func (blk Integer64Block) RawString() string {
	return fmt.Sprintf("%d", blk.Content)
}

//Parse bytes to integer
func (blk IntegerBlock) Parse(bytes []byte, rd io.Reader) (BlockParseResult, error) {
	var result = BlockParseResult{}

	if string(bytes[0]) != ":" {
		return result, errors.New("invalid format: " + string(bytes[0]))
	}

	line, bytes, err := readLine(bytes, rd)
	if err != nil {
		return result, err
	}

	i, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return result, err
	}

	result.Value = IntegerBlock{Content: i}
	result.remainBytes = bytes

	return result, nil
}
