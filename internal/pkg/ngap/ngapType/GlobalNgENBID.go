package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type GlobalNgENBID struct {
	PLMNIdentity PLMNIdentity
	NgENBID      NgENBID                                        `aper:"valueLB:0,valueUB:3"`
	IEExtensions *ProtocolExtensionContainerGlobalNgENBIDExtIEs `aper:"optional"`
}
