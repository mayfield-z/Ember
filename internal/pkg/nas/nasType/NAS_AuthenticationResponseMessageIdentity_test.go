package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayfield-z/ember/internal/pkg/nas/nasMessage"
	"github.com/mayfield-z/ember/internal/pkg/nas/nasType"
)

type nasTypeResponseMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeResponseMessageIdentityTable = []nasTypeResponseMessageIdentityData{
	{nasMessage.AuthenticationResponseEAPMessageType, nasMessage.AuthenticationResponseEAPMessageType},
}

func TestNasTypeNewAuthenticationResponseMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationResponseMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetAuthenticationResponseMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationResponseMessageIdentity()
	for _, table := range nasTypeResponseMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
