package ai

import (
	"testing"
)

func Test_stateCondition_SatisfiedState(t *testing.T) {
	type fields struct {
		entityID string
		state    string
	}
	type args struct {
		id string
		v  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "'inactive' on mine", fields: fields{entityID: "comp1", state: "active"}, args: args{id: "comp1", v: "inactive"}, want: false},
		{name: "'active' on mine", fields: fields{entityID: "comp1", state: "active"}, args: args{id: "comp1", v: "active"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newStateCondition(tt.fields.entityID, tt.fields.state)
			if got := c.SatisfiedState(tt.args.id, tt.args.v); got != tt.want {
				t.Errorf("stateCondition.SatisfiedState() = %v, want %v", got, tt.want)
			}
		})
	}
}
