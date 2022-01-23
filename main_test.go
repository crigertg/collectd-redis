package main

import (
	"os"
	"testing"
)

func Test_getCollectdIntervalDefault(t *testing.T) {
	result := getCollectdInterval()
	if result != 10 {
		t.Errorf("Expected 10 put got %f.", result)
	}
}

func Test_getCollectdIntervalAlteredEnv(t *testing.T) {
	os.Setenv("COLLECTD_INTERVAL", "7")
	defer os.Unsetenv("COLLECTD_INTERVAL")
	result := getCollectdInterval()
	if result != 7 {
		t.Errorf("Expected 10 put got %f.", result)
	}
}
