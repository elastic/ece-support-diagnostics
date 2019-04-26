package systemInfo

// func Test_systemCmd_run(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		cmd     systemCmd
// 		want    []byte
// 		wantErr bool
// 	}{
// 		{
// 			name:    "echo Test",
// 			cmd:     systemCmd{RawCmd: "printf hello"},
// 			want:    []byte(`hello`),
// 			wantErr: false,
// 		},
// 		{
// 			name:    "No command",
// 			cmd:     systemCmd{RawCmd: "UNKNOWN command"},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := tt.cmd
// 			got, err := c.run()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("systemCmd.run() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("systemCmd.run() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
