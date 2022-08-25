package object

var table = []*InternalName{
	{Name: "len", Fun: &InternalFun{Fun_: Len}},
	{Name: "println", Fun: &InternalFun{Fun_: Println}},
	{Name: "echo", Fun: &InternalFun{Fun_: Echo}},
	{Name: "put", Fun: &InternalFun{Fun_: Put}},
	{Name: "str_rev", Fun: &InternalFun{Fun_: StringReversal}},
}

var funName map[string]int

func init() {
	funName = make(map[string]int)
	for index, list := range table {
		funName[list.Name] = index
	}
}
func GetNameIndex(name string) (int, bool) {
	v, ok := funName[name]
	return v, ok
}
func GetFun(index int) *InternalFun {
	return table[index].Fun
}
