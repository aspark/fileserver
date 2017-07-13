package test

import (
	"testing"

	"github.com/aspark/fileserver/command"
	"github.com/aspark/fileserver/config"
	"github.com/aspark/fileserver/protocal"
)

func execute(t *testing.T, bytes []byte) ([]protocal.IBlock, error) {
	// var block protocal.Block
	blocks, err := protocal.ParseAll(bytes, nil)
	if err != nil {
		t.Fatal(err)
		return nil, err
	}

	ctx := config.CreateContext(nil)
	ctx.HasAuth = true

	result, err := command.ExecuteBlocks(ctx, blocks)
	if err != nil {
		t.Fatal(err)
		return nil, err
	}

	results, err := protocal.ParseAll(result, nil)

	return results, err
}

func Test_ExistCommand(t *testing.T) {
	// var cmd = "$6\r\nEXISTS\r\n$12\r\nfileFullName\r\n"
	var bytes = protocal.ConvertToBulkStringBytesEx("EXISTS", "test.txt")
	blk, err := execute(t, bytes)
	// if err != nil {
	// 	return
	// }

	if blk[0].GetContent() == 0 {
		t.Log("Exist Passed")
	} else {
		t.Fatal(err)
	}

	bytes = protocal.ConvertToBulkStringBytesEx("EXISTS", "../../attachments/test.txt")
	blk, err = execute(t, bytes)

	if blk[0].GetContent() == 1 {
		t.Log("Exist Passed")
	} else {
		t.Fatal(err)
	}
}

func Test_SaveCommand(t *testing.T) {
	fileName := "../../attachments/save.test"
	var bytes = protocal.ConvertToBulkStringBytesEx("save", fileName)
	bytes = append(bytes, protocal.ConvertToBlockBytes([]byte{'a', 'b', 'c'})...)
	blk, err := execute(t, bytes)

	if blk[0].GetContent() == "OK" {
		bytes = protocal.ConvertToBulkStringBytesEx("EXISTS", fileName)
		blk, err = execute(t, bytes)

		if blk[0].GetContent() == 1 {
			DelFile(t, fileName)

			t.Log("Save Passed")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}

func DelFile(t *testing.T, fileName string) {
	var bytes = protocal.ConvertToBulkStringBytesEx("del", fileName)
	blk, err := execute(t, bytes)

	if err != nil {
		t.Fatal(err)
	}

	if blk[0].GetContent() == "OK" {
		bytes = protocal.ConvertToBulkStringBytesEx("EXISTS", fileName)
		blk, err = execute(t, bytes)

		if blk[0].GetContent() == 0 {
			t.Log("del Passed")
		} else {
			t.Fatal("del file failed")
		}
	} else {
		t.Fatal("del file failed")
	}
}

func Test_LenCommand(t *testing.T) {
	var bytes = protocal.ConvertToBulkStringBytesEx("LEN", "../../attachments/sub/test.txt")
	blk, err := execute(t, bytes)

	if blk[0].GetContent() == 27 {
		t.Log("Len Passed")
	} else {
		t.Fatal(err)
	}

}

func Test_GetCommand(t *testing.T) {
	var bytes = protocal.ConvertToBulkStringBytesEx("GET", "../../attachments/sub/test.txt")
	blk, err := execute(t, bytes)

	contents, ok := blk[0].GetContent().([]byte)
	if ok && len(contents) == 27 && contents[0] == 't' {
		t.Log("Len Passed")
	} else {
		t.Fatal(err)
	}
}

func Test_CGetCommand(t *testing.T) {
	var bytes = protocal.CombineBlockBytes(
		protocal.ConvertToBulkStringBytesEx("CGET", "../../attachments/sub/test.txt"),
		protocal.ConvertToIntegerBytes(1),
		protocal.ConvertToIntegerBytes(3),
	)
	blk, err := execute(t, bytes)

	results := blk[0].GetContent().([]protocal.IBlock)

	length, offsize, realLength := results[0].GetContent().(int), results[1].GetContent().(int), results[2].GetContent().(int)
	contents, ok := results[3].GetContent().([]byte)
	if ok && length == 27 && offsize == 1 && realLength == 3 && len(contents) == 3 && contents[0] == 'h' {
		t.Log("Len Passed")
	} else {
		t.Fatal(err)
	}
}

func Test_PingCommand(t *testing.T) {
	var bytes = protocal.ConvertToBulkStringBytesEx("PING", "lalala")
	blk, err := execute(t, bytes)

	if blk[0].GetContent() == "PONG" {
		t.Log("Len Passed")
	} else {
		t.Fatal(err)
	}
}

func Test_DirsCommand(t *testing.T) {
	var bytes = protocal.ConvertToBulkStringBytesEx("DIR", "../../attachments")
	blk, err := execute(t, bytes)

	var items = blk[0].GetContent().([]protocal.IBlock)
	if len(items) == 2 && items[1].GetContent() == "test.txt" {
		t.Log("Len Passed")
	} else {
		t.Fatal(blk, err)
	}
}
