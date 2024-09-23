// utils/globalvar.go
package utils

// GlobalVar struct dengan private variabel ReferenceId dan version
type GlobalVar struct {
	ReferenceId string
}

// Constructor untuk GlobalVar
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
