package docker

// func Test_safeFilename(t *testing.T) {
// 	type args struct {
// 		names []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "Multi strings to combine",
// 			args: args{names: []string{"hello", "small", "world?"}},
// 			want: "hello__small__world_",
// 		},
// 		{
// 			name: "Back Slash",
// 			args: args{[]string{"hello\\"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Forward Slash",
// 			args: args{[]string{"hello/"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Colon",
// 			args: args{[]string{"hello:"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Asterisk",
// 			args: args{[]string{"hello*"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Question Mark",
// 			args: args{[]string{"hello?"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Double Quote",
// 			args: args{[]string{"hello\""}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Less Than",
// 			args: args{[]string{"hello<"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Greater Than",
// 			args: args{[]string{"hello>"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Verticle Bar",
// 			args: args{[]string{"hello|"}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Period",
// 			args: args{[]string{"hello."}},
// 			want: "hello_",
// 		},
// 		{
// 			name: "Docker.elastic.co",
// 			args: args{[]string{"hello.docker.elastic.co"}},
// 			want: "hello_",
// 		},

// 		// "docker.elastic.co", "",
// 		// "\\", "_",
// 		// "/", "_",
// 		// ":", "_",
// 		// "*", "_",
// 		// "?", "_",
// 		// "\"", "_",
// 		// "<", "_",
// 		// ">", "_",
// 		// "|", "_",
// 		// ".", "_",
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := safeFilename(tt.args.names...); got != tt.want {
// 				t.Errorf("safeFilename() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_fp(t *testing.T) {
// 	type args struct {
// 		filename []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := fp(tt.args.filename...); got != tt.want {
// 				t.Errorf("fp() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
