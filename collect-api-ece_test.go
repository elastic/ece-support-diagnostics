package ecediag

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_readJSON(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want interface{}
	}{
		{
			name: "json_object",
			in:   []byte(`{"json":"object"}`),
			want: map[string]interface{}{"json": "object"},
		},
		{
			name: "json_array",
			in:   []byte(`[{"json":"array"}]`),
			want: []interface{}{
				map[string]interface{}{"json": "array"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readJSON(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runTemplate(t *testing.T) {

	type args struct {
		item string
		Obj  interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				item: "{{ .t }}_test",
				Obj:  readJSON([]byte(`{"t":"test"}`)),
			},
			want: "test_test",
		},
		{
			name: "nested",
			args: args{
				item: "{{ .t.t1 }}_test",
				Obj:  readJSON([]byte(`{"t":{"t1":"nested"}}`)),
			},
			want: "nested_test",
		},
		{
			name: "array",
			args: args{
				item: "{{ (index .levelone.leveltwo 1).item }}_test",
				Obj:  readJSON([]byte(`{"levelone":{"leveltwo":[{"item":"first"},{"item":"second"}]}}`)),
			},
			want: "second_test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runTemplate(tt.args.item, tt.args.Obj); got != tt.want {
				t.Errorf("runTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRest_templater(t *testing.T) {
	tests := []struct {
		name string
		item Rest
		Obj  interface{}
	}{
		{
			name: "test1",
			item: Rest{
				Filename: "{{ .test1 }}/test1",
				URI:      "{{ .test2 }}/test2",
			},
			Obj: readJSON([]byte(`{"test1":"value1","test2":"value2"}`)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.templater(tt.Obj)
			if tt.item.Filename != "value1/test1" {
				t.Errorf("expected value1/test1, got %s", tt.item.Filename)
			}
			if tt.item.URI != "value2/test2" {
				t.Errorf("expected value2/test2, got %s", tt.item.URI)
			}
		})
	}
}

func Test_getCredentials(t *testing.T) {

	content := []byte("myusername\nblah\n")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	user, pass := getCredentials()
	fmt.Printf("\nuser: %s, pass: %s\n", user, pass)
	if user != "myusername" {
		t.Errorf("userInput failed: %v", err)
	}
	if pass != "THIS DOES NOT WORK" {
		t.Errorf("password input failed: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}
