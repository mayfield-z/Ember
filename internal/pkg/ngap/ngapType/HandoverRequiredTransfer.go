package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type HandoverRequiredTransfer struct {
	DirectForwardingPathAvailability *DirectForwardingPathAvailability                         `aper:"optional"`
	IEExtensions                     *ProtocolExtensionContainerHandoverRequiredTransferExtIEs `aper:"optional"`
}
