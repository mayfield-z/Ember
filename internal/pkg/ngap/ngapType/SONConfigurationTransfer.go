package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SONConfigurationTransfer struct {
	TargetRANNodeID        TargetRANNodeID                                           `aper:"valueExt"`
	SourceRANNodeID        SourceRANNodeID                                           `aper:"valueExt"`
	SONInformation         SONInformation                                            `aper:"valueLB:0,valueUB:2"`
	XnTNLConfigurationInfo *XnTNLConfigurationInfo                                   `aper:"valueExt,optional"`
	IEExtensions           *ProtocolExtensionContainerSONConfigurationTransferExtIEs `aper:"optional"`
}
