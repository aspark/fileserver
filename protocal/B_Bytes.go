package protocal

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

//BytesBlock struct
type BytesBlock struct {
	Content []byte
}

func (blk BytesBlock) String() string {
	if blk.Content == nil {
		return "nil bytes array"
	}

	if len(blk.Content) > 10 {
		return fmt.Sprintf("%v", blk.Content[:10]) + "...len(" + strconv.Itoa(len(blk.Content)) + ")"
	}

	return fmt.Sprintf("%v", blk.Content)
}

//GetContent implement IBlock
func (blk BytesBlock) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk BytesBlock) GetBytes() []byte {
	return ConvertToBlockBytes(blk.Content)
}

//
func (blk BytesBlock) RawString() string {
	return fmt.Sprintf("%v", blk.Content)
}

//Parse bytes to block
func (blk BytesBlock) Parse(bytes []byte, rd io.Reader) (BlockParseResult, error) {
	var result = BlockParseResult{}

	if string(bytes[0]) != "!" {
		return result, errors.New("invalid format: " + string(bytes[0]))
	}

	// log.Println("begin parse array block")

	content, remain, err := readBytesByLength(bytes, rd)
	if err != nil {
		return result, err
	}

	result.Value = BytesBlock{Content: content}
	result.remainBytes = remain

	return result, nil
}
