package test

import (
	"testing"

	"log"

	"github.com/aspark/fileserver/protocal"
)

// type EmptyReader struct{}

// func (r EmptyReader) Read(p []byte) (n int, err error) {
// 	return 0, errors.New("not implement")
// }

func Test_SimpleStringParser(t *testing.T) {
	// var bytes = []byte("+abcd\r\n+1234\r\n")
	str1, str2 := "abcd", "1234"
	var bytes = protocal.ConvertToSimpleStringBytes(str1)
	bytes = append(bytes, protocal.ConvertToSimpleStringBytes(str2)...)

	block, err := protocal.ParseOne(bytes)

	if err == nil && block.GetContent() == str1 {
		t.Log("SimpleString ParseOne Passed")
	} else {
		t.Error(err, block.GetContent())
	}

	blocks, err := protocal.ParseAll(bytes, nil)

	if err == nil && len(blocks) == 2 && blocks[1].GetContent() == str2 {
		t.Log("SimpleString ParseAll Passed")
	} else {
		t.Error(err, len(blocks), blocks[1].GetContent())
	}
}

func Test_BulkStringParser(t *testing.T) {

	str := "12abcde547d1\ra"
	bytes := protocal.ConvertToBulkStringBytes(&str)

	block, err := protocal.ParseOne(bytes)

	if err == nil && block.GetContent() == str {
		t.Log("BulkString ParseOne Passed")
	} else {
		t.Error(err, block.GetContent())
	}

	str1 := ""
	block, err = protocal.ParseOne(protocal.ConvertToBulkStringBytes(&str1))

	if err == nil && block.GetContent() == "" {
		t.Log("SimpleString Empty ParseOne Passed")
	} else {
		t.Error(err, block.GetContent())
	}

	block, err = protocal.ParseOne(protocal.ConvertToBulkStringBytes(nil))

	if err == nil && block.GetContent() == nil {
		t.Log("SimpleString nil ParseOne Passed")
	} else {
		t.Error(err, block.GetContent())
	}

}

func Test_BytesParser(t *testing.T) {

	orignial := []byte{0, 1, 2, 3, 4, 5, 6, 0}
	bytes := protocal.ConvertToBlockBytes(orignial)

	block, err := protocal.ParseOne(bytes)

	receive, ok := block.GetContent().([]byte)
	var originalLength = len(orignial)
	if ok && err == nil &&
		len(receive) == originalLength &&
		receive[0] == orignial[0] &&
		receive[originalLength-1] == orignial[originalLength-1] {

		t.Log("BytesParser ParseOne Passed")
	} else {
		t.Error(err, block.GetContent())
	}

}

func Test_IntegerParser(t *testing.T) {
	var i = 987654321
	var bytes = protocal.ConvertToIntegerBytes(i)

	block, err := protocal.ParseOne(bytes)
	n, ok := block.GetContent().(int)
	if ok && err == nil && n == i {
		t.Log("IntegerParser ParseOne Passed")
	} else {
		t.Error(err, block.GetContent())
	}
}

func Test_ErrorConverter(t *testing.T) {
	var str = "error occur"
	var bytes = protocal.ConvertToErrorBytes(str)

	if len(bytes) == len(str)+3 && bytes[0] == '-' {
		t.Log("ErrorConverter ParseOne Passed")
	} else {
		t.Error(string(bytes))
	}
}

func Test_ArrayParse(t *testing.T) {
	var bs = protocal.ConvertToArrayBytes_ForString("a", "b", "c", "d")

	blk, err := protocal.ParseOne(bs)

	var content = blk.GetContent()
	if err == nil && len(content.([]protocal.IBlock)) == 4 {
		log.Println("ArrayParse Passed")
	} else {
		log.Fatal(blk, err)
	}

}

func Test_CompositParse(t *testing.T) {
	str1, str2 := "abcd", "12\t34"
	i := 90
	bs := []byte{0, 1, 2, 3, 4, 5}
	var bytes = protocal.ConvertToSimpleStringBytes(str1)

	bytes = append(bytes, protocal.ConvertToBulkStringBytes(&str2)...)

	bytes = append(bytes, protocal.ConvertToIntegerBytes(i)...)

	bytes = append(bytes, protocal.ConvertToBlockBytes(bs)...)

	blocks, err := protocal.ParseAll(bytes, nil)

	if err == nil && len(blocks) == 4 &&
		blocks[0].GetContent() == str1 &&
		blocks[1].GetContent() == str2 &&
		blocks[2].GetContent() == i {
		nbs, ok := blocks[3].GetContent().([]byte)
		if ok && len(nbs) == len(bs) {
			t.Log("CompositParse ParseAll Passed")
		}
	} else {
		t.Error(err, len(blocks), blocks)
	}
}
