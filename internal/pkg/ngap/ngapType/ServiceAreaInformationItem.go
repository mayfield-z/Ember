package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type ServiceAreaInformationItem struct {
	PLMNIdentity   PLMNIdentity
	AllowedTACs    *AllowedTACs                                                `aper:"optional"`
	NotAllowedTACs *NotAllowedTACs                                             `aper:"optional"`
	IEExtensions   *ProtocolExtensionContainerServiceAreaInformationItemExtIEs `aper:"optional"`
}
