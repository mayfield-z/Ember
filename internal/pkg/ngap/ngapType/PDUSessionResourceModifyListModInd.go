package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceModifyListModInd */
/* PDUSessionResourceModifyItemModInd */
type PDUSessionResourceModifyListModInd struct {
	List []PDUSessionResourceModifyItemModInd `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
