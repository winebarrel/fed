package subcmd_test

import (
	"testing"

	"github.com/kanmu/kasa"
	"github.com/kanmu/kasa/esa/model"
	"github.com/kanmu/kasa/subcmd"
	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	assert := assert.New(t)

	stats := &subcmd.StatsCmd{}
	driver := NewMockDriver(t)
	printer := &MockPrinterImpl{}

	driver.MockGetStats = func() (*model.Stats, error) {
		return &model.Stats{}, nil
	}

	err := stats.Run(&kasa.Context{
		Driver: driver,
		Fmt:    printer,
	})

	assert.NoError(err)

	assert.Equal(`{
  "members": 0,
  "posts": 0,
  "posts_wip": 0,
  "posts_shipped": 0,
  "comments": 0,
  "stars": 0,
  "daily_active_users": 0,
  "weekly_active_users": 0,
  "monthly_active_users": 0
}
`, printer.String())
}
