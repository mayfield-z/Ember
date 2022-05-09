package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type ERABInformationItem struct {
	ERABID       ERABID
	DLForwarding *DLForwarding                                        `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerERABInformationItemExtIEs `aper:"optional"`
}
