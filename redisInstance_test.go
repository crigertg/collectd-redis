package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_validateConnectionString(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "localhost",
			args: args{
				"local:localhost:6379",
			},
			wantErr: false,
		},
		{
			name: "remotehost",
			args: args{
				"remotehost:redis-remote-intern:6379",
			},
			wantErr: false,
		},
		{
			name: "iphost",
			args: args{
				"ip:192.168.1.233:6379",
			},
			wantErr: false,
		},
		{
			name: "hostWithPassword",
			args: args{
				"pw:localhost:6379:secret",
			},
			wantErr: false,
		},
		{
			name: "missing_field",
			args: args{
				"localhost:6379",
			},
			wantErr: true,
		},
		{
			name: "invalidPort",
			args: args{
				"pw:localhost:234234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConnectionString(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("validateConnectionString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseArgToInstance(t *testing.T) {
	type args struct {
		connStr string
	}
	tests := []struct {
		name string
		args args
		want redisInstance
	}{
		{
			name: "localhost",
			args: args{
				"local:localhost:6379",
			},
			want: redisInstance{
				name: "local",
				host: "localhost",
				port: 6379,
				pw:   "",
			},
		},
		{
			name: "hostWithPassword",
			args: args{
				"pw:localhost:6379:secret",
			},
			want: redisInstance{
				name: "pw",
				host: "localhost",
				port: 6379,
				pw:   "secret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseArgToInstance(tt.args.connStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConnStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestRedisInfoOutput() string {
	// load content from static file representing redis INFO output
	content, err := ioutil.ReadFile("test/redis_output")
	errCheckFatal(err)
	return string(content)
}
