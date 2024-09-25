// utils/globalvar.go
package utils

// GlobalVar struct
type GlobalVar struct {
	ReferenceId string
}

// Constructor for GlobalVar
func NewGlobalVar() *GlobalVar {
	return &GlobalVar{
		ReferenceId: "",
	}
}

// Getter setter
func (g *GlobalVar) GetReferenceId() string {
	return g.ReferenceId
}

func (g *GlobalVar) SetReferenceId(ReferenceId string) {
	g.ReferenceId = ReferenceId
}

var GlobalVarInstance = NewGlobalVar()
