package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceListHORqd */
/* PDUSessionResourceItemHORqd */
type PDUSessionResourceListHORqd struct {
	List []PDUSessionResourceItemHORqd `aper:"valueExt,sizeLB:1,sizeUB:256"`
}