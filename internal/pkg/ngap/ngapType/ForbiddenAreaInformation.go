package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ForbiddenAreaInformation */
/* ForbiddenAreaInformationItem */
type ForbiddenAreaInformation struct {
	List []ForbiddenAreaInformationItem `aper:"valueExt,sizeLB:1,sizeUB:16"`
}
