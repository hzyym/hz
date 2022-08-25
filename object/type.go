package object

type ObjectType int

const (
	INT ObjectType = iota
	BOOL
	NULL
	RETURN
	ERROR
	FUN
	STRING
	BUILTFun
	ARRAY
	HASH
	CompiledFun
)

var typeString = map[ObjectType]string{
	INT:    "int",
	BOOL:   "bool",
	NULL:   "null",
	STRING: "string",
	ARRAY:  "array",
	HASH:   "hash",
}

func (o ObjectType) String() string {
	return typeString[o]
}

type InsideFun func(ages ...Object) Object
