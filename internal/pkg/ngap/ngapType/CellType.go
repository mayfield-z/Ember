package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type CellType struct {
	CellSize     CellSize
	IEExtensions *ProtocolExtensionContainerCellTypeExtIEs `aper:"optional"`
}
