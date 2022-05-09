package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SONInformationReply struct {
	XnTNLConfigurationInfo *XnTNLConfigurationInfo                              `aper:"valueExt,optional"`
	IEExtensions           *ProtocolExtensionContainerSONInformationReplyExtIEs `aper:"optional"`
}
