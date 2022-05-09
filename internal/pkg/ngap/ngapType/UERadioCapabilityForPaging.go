package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type UERadioCapabilityForPaging struct {
	UERadioCapabilityForPagingOfNR    *UERadioCapabilityForPagingOfNR                             `aper:"optional"`
	UERadioCapabilityForPagingOfEUTRA *UERadioCapabilityForPagingOfEUTRA                          `aper:"optional"`
	IEExtensions                      *ProtocolExtensionContainerUERadioCapabilityForPagingExtIEs `aper:"optional"`
}
