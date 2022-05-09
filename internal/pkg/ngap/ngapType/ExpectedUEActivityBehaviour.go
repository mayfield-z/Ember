package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type ExpectedUEActivityBehaviour struct {
	ExpectedActivityPeriod                 *ExpectedActivityPeriod                                      `aper:"optional"`
	ExpectedIdlePeriod                     *ExpectedIdlePeriod                                          `aper:"optional"`
	SourceOfUEActivityBehaviourInformation *SourceOfUEActivityBehaviourInformation                      `aper:"optional"`
	IEExtensions                           *ProtocolExtensionContainerExpectedUEActivityBehaviourExtIEs `aper:"optional"`
}
