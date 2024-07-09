package validator

import "testing"

func Test_sortOptionsValidator_Validate(t *testing.T) {
	type fields struct {
		baseValidator baseValidator
		criterion     string
		order         string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Valid sort options - date asc",
			fields: fields{
				criterion: "date",
				order:     "asc",
			},
			want: true,
		},
		{
			name: "Valid sort options - source desc",
			fields: fields{
				criterion: "source",
				order:     "desc",
			},
			want: true,
		},
		{
			name: "Invalid sort criterion",
			fields: fields{
				criterion: "invalid",
				order:     "asc",
			},
			want: false,
		},
		{
			name: "Invalid sort order",
			fields: fields{
				criterion: "date",
				order:     "invalid",
			},
			want: false,
		},
		{
			name: "Valid sort options - empty criterion and order",
			fields: fields{
				criterion: "",
				order:     "",
			},
			want: true,
		},
		{
			name: "Invalid sort criterion and order",
			fields: fields{
				criterion: "invalid",
				order:     "invalid",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := sortOptionsValidator{
				baseValidator: tt.fields.baseValidator,
				criterion:     tt.fields.criterion,
				order:         tt.fields.order,
			}
			if got := v.Validate(); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
