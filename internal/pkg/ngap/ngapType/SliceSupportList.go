package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct SliceSupportList */
/* SliceSupportItem */
type SliceSupportList struct {
	List []SliceSupportItem `aper:"valueExt,sizeLB:1,sizeUB:1024"`
}