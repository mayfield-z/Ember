package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct XnGTP_TLAs */
/* TransportLayerAddress */
type XnGTPTLAs struct {
	List []TransportLayerAddress `aper:"sizeLB:1,sizeUB:16"`
}
