package api

import (
	"os"
	"testing"
	"time"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/danielmoisa/neobank/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
