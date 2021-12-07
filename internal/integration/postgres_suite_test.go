//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/duyquang6/quote-today/internal/database"
	"github.com/duyquang6/quote-today/internal/serverenv"
	"github.com/duyquang6/quote-today/internal/setup"

	"github.com/stretchr/testify/suite"
)

type PostgresRepositoryTestSuite struct {
	env    *serverenv.ServerEnv
	config database.Config
	suite.Suite
}

func (p *PostgresRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	var config database.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		panic(err)
	}
	p.env = env
}

func TestMySqlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &PostgresRepositoryTestSuite{})
}

func (p *PostgresRepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	if err := p.env.Database().Migrate(ctx); err != nil {
		panic(err)
	}
}

func (p *PostgresRepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	if err := p.env.Database().MigrateDown(ctx); err != nil {
		panic(err)
	}
}
