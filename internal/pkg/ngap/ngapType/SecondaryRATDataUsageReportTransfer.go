package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SecondaryRATDataUsageReportTransfer struct {
	SecondaryRATUsageInformation *SecondaryRATUsageInformation                                        `aper:"valueExt,optional"`
	IEExtensions                 *ProtocolExtensionContainerSecondaryRATDataUsageReportTransferExtIEs `aper:"optional"`
}
