package tickerpool

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

var dummy *TickerPool
var clipboard int32

func TestMain(t *testing.M) {
	dummy, _ = NewTickerPool(time.Second)

	os.Exit(t.Run())
}

func dummyA() {
	fmt.Print("dummyA called\n")
}

func dummyB() {
	fmt.Print("dummyB called\n")
}

func dummyC() {
	fmt.Print("dummyC called\n")
}

func TestNewTickerPool(t *testing.T) {
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{time.Second}, false},
		{"zero", args{0}, true},
		{"neg", args{-10 * time.Second}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTickerPool(tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTickerPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTickerPool_Add(t *testing.T) {
	type args struct {
		name string
		task func()
	}
	tests := []struct {
		name      string
		tp        *TickerPool
		args      args
		wantExist bool
	}{
		{"1", dummy, args{"a", dummyA}, false},
		{"2", dummy, args{"b", dummyB}, false},
		{"3", dummy, args{"c", dummyC}, false},
		{"4", dummy, args{"c", dummyC}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := tt.tp.Add(tt.args.name, tt.args.task)
			if exists != tt.wantExist {
				t.Errorf("'%s' exists when it shouldn't!", tt.args.name)
			}
			time.Sleep(time.Second)
		})
	}

	time.Sleep(time.Second * 5)

	// TODO: Add post test checks
}

func TestTickerPool_Remove(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		tp   *TickerPool
		args args
	}{
		{"1", dummy, args{"a"}},
		{"2", dummy, args{"b"}},
		{"3", dummy, args{"c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tp.Remove(tt.args.name)
		})
	}

	time.Sleep(time.Second * 2)

	// TODO: Add post test checks
}

func Test_calculateWorkerInterval(t *testing.T) {
	type args struct {
		interval time.Duration
		workers  int64
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"1", args{time.Second, 10}, time.Millisecond * 100},
		{"2", args{time.Second, 100}, time.Millisecond * 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateWorkerInterval(tt.args.interval, tt.args.workers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateWorkerInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
