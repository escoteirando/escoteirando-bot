package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeParse(t *testing.T) {
	type args struct {
		timeString string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Valid date",
			args: args{timeString: "Sat Feb 05 1977 00:00:00 GMT+0000 (UTC)"},
			want: time.Date(1977, 02, 05, 0, 0, 0, 0, time.UTC),
		}, {
			name: "Valid date",
			args: args{timeString: "2018-05-26T22:13:14.000Z"},
			want: time.Date(2018, 05, 26, 22, 13, 14, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeParse(tt.args.timeString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundStart(t *testing.T) {
	type args struct {
		from     time.Time
		interval time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "5 minutes",
			args: args{from: time.Date(2000, 1, 1, 0, 3, 2, 0, time.UTC),
				interval: time.Duration(5) * time.Minute},
			want: time.Date(2000, 1, 1, 0, 5, 0, 0, time.UTC),
		}, {
			name: "1 hour",
			args: args{from: time.Date(2000, 1, 1, 0, 3, 2, 0, time.UTC),
				interval: time.Duration(1) * time.Hour},
			want: time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
		},
		{
			name: "1 day",
			args: args{from: time.Date(2000, 1, 1, 0, 3, 2, 0, time.UTC),
				interval: time.Duration(24) * time.Hour},
			want: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundStart(tt.args.from, tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundStart() = %v, want %v", got, tt.want)
			}
		})
	}
}
