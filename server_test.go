package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestValidateAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name     string
		args     args
		wantErrs []error
	}{
		{"valid", args{"192.168.1.2"}, nil},
		{"valid.port", args{"192.168.1.2:7777"}, nil},
		{"valid.scheme", args{"samp://192.168.1.2"}, nil},
		{"invalid.empty", args{""}, []error{fmt.Errorf("address is empty")}},
		{"invalid.port", args{"192.168.1.2:port"}, []error{fmt.Errorf("invalid port 'port' specified")}},
		{"invalid.scheme", args{"http://192.168.1.2"}, []error{fmt.Errorf("address contains invalid scheme 'http', must be either empty or 'samp://'")}},
		{"invalid.user", args{"user:pass@192.168.1.2"}, []error{fmt.Errorf("address contains a user:password component")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := ValidateAddress(tt.args.address); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("ValidateAddress() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}
