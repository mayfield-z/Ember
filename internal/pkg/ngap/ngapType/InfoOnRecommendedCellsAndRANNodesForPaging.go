package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type InfoOnRecommendedCellsAndRANNodesForPaging struct {
	RecommendedCellsForPaging  RecommendedCellsForPaging                                                   `aper:"valueExt"`
	RecommendRANNodesForPaging RecommendedRANNodesForPaging                                                `aper:"valueExt"`
	IEExtensions               *ProtocolExtensionContainerInfoOnRecommendedCellsAndRANNodesForPagingExtIEs `aper:"optional"`
}
