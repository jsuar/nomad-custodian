package nomadhelper

import (
	"fmt"
	"testing"

	nomad "github.com/hashicorp/nomad/api"
)

func TestDisplayJobDiff(t *testing.T) {
	type args struct {
		diff nomad.JobDiff
	}
	fieldDiff := make([]*nomad.FieldDiff, 0)
	fmt.Println("fieldDiff", fieldDiff)
	fieldDiff = append(fieldDiff, &nomad.FieldDiff{"fieldType", "fieldName", "old", "new", make([]string, 0)})
	fmt.Println("fieldDiff", fieldDiff)
	tests := []struct {
		name string
		args args
	}{
		{
			"Nil Test",
			args{
				nomad.JobDiff{"one", "two", nil, nil, nil},
			},
		},
		{
			"Field Test",
			args{
				nomad.JobDiff{
					"one", "two",
					fieldDiff,
					nil, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DisplayJobDiff(tt.args.diff)
		})
	}
}

// func TestNomadHelper_Init(t *testing.T) {
// 	type fields struct {
// 		Client *nomad.Client
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			"Nil",
// 			fields{nil},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			n := &NomadHelper{
// 				Client: tt.fields.Client,
// 			}
// 			n.Init()
// 		})
// 	}
// }

func TestNomadHelper_InitConfig(t *testing.T) {
	type fields struct {
		Client *nomad.Client
		Config *nomad.Config
	}
	type args struct {
		config *nomad.Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"Nil Config",
			fields{nil, nil},
			args{nil},
		},
		{
			"Default Config",
			fields{nil, nil},
			args{nomad.DefaultConfig()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// n := &NomadHelper{
			// 	Client: tt.fields.Client,
			// 	Config: tt.fields.Config,
			// }
			n := new(NomadHelper)
			n.InitConfig(tt.args.config)
		})
	}
}
