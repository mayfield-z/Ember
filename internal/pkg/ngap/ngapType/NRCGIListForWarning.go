package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct NR_CGIListForWarning */
/* NRCGI */
type NRCGIListForWarning struct {
	List []NRCGI `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
