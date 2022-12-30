package paas

import (
	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInWindow(t *testing.T) {
	sundayOrMondayAfternoon := "* 14,15,16,17,18,19 * * SUN,MON"
	sundayAtFivePM := time.UnixMilli(1672552800000)
	tuesdayAtFivePM := time.UnixMilli(1672725600000)

	c, err := cron.ParseStandard(sundayOrMondayAfternoon)
	require.NoError(t, err)

	assert.True(t, inWindow(c, sundayAtFivePM))
	assert.False(t, inWindow(c, tuesdayAtFivePM))
}

func TestInWindow2(t *testing.T) {
	anyTime := "* * * * *"
	sundayAtFivePM := time.UnixMilli(1672552800000)
	tuesdayAtFivePM := time.UnixMilli(1672725600000)

	c, err := cron.ParseStandard(anyTime)
	require.NoError(t, err)

	assert.True(t, inWindow(c, sundayAtFivePM))
	assert.True(t, inWindow(c, tuesdayAtFivePM))
}
