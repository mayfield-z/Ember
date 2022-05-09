package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct UPTransportLayerInformationList */
/* UPTransportLayerInformationItem */
type UPTransportLayerInformationList struct {
	List []UPTransportLayerInformationItem `aper:"valueExt,sizeLB:1,sizeUB:3"`
}
