package modifiers

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noxworld-dev/opennox-lib/noxtest"
)

const noxModifier = "modifier.bin"

func TestParse(t *testing.T) {
	f, err := ReadFile(noxtest.DataPath(t, noxModifier))
	require.NoError(t, err)

	data, _ := json.MarshalIndent(f, "", "\t")
	os.WriteFile("modifiers.json", data, 0644)
}
