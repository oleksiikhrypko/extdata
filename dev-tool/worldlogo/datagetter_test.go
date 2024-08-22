package worldlogo

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_getData(t *testing.T) {
	if os.Getenv("DEV_LOCAL_RUN_WORLDLOGO") != "YES" {
		t.Skip()
	}

	ctx := context.Background()

	hasMore := true
	for i := 1; hasMore; i++ {
		res, err := TakeWorldLogo(ctx, fmt.Sprintf("https://worldvectorlogo.com/alphabetical/a/%d", i))
		require.NoError(t, err)
		hasMore = len(res) != 0
		err = WriteToSCV("./data/worldlogo.csv", res)
		require.NoError(t, err)
		<-time.After(2 * time.Second)
	}
}
