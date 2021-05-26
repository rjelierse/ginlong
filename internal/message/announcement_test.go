package message

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnouncement(t *testing.T) {
	messages := []string{
		"02f9e174005b01000000000000053c780164034d575f30385f303530315f312e35380000000000000000000000000000000000000000000000000098d863e473b23139322e3136382e312e323532000000f64f010105",
		"027de374005901000000000000053c780164034d575f30385f303530315f312e35380000000000000000000000000000000000000000000000000098d863e473b23139322e3136382e312e323532000000f84f010105",
	}
	for idx, payload := range messages {
		t.Run(fmt.Sprintf("Message %d", idx+1), func(t *testing.T) {
			data, err := hex.DecodeString(payload)
			require.NoError(t, err)

			var message Announcement
			require.NoError(t, message.UnmarshalBinary(data))

			assert.Equal(t, "MW_08_0501_1.58", message.FirmwareVersion)
			assert.Equal(t, "192.168.1.252", message.IPAddress.String())
			assert.Equal(t, "98:d8:63:e4:73:b2", message.HardwareAddress.String())
		})
	}
}
