package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AMF_TNLAssociationToUpdateList */
/* AMFTNLAssociationToUpdateItem */
type AMFTNLAssociationToUpdateList struct {
	List []AMFTNLAssociationToUpdateItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}
