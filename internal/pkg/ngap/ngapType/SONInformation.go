package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	SONInformationPresentNothing int = iota /* No components present */
	SONInformationPresentSONInformationRequest
	SONInformationPresentSONInformationReply
	SONInformationPresentChoiceExtensions
)

type SONInformation struct {
	Present               int
	SONInformationRequest *SONInformationRequest
	SONInformationReply   *SONInformationReply `aper:"valueExt"`
	ChoiceExtensions      *ProtocolIESingleContainerSONInformationExtIEs
}
