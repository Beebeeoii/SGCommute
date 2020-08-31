package main

import (
	"reflect"
	"testing"
)

func Test_retrieveBusDetailsFromLTA(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name     string
		args     args
		wantData AllBusesDetails
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := retrieveBusDetailsFromLTA(tt.args.apiKey); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("retrieveBusDetailsFromLTA() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
