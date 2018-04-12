package transaction

type Output interface {
	GetTransaction() *Transaction
}

type Transaction interface {
	GetKey() *Key
	SetKey(*Key)
	SetBlock(*Block)
}

type Key interface {

}

type Block interface {

}
