package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PagingAttemptInformation struct {
	PagingAttemptCount             PagingAttemptCount
	IntendedNumberOfPagingAttempts IntendedNumberOfPagingAttempts
	NextPagingAreaScope            *NextPagingAreaScope                                      `aper:"optional"`
	IEExtensions                   *ProtocolExtensionContainerPagingAttemptInformationExtIEs `aper:"optional"`
}
