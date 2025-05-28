package humantime

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSpan(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{
			name: "1時間09分のとき、1h10mと返る",
			d:    time.Hour + 9*time.Minute,
			want: "1h09m",
		},
		{
			name: "1時間ちょうどのとき、1hと返る(1h0mにならない)",
			d:    time.Hour,
			want: "1h",
		},
		{
			name: "1時間未満(59分)のとき、59mと返る(0h59mにならない)",
			d:    59 * time.Minute,
			want: "59m",
		},
		{
			name: "30秒のとき、四捨五入されて01mと返る",
			d:    59 * time.Second,
			want: "01m",
		},
		{
			name: "29秒のとき、四捨五入されて0mと返る",
			d:    time.Second * 29,
			want: "0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Span(tt.d)
			require.Equal(t, tt.want, got)
		})
	}
}
