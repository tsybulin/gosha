package ai

import (
	"testing"
	"time"

	"github.com/tsybulin/gosha/aut"
)

func Test_lessThanCondition_SatisfiedCompare(t *testing.T) {
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
		{name: "0 < 1", fields: fields{entityID: "comp1", value: 1}, args: args{id: "comp1", value: 0}, want: true},
		// {name: "0 < 1", fields: fields{entityID: "comp1", value: 1}, args: args{id: "comp1", value: "0"}, want: true},
		{name: "5 < 3", fields: fields{entityID: "comp1", value: 3}, args: args{id: "comp1", value: 5}, want: false},
		{name: "0.1 < 1.1", fields: fields{entityID: "comp1", value: 1.1}, args: args{id: "comp1", value: 0.1}, want: true},
		{name: "5.5 < 3.3", fields: fields{entityID: "comp1", value: 3.3}, args: args{id: "comp1", value: 5.5}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ltc := &lessThanCondition{
				Condition:        tt.fields.Condition,
				entityID:         tt.fields.entityID,
				value:            tt.fields.value,
				firstSatisfiedAt: time.Time{},
				forSeconds:       0,
			}
			if got := ltc.SatisfiedCompare(tt.args.id, tt.args.value); got != tt.want {
				t.Errorf("lessThanCondition.SatisfiedLessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
