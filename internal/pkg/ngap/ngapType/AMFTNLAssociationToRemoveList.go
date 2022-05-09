package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AMF_TNLAssociationToRemoveList */
/* AMFTNLAssociationToRemoveItem */
type AMFTNLAssociationToRemoveList struct {
	List []AMFTNLAssociationToRemoveItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}
