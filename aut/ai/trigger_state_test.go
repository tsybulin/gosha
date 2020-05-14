package ai

import (
	"testing"
)

/*
func Test_stateTrigger_FireState(t *testing.T) {
	type args struct {
		id string
		v  interface{}
	}
	tests := []struct {
		name    string
		trigger aut.StateTrigger
		args    args
		want    bool
	}{
		{name: "nill to false ", args: args{id: "id1", v: false}, want: false, trigger: newStateTrigger("id1", false, "on")},
		{name: "false to false", args: args{id: "id1", v: false}, want: false},
		{name: "false to true ", args: args{id: "id1", v: true}, want: true},
		{name: "true to true  ", args: args{id: "id1", v: true}, want: false},
		{name: "true to false ", args: args{id: "id1", v: false}, want: false},
		{name: "any true ", args: args{id: "id1", v: false}, want: true, trigger: newStateTrigger("id1", nil, nil)},
		{name: "any string ", args: args{id: "id1", v: "active"}, want: true},
		{name: "to true ", args: args{id: "id1", v: true}, want: true, trigger: newStateTrigger("id1", nil, true)},
		{name: "to false ", args: args{id: "id1", v: false}, want: false},
	}

	var st aut.StateTrigger

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.trigger != nil {
				st = tt.trigger
			}

			if got := st.FireState(tt.args.id, tt.args.v); got != tt.want {
				t.Errorf("stateTrigger.FireState() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func Test_newStateTriggers(t *testing.T) {
	type args struct {
		cfg map[string]string
	}
	tests := []struct {
		name string
		args map[string]string
		want int
	}{
		{
			name: "1 trigger",
			args: map[string]string{
				"components": "cmp1",
				"from":       "off",
				"to":         "on",
			},
			want: 1,
		},
		{
			name: "2 triggers",
			args: map[string]string{
				"components": "cmp1     cmp2",
				"from":       "off",
				"to":         "on",
			},
			want: 2,
		},
		{
			name: "no triggers",
			args: map[string]string{
				"components": " ",
				"from":       "off",
				"to":         "on",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newStateTriggers(tt.args); len(got) != tt.want {
				t.Errorf("newStateTriggers() = %v, want %v", got, tt.want)
			}
		})
	}
}
