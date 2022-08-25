package compiler

const (
	Global SymbolType = iota
	Local
	Free
)

type SymbolType int
type SymbolTable struct {
	top   *SymbolTable
	table map[string]*Symbol
	index int
	free  []*Symbol
}
type Symbol struct {
	index int
	Name  string
	types SymbolType
}

func NewSymbolTable(top *SymbolTable) *SymbolTable {
	return &SymbolTable{table: map[string]*Symbol{}, top: top}
}
func (s *SymbolTable) SetSymbol(name string) *Symbol {
	symbol := &Symbol{Name: name, index: s.index}
	if s.top == nil {
		symbol.types = Global
	} else {
		symbol.types = Local
	}
	s.table[name] = symbol
	s.index++
	return symbol
}
func (s *SymbolTable) GetSymbol(name string) (*Symbol, bool) {
	v, ok := s.table[name]
	if !ok && s.top != nil {
		v, ok = s.top.GetSymbol(name)
		if !ok {
			return v, ok
		}
		if v.types == Global {
			return v, ok
		}
		return s.setFreeSymbol(v), true
	}
	return v, ok
}
func (s *SymbolTable) setFreeSymbol(symbol *Symbol) *Symbol {
	s.free = append(s.free, symbol)
	sy := &Symbol{
		index: len(s.free) - 1,
		Name:  symbol.Name,
		types: Free,
	}
	s.table[symbol.Name] = sy
	return sy
}
func (s *SymbolTable) DelSymbol(name string) {
	delete(s.table, name)
}
