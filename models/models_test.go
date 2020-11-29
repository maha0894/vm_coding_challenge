package models

import (
	"reflect"
	"testing"
)

func TestGetStats(t *testing.T) {
	t.Run("gets stats", func(t *testing.T) {
		// TODO setup data for stats
		var expStats Stats

		stats, _ := GetStats("customer", 0, 1)
		if !reflect.DeepEqual(expStats, stats) {
			t.Errorf("got %q, want %q", stats, expStats)
		}
	})
}
