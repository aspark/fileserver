package protocal

import (
	"io"
	"strconv"
)

//BulkStringBlock struct
type BulkStringBlock struct {
	Content interface{} //string
}

func (blk BulkStringBlock) String() string {
	if blk.Content == nil {
		return "nil bulk string"
	}

	var content = blk.Content.(string)
	if len(content) > 30 {
		return content[:30] + "...len(" + strconv.Itoa(len(content)) + ")"
	}

	return content
}

//GetContent implement IBlock
func (blk BulkStringBlock) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk BulkStringBlock) GetBytes() []byte {
	var str, ok = blk.Content.(string)
	if ok {
		return ConvertToBulkStringBytes(&str)
	}

	return ConvertToBulkStringBytes(nil)
}

//
func (blk BulkStringBlock) RawString() string {
	if blk.Content == nil {
		return ""
	}

	return blk.Content.(string)
}

//Parse BulkString
func (blk BulkStringBlock) Parse(bytes []byte, rd io.Reader) (BlockParseResult, error) {
	var result = BlockParseResult{}

	content, reamin, err := readBytesByLength(bytes, rd)
	if err != nil {
		return result, err
	}

	if content != nil {
		result.Value = BulkStringBlock{Content: string(content)}
	} else {
		result.Value = BulkStringBlock{Content: nil}
	}

	result.remainBytes = reamin

	return result, nil
}
