package protocal

import (
	"errors"
	"io"
	"strconv"
)

//IBlockParser block praser common base
type IBlockParser interface {
	//parse bytes to Block, return true if complete parse
	Parse(bytes []byte, rd io.Reader) (BlockParseResult, error)

	// Convert(content interface{}) []byte
}

//
type BlockParseResult struct {
	Value IBlock
	// Position    int
	remainBytes []byte
}

//get block parser mapper
func GetBlockPraserMapper() map[rune]IBlockParser {
	mapper := map[rune]IBlockParser{
		'+': SimpleStringBlock{},
		'-': ErrorBlock{},
		':': IntegerBlock{},
		'$': BulkStringBlock{},
		'*': ArrayBlock{},
		'!': BytesBlock{},
	}

	return mapper
}

var parserMapper = GetBlockPraserMapper()

//GetBlockParser method
func GetBlockParser(start rune) (IBlockParser, error) {
	parser, ok := parserMapper[start]
	if !ok {
		return nil, errors.New("not support prefix: " + string(start))
	}

	return parser, nil
}

//ParseOne Parse bytes and get the first block
func ParseOne(bytes []byte) (IBlock, error) {
	if bytes == nil || len(bytes) == 0 {
		return nil, errors.New("bytes can not be nil or empty")
	}

	parser, err := GetBlockParser(rune(bytes[0]))
	if err != nil {
		return nil, err
	}

	result, err := parser.Parse(bytes, nil)

	if err != nil {
		return result.Value, err
	}

	return result.Value, nil
}

//ParseAll Parse bytes to block array
//blocks, remainBytes, error
func ParseAll(bytes []byte, rd io.Reader) ([]IBlock, error) { //, []byte
	results, err := parseAllBlockResults(bytes, rd)
	if err != nil {
		return nil, err
	}

	blocks, _ := getBlockFromResults(results)

	return blocks, nil
}

func getBlockFromResults(results []BlockParseResult) ([]IBlock, []byte) {
	var blocks = make([]IBlock, 0, 8)

	var remain []byte
	for _, result := range results {
		remain = result.remainBytes
		blocks = append(blocks, result.Value)
	}

	return blocks, remain
}

func parseAllBlockResults(bytes []byte, rd io.Reader) ([]BlockParseResult, error) {
	var results = make([]BlockParseResult, 0, 8)

	if bytes == nil || len(bytes) == 0 {
		bytes = make([]byte, 1)
		_, err := rd.Read(bytes)
		if err != nil {
			return nil, err
		}
	}

	var data = bytes[:]

	for {
		parser, err := GetBlockParser(rune(data[0]))
		if err != nil {
			return nil, err
		}

		result, err := parser.Parse(data, rd)

		if err != nil {
			return nil, err
		}

		results = append(results, result)

		if len(result.remainBytes) == 0 {
			break
		}
		data = result.remainBytes
	}

	return results, nil
}

//line bytes, remain bytes, error
func readLine(bytes []byte, rd io.Reader) ([]byte, []byte, error) {
	end, length := -1, len(bytes)
	for i, b := range bytes {
		if b == '\r' && i+1 < length && bytes[i+1] == '\n' { //eof
			end = i
			break
		}
	}

	if end > 0 {
		line, newBytes := bytes[:end], make([]byte, 0)
		if end+2 < length {
			newBytes = bytes[end+2:]
		}

		return line, newBytes, nil
	}

	if rd == nil {
		return nil, nil, errors.New("reader is null")
	}

	buff := make([]byte, 1024)
	readN, err := rd.Read(buff)
	if err != nil {
		return nil, nil, err
	}

	return readLine(append(bytes, buff[0:readN]...), rd)
}

//read bytes, remain bytes, error
func readBytesByLength(bytes []byte, rd io.Reader) ([]byte, []byte, error) {
	line, remain, err := readLine(bytes, rd)
	if err != nil {
		return nil, nil, err
	}

	bytesCount, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, nil, err
	}

	if bytesCount < 0 {
		return nil, nil, err
	}

	if len(remain) >= bytesCount+2 {
		left := make([]byte, 0)
		if len(remain) > bytesCount+2 {
			left = remain[bytesCount+2:]
		}
		return remain[:bytesCount], left, nil // //BlockBytes: bytes[:end],
	}

	if rd == nil {
		return nil, nil, errors.New("reader is null")
	}

	buff := make([]byte, bytesCount)
	readN, rerr := rd.Read(buff)
	if rerr != nil {
		return nil, nil, rerr
	}

	// log.Println("readed bytes count: ", readN)

	return readBytesByLength(append(bytes, buff[:readN]...), rd)
}
