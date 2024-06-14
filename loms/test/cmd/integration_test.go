//go:build integration

package cmd

import (
	"testing"

	"route256/loms/test/suits"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(suits.InmemSuite))
}
