package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AdditionalDLUPTNLInformationForHOList */
/* AdditionalDLUPTNLInformationForHOItem */
type AdditionalDLUPTNLInformationForHOList struct {
	List []AdditionalDLUPTNLInformationForHOItem `aper:"valueExt,sizeLB:1,sizeUB:3"`
}