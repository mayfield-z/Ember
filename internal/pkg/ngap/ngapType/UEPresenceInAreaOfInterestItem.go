package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type UEPresenceInAreaOfInterestItem struct {
	LocationReportingReferenceID LocationReportingReferenceID
	UEPresence                   UEPresence
	IEExtensions                 *ProtocolExtensionContainerUEPresenceInAreaOfInterestItemExtIEs `aper:"optional"`
}
