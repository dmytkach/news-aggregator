package validator

import (
	"errors"
	"testing"
)

func Test_sortOptionsValidator_Validate(t *testing.T) {
	type fields struct {
		baseValidator baseValidator
		criterion     string
		order         string
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "Valid sort options - date asc",
			fields: fields{
				criterion: "date",
				order:     "asc",
			},
			want: nil,
		},
		{
			name: "Valid sort options - source desc",
			fields: fields{
				criterion: "source",
				order:     "desc",
			},
			want: nil,
		},
		{
			name: "Invalid sort criterion",
			fields: fields{
				criterion: "invalid",
				order:     "asc",
			},
			want: errors.New("invalid sort criterion. Please use `date` or `source`"),
		},
		{
			name: "Invalid sort order",
			fields: fields{
				criterion: "date",
				order:     "invalid",
			},
			want: errors.New("invalid order criterion. Please use `asc` or `desc`"),
		},
		{
			name: "Valid sort options - empty criterion and order",
			fields: fields{
				criterion: "",
				order:     "",
			},
			want: nil,
		},
		{
			name: "Invalid sort criterion and order",
			fields: fields{
				criterion: "invalid",
				order:     "invalid",
			},
			want: errors.New("invalid sort criterion. Please use `date` or `source`"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := sortOptionsValidator{
				baseValidator: tt.fields.baseValidator,
				criterion:     tt.fields.criterion,
				order:         tt.fields.order,
			}
			result := v.Validate()
			if tt.want != nil {
				if result == nil || result.Error() != tt.want.Error() {
					t.Errorf("Expected %v, but got %v", tt.want, result)
				}
			}
		})
	}
}
