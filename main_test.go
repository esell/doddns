package main

import "testing"

func TestSubExists(t *testing.T) {
	dummyRecords := make([]DomainRecord, 0)

	doesExist, recID := subExists(dummyRecords)

	if doesExist {
		t.Errorf("subExists() returned true, should have been false")
	}

	if recID != 0 {
		t.Errorf("subExists() returned %v for recID, should be 0", recID)
	}
}
