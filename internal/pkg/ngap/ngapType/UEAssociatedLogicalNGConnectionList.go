package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct UE_associatedLogicalNG_connectionList */
/* UEAssociatedLogicalNGConnectionItem */
type UEAssociatedLogicalNGConnectionList struct {
	List []UEAssociatedLogicalNGConnectionItem `aper:"valueExt,sizeLB:1,sizeUB:65536"`
}