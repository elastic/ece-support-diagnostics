package helpers

import (
	"testing"
)

func TestByteCountDecimal(t *testing.T) {
	type args struct {
		b int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "100B",
			args: args{b: 100},
			want: "100 B",
		},
		{
			name: "100kB",
			args: args{b: 100000},
			want: "100.0 kB",
		},
		{
			name: "100MB",
			args: args{b: 100000000},
			want: "100.0 MB",
		},
		{
			name: "100GB",
			args: args{b: 1e+11},
			want: "100.0 GB",
		},
		{
			name: "100TB",
			args: args{b: 1e+14},
			want: "100.0 TB",
		},
		{
			name: "100PB",
			args: args{b: 1e+17},
			want: "100.0 PB",
		},
		{
			name: "1EB",
			args: args{b: 1e+18},
			want: "1.0 EB",
			// 100EB (1e+20) is too large for int64
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteCountDecimal(tt.args.b); got != tt.want {
				t.Errorf("ByteCountDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteCountBinary(t *testing.T) {
	type args struct {
		b int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "100B",
			args: args{b: 100},
			want: "100 B",
		},
		{
			name: "100 KiB",
			args: args{b: 102400},
			want: "100.0 KiB",
		},
		{
			name: "100 MiB",
			args: args{b: 104857600},
			want: "100.0 MiB",
		},
		{
			name: "100 GiB",
			args: args{b: 107374182400},
			want: "100.0 GiB",
		},
		{
			name: "100 TiB",
			args: args{b: 1.09951162778e+14},
			want: "100.0 TiB",
		},
		{
			name: "100 PiB",
			args: args{b: 1.12589990684e+17},
			want: "100.0 PiB",
		},
		{
			name: "1 EiB",
			args: args{b: 1.15292150461e+18},
			want: "1.0 EiB",
			// 100 EiB (1.44115188076E+19) is too large for int64
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteCountBinary(tt.args.b); got != tt.want {
				t.Errorf("ByteCountBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}
