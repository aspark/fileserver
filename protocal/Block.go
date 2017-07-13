package protocal

//
type IBlock interface {
	GetContent() interface{}
	GetBytes() []byte
	RawString() string
}

//
type IBlockResult interface {
	GetBytes() []byte
}

//
// type Block struct {
// 	// BlockBytes []byte
// 	Content interface{}
// }

// //
// func (blk Block) GetContent() interface{} {
// 	return blk.Content
// }
