package ai

import (
	"testing"

	"github.com/tsybulin/gosha/aut"
)

func Test_greaterThanCondition_SatisfiedGreaterThan(t *testing.T) {
	type fields struct {
		Condition aut.Condition
		entityID  string
		value     float64
	}
	type args struct {
		id    string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "0 > 1", fields: fields{entityID: "comp1", value: 1}, args: args{id: "comp1", value: 0}, want: false},
		{name: "5 > 3", fields: fields{entityID: "comp1", value: 3}, args: args{id: "comp1", value: 5}, want: true},
		{name: "0.1 > 1.1", fields: fields{entityID: "comp1", value: 1.1}, args: args{id: "comp1", value: 0.1}, want: false},
		{name: "5.5 > 3.3", fields: fields{entityID: "comp1", value: 3.3}, args: args{id: "comp1", value: 5.5}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gtc := &greaterThanCondition{
				Condition: tt.fields.Condition,
				entityID:  tt.fields.entityID,
				value:     tt.fields.value,
			}
			if got := gtc.SatisfiedCompare(tt.args.id, tt.args.value); got != tt.want {
				t.Errorf("greaterThanCondition.SatisfiedGreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
