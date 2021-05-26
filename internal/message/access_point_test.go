package message

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessPointInfo(t *testing.T) {
	data, err := hex.DecodeString("81ffe075000500000000000000000048617a6573000000000000000000000000000000000000000000000000006401")
	require.NoError(t, err)

	fmt.Print(hex.Dump(data))

	var message AccessPointInfo
	require.NoError(t, message.UnmarshalBinary(data))

	assert.Equal(t, "Hazes", message.SSID)
}
