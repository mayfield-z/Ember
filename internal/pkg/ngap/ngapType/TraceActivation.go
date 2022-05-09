package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TraceActivation struct {
	NGRANTraceID                   NGRANTraceID
	InterfacesToTrace              InterfacesToTrace
	TraceDepth                     TraceDepth
	TraceCollectionEntityIPAddress TransportLayerAddress
	IEExtensions                   *ProtocolExtensionContainerTraceActivationExtIEs `aper:"optional"`
}
