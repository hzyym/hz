package object

type InternalFun struct {
	Fun_ InsideFun
}

func (b *InternalFun) Type() ObjectType {
	return BUILTFun
}

func (b *InternalFun) Inspect() string {
	return "InternalFun"
}

type InternalName struct {
	Name string
	Fun  *InternalFun
}
