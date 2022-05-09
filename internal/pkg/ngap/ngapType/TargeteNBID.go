package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TargeteNBID struct {
	GlobalENBID    GlobalNgENBID                                `aper:"valueExt"`
	SelectedEPSTAI EPSTAI                                       `aper:"valueExt"`
	IEExtensions   *ProtocolExtensionContainerTargeteNBIDExtIEs `aper:"optional"`
}
