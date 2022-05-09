package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayfield-z/ember/internal/pkg/nas/nasMessage"
	"github.com/mayfield-z/ember/internal/pkg/nas/nasType"
)

type nasTypeResultMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeResultMessageIdentityTable = []nasTypeResultMessageIdentityData{
	{nasMessage.PDUSessionAuthenticationResultEAPMessageType, nasMessage.PDUSessionAuthenticationResultEAPMessageType},
}

func TestNasTypeNewAuthenticationResultMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationResultMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetAuthenticationResultMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationResultMessageIdentity()
	for _, table := range nasTypeResultMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
