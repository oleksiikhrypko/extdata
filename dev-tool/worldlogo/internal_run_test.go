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
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for _, l := range alphabet {
		hasMore := true
		for i := 1; hasMore; i++ {
			requrl := fmt.Sprintf("https://worldvectorlogo.com/alphabetical/%s/%d", l, i)
			fmt.Printf("requrl: %s\n", requrl)
			res, err := TakeWorldLogo(ctx, requrl)
			require.NoError(t, err)
			hasMore = len(res) != 0
			err = WriteToSCV("./data/worldlogo.csv", res)
			require.NoError(t, err)
			<-time.After(10 * time.Second)
		}
	}
}

func Test_SendData(t *testing.T) {
	if os.Getenv("DEV_LOCAL_RUN_WORLDLOGO") != "YES" {
		t.Skip()
	}

	err := WriteDataFromCSVToAPI("./data/worldlogo.csv", "http://localhost:8080")
	require.NoError(t, err)
}
