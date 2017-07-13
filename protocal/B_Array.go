package protocal

import (
	"errors"
	"io"
	"strconv"
)

//
type ArrayBlock struct {
	// Block
	Content []IBlock
}

//GetContent implement IBlock
func (blk ArrayBlock) GetContent() interface{} {
	return blk.Content
}

//GetBytes implement IBlockResult
func (blk ArrayBlock) GetBytes() []byte {
	return ConvertToArrayBytes(blk.Content)
}

//
func (blk ArrayBlock) RawString() string {
	var str string

	for _, item := range blk.Content {
		str += (item.RawString() + "\r\n")
	}

	return str
}

//Parse array
func (blk ArrayBlock) Parse(bytes []byte, rd io.Reader) (BlockParseResult, error) {
	var result = BlockParseResult{}

	line, remain, err := readLine(bytes, rd)

	if err != nil {
		return result, err
	}

	count, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return result, err
	}

	// log.Println("array items count: ", count)

	if count == 0 {
		result.Value = ArrayBlock{Content: make([]IBlock, 0)}
	} else {
		results, err := parseAllBlockResults(remain, rd)
		if err != nil {
			return result, err
		}

		if len(results) != count {
			return result, errors.New("array element count mismatched")
		}

		resultBytes, resultRemain := getBlockFromResults(results)
		result.Value = ArrayBlock{Content: resultBytes}
		result.remainBytes = resultRemain
		// result.remainBytes = results[len(results)-1].remainBytes
	}

	return result, nil
}
