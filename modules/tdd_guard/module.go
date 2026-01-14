package tddguard

var Module = &TddGuardModule{
	Name:     "tdd-guard",
	Version:  1,
	Category: "claude",
}

type TddGuardModule struct {
	Name     string
	Version  int
	Category string
}

func (tm *TddGuardModule) GetName() string {
	return tm.Name
}

func (tm *TddGuardModule) GetCategory() string {
	return tm.Category
}
func (tm *TddGuardModule) GetVersion() int {
	return tm.Version
}
func (tm *TddGuardModule) IsInstalled() bool {
	return false
}

func (tm *TddGuardModule) Install() bool {
	return false
}

func (tm *TddGuardModule) Uninstall() bool {
	return false
}
