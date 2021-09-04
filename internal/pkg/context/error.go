package context

import "fmt"

type UERegisterError struct {
	UeSUPI  string
	GNBName string
}

func (e UERegisterError) Error() string {
	return fmt.Sprintf("UE register failed, supi: %v, gNB: %v", e.UeSUPI, e.GNBName)
}

func (u *UE) UeRegisterError() error {
	return UERegisterError{
		UeSUPI:  u.supi,
		GNBName: u.connectedGnb.name,
	}
}
