package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AreaOfInterestRANNodeItem struct {
	GlobalRANNodeID GlobalRANNodeID                                            `aper:"valueLB:0,valueUB:3"`
	IEExtensions    *ProtocolExtensionContainerAreaOfInterestRANNodeItemExtIEs `aper:"optional"`
}
