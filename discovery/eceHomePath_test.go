package discovery

import "testing"

func TestCheckStoragePathError(t *testing.T) {
	testFolder := "test_data"
	_, err := CheckStoragePath(testFolder)
	if err.Error() != "Could not find a valid ECE install location" {
		t.Errorf("%s", err)
	}
}

func TestCheckStoragePath(t *testing.T) {
	testFolder := "test_data/elastic"
	expectedRunner := "172.16.0.71"
	RunnerName, _ := CheckStoragePath(testFolder)
	if RunnerName != "172.16.0.71" {
		t.Errorf("RunnerName was incorrect, got: %s, expected: %s.", RunnerName, expectedRunner)
	}
}
