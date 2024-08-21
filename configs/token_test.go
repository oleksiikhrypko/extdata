package configs

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

type claims map[string]interface{}

func (c claims) Valid() error {
	return nil
}

func generateToken(t *testing.T, path string) string {
	// Create and sign the token using the HMAC key.
	unsignedToken := jwt.New(jwt.SigningMethodHS256)
	unsignedToken.Header["kid"] = "IntegrationTests"

	b, err := os.ReadFile(path)
	require.NoError(t, err, "can't open token data")

	var cl claims
	err = json.Unmarshal(b, &cl)
	require.NoError(t, err, "can't unmarshal token data")

	unsignedToken.Claims = cl
	jwtB64, err := unsignedToken.SignedString([]byte("TestSecret"))
	require.NoError(t, err, "failed to sign")

	return jwtB64
}

func Test_GenLocalToken(t *testing.T) {
	t.Skip()
	t.Log(generateToken(t, "./token.user.local.json"))
	// eyJhbGciOiJIUzI1NiIsImtpZCI6IkludGVncmF0aW9uVGVzdHMiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiaHR0cHM6Ly9pZGVhLWRvbWFpbi5kZXYuc2x5bmdzaG90LmFpLyJdLCJhenAiOiJGTDJPNkpNS3ZDbUJKdEhTMnVPdzdZSEw0c3pnNFFoeSIsImV4cCI6MjY4Nzk4MjgyNywiaWF0IjoxNjg3ODc1NjI3LCJpc3MiOiJodHRwczovL2F1dGhvcml6YXRpb24tcHJveHkuZGV2LnNseW5nc2hvdC5haS8iLCJwZXJtaXNzaW9ucyI6WyJzbHluZ3Nob3QuaWRlYS1kb21haW46cXVlcnk6dXNlciJdLCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic3ViIjoiYXV0aDB8NjQ4OTk5NzJhNWM2ZjExMTU2ZGMzYTZkIiwieC11c2VyLWlkIjoiMDJHR0dHNTQxWFRWMktZTkEyMTExU1VTUjEifQ.Zhlwf1fxGPqQRQp5qNlqrclR0qgaf-EquBmRkJHtOxc
	t.Log(generateToken(t, "./token.admin.local.json"))
	// eyJhbGciOiJIUzI1NiIsImtpZCI6IkludGVncmF0aW9uVGVzdHMiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiaHR0cHM6Ly9pZGVhLWRvbWFpbi5kZXYuc2x5bmdzaG90LmFpLyJdLCJhenAiOiJGTDJPNkpNS3ZDbUJKdEhTMnVPdzdZSEw0c3pnNFFoeSIsImV4cCI6MjY4Nzk4MjgyNywiaWF0IjoxNjg3ODc1NjI3LCJpc3MiOiJodHRwczovL2F1dGhvcml6YXRpb24tcHJveHkuZGV2LnNseW5nc2hvdC5haS8iLCJwZXJtaXNzaW9ucyI6WyJzbHluZ3Nob3QuaWRlYS1kb21haW46cXVlcnk6YWRtaW4iXSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInN1YiI6ImF1dGgwfDY0ODk5OTcyYTVjNmYxMTE1NmRjM2E2ZCIsIngtdXNlci1pZCI6IjAyR0dHRzU0MVhUVjJLWU5BMjExMVNBRE0xIn0.pTg2gmfx1wrOunVZ1cp3tHZkoTl_TxFYMgWuxkwZ7x0
}
