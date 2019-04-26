package tar

// func TestTarball_Create(t *testing.T) {
// 	type args struct {
// 		filePath string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "create tar.gz",
// 			args:    args{filePath: "tmp/unittest.create.tar.gz"},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tw := new(Tarball)
// 			tw.Create(tt.args.filePath)

// 			defer os.Remove(tt.args.filePath) // clean up
// 			tw.t.Close()
// 			tw.g.Close()

// 			if _, err := os.Stat(tt.args.filePath); os.IsNotExist(err) {
// 				t.Errorf("Tarball.Create() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
