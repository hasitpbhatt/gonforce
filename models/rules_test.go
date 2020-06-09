package models

import "testing"

func TestPackageRule_Validate(t *testing.T) {
	type fields struct {
		Type    RuleType
		Imports []string
		Except  []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"allowlist", fields{Type: "allowlist"}, false},
		{"blocklist", fields{Type: "blocklist"}, false},
		{"random_type", fields{Type: "random_type"}, true},
		{
			"conflicting import and exceptions",
			fields{
				Type:    "allowlist",
				Imports: []string{"a", "b"},
				Except:  []string{"b"},
			},
			true,
		},
		{
			"valid import and exceptions",
			fields{
				Type:    "allowlist",
				Imports: []string{"a", "b"},
				Except:  []string{"a/b", "b/a", "a/a"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PackageRule{
				Type:    tt.fields.Type,
				Imports: tt.fields.Imports,
				Except:  tt.fields.Except,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PackageRule.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
