package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_fetchMetricValue(t *testing.T) {
	redisInfo := getTestRedisInfoOutput()
	type args struct {
		reIn    string
		metrNam string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "fetch_total_net_input_bytes",
			args: args{
				reIn:    redisInfo,
				metrNam: "total_net_input_bytes",
			},
			want: 72,
		},
		{
			name: "fetch_allocator_frag_bytes",
			args: args{
				reIn:    redisInfo,
				metrNam: "allocator_frag_bytes",
			},
			want: 173736,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchMetricValue(tt.args.reIn, tt.args.metrNam)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchMetricValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchMetricValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateRecordsMetrics(t *testing.T) {
	result := generateRecordsMetrics(getTestRedisInfoOutput())
	expected := []redisMetric{
		{
			name:  "records-0",
			value: 1,
		},
		{
			name:  "records-2",
			value: 2,
		},
	}
	if diff := cmp.Diff(expected, result, cmp.AllowUnexported(redisMetric{})); diff != "" {
		t.Error(diff)
	}
}

func Test_generateUniqueMetrics(t *testing.T) {
	result := generateUniqueMetrics(getTestRedisInfoOutput())
	expected := []redisMetric{
		{
			name:  "blocked_clients",
			value: 0,
		},
		{
			name:  "current_connections-clients",
			value: 1,
		},
		{
			name:  "current_connections-slaves",
			value: 0,
		},
		{
			name:  "evicted_keys",
			value: 0,
		},
		{
			name:  "expired_keys",
			value: 0,
		},
		{
			name:  "cache_result-hits",
			value: 0,
		},
		{
			name:  "cache_result-misses",
			value: 0,
		},
		{
			name:  "pubsub-channels",
			value: 0,
		},
		{
			name:  "pubsub-patterns",
			value: 0,
		},
		{
			name:  "volatile_changes",
			value: 4,
		},
		{
			name:  "total_operations",
			value: 5,
		},
		{
			name:  "total_connections",
			value: 1,
		},
		{
			name:  "total_bytes-input",
			value: 72,
		},
		{
			name:  "total_bytes-output",
			value: 24,
		},
		{
			name:  "uptime",
			value: 41,
		},
		{
			name:  "ps_cputime-daemon/syst",
			value: 0.053338,
		},
		{
			name:  "ps_cputime-children/syst",
			value: 0.003915,
		},
		{
			name:  "ps_cputime-daemon/user",
			value: 0.070197,
		},
		{
			name:  "ps_cputime-children/user",
			value: 0,
		},
		{
			name:  "memory",
			value: 873896,
		},
		{
			name:  "memory_lua",
			value: 37888,
		},
	}
	if diff := cmp.Diff(expected, result, cmp.AllowUnexported(redisMetric{})); diff != "" {
		t.Error(diff)
	}
}

func Test_parsePutvalString(t *testing.T) {
	type args struct {
		instance string
		metric   redisMetric
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Put CPU metric",
			args: args{
				instance: "caching",
				metric: redisMetric{
					name:  "ps_cputime-daemon/syst",
					value: 3.563614,
				},
			},
			want: "PUTVAL \"localhost/redis-caching/ps_cputime-daemon/syst\" interval=10.000000 N:4\n",
		},
		{
			name: "Put Connection Metric",
			args: args{
				instance: "caching",
				metric: redisMetric{
					name:  "total_connections",
					value: 236,
				},
			},
			want: "PUTVAL \"localhost/redis-caching/total_connections\" interval=10.000000 N:236\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePutvalString(tt.args.instance, tt.args.metric); got != tt.want {
				t.Errorf("parsePutvalString() = %v, want %v", got, tt.want)
			}
		})
	}
}
