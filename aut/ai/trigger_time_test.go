package ai

import (
	"testing"
	"time"
)

func Test_timeTrigger_FireTime(t *testing.T) {
	tests := []struct {
		name string
		at   time.Time
		args time.Time
		want bool
	}{
		{name: "now and before", at: time.Now(), args: time.Now().Add(-time.Minute), want: false},
		{name: "now and now", at: time.Now(), args: time.Now(), want: true},
		{name: "now and after", at: time.Now(), args: time.Now().Add(time.Minute), want: false},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tr := newTimeTriggerAt(tt.at)

			if got := tr.FireTime(tt.args); got != tt.want {
				t.Errorf("timeTrigger.FireTime() = %v, want %v", got, tt.want)
			}

			if tt.want {
				tt.want = false
			}

			if got := tr.FireTime(tt.args); got != tt.want {
				t.Errorf("timeTrigger.FireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTimeTrigger(t *testing.T) {
	type want struct {
		h int
		m int
	}
	tests := []struct {
		name string
		cfg  map[string]string
		want want
	}{
		{name: "", cfg: map[string]string{"at": "19:10"}, want: want{h: 19, m: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTimeTriggers(tt.cfg)[0]; got.GetAt().Hour() != tt.want.h || got.GetAt().Minute() != tt.want.m {
				t.Errorf("newTimeTrigger() = %v %v, want %v %v", got.GetAt().Hour(), got.GetAt().Minute(), tt.want.h, tt.want.m)
			}
		})
	}
}
