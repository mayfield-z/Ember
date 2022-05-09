package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TargetNGRANNodeToSourceNGRANNodeTransparentContainer struct {
	RRCContainer RRCContainer
	IEExtensions *ProtocolExtensionContainerTargetNGRANNodeToSourceNGRANNodeTransparentContainerExtIEs `aper:"optional"`
}
