package options

import (
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/zalando/planb-tokeninfo/processor"
)

func TestGetString(t *testing.T) {
	for _, test := range []struct {
		envSet string
		envGet string
		value  string
		def    string
		want   string
	}{
		{"T1", "T1", "to-be", "or-not", "to-be"},
		{"", "T2", "", "default", "default"},
		{"T3", "SHOULD_NOT_BE_FOUND_IN_ENV", "foo", "bar", "bar"},
	} {
		os.Clearenv()
		if test.envSet != "" {
			os.Setenv(test.envSet, test.value)
		}
		if s := getString(test.envGet, test.def); s != test.want {
			t.Errorf("Failed to retrieve the correct value from the environment. Wanted %q, got %q", test.want, s)
		}
	}
}

func TestGetUrl(t *testing.T) {
	for _, test := range []struct {
		name      string
		value     string
		want      string
		wantError bool
	}{
		{"", "localhost", "", true},
		{"DIFFICULT_TO_GUESS", "localhost", "localhost", false},
		{"DIFFICULT_TO_GUESS", "http://192.168.0.%31/", "", true},
		{"DIFFICULT_TO_GUESS", "", "", true},
		{"DIFFICULT_TO_GUESS", "http://example.com", "http://example.com", false},
	} {
		os.Clearenv()
		if test.name != "" {
			os.Setenv(test.name, test.value)
		}
		u, err := getURL(test.name)
		if test.wantError {
			if err == nil {
				t.Error("Expected an error but call succeeded: ", test)
			}
		} else {
			if u.String() != test.want {
				t.Errorf("Unexpected URL. Wanted %q, got %v", test.want, u)
			}
		}
	}
}

func TestGetInt(t *testing.T) {
	for _, test := range []struct {
		envSet string
		value  string
		envGet string
		def    int
		want   int
	}{
		{"T1", "", "T1", 42, 42},
		{"T1", "invalid-int", "T1", 15, 15},
		{"", "", "DIFFICULT_TO_GUESS", 0, 0},
		{"T1", "7", "T1", 0, 7},
	} {
		os.Clearenv()
		if test.envSet != "" {
			os.Setenv(test.envSet, test.value)
		}
		if s := getInt(test.envGet, test.def); s != test.want {
			t.Errorf("Failed to retrieve the correct value from the environment. Wanted %q, got %q", test.want, s)
		}
	}
}

func TestGetDuration(t *testing.T) {
	for _, test := range []struct {
		envSet string
		value  string
		envGet string
		def    time.Duration
		want   time.Duration
	}{
		{"T1", "", "T1", time.Millisecond, time.Millisecond},
		{"T1", "invalid-duration", "T1", time.Second, time.Second},
		{"", "", "DIFFICULT_TO_GUESS", 0, 0},
		{"T1", "7ns", "T1", 0, 7},
		{"T1", "7ms", "T1", 0, 7 * time.Millisecond},
		{"T1", "30s", "T1", 0, 30 * time.Second},
		{"T1", "30m", "T1", 0, 30 * time.Minute},
		{"T1", "1h", "T1", 0, time.Hour},
		{"T1", "10", "T1", 0, 10 * time.Second},
	} {
		os.Clearenv()
		if test.envSet != "" {
			os.Setenv(test.envSet, test.value)
		}
		if s := getDuration(test.envGet, test.def); s != test.want {
			t.Errorf("Failed to retrieve the correct value from the environment. Wanted %q, got %q", test.want, s)
		}
	}
}

func TestLoading(t *testing.T) {
	exampleCom, _ := url.Parse("http://example.com")
	for _, test := range []struct {
		name     string
		env      map[string]string
		want     *Settings
		wantFail bool
	}{
		{"empty", map[string]string{}, nil, true},
		{
			"UPSTREAM_TOKENINFO_URL empty",
			map[string]string{"UPSTREAM_TOKENINFO_URL": ""},
			&Settings{
				UpstreamTokenInfoURL:              nil,
				OpenIDProviderConfigurationURL:    nil,
				RevocationProviderUrl:             nil,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"UPSTREAM_TOKENINFO_URL set",
			map[string]string{"UPSTREAM_TOKENINFO_URL": "http://example.com"},
			&Settings{
				UpstreamTokenInfoURL:              nil,
				OpenIDProviderConfigurationURL:    nil,
				RevocationProviderUrl:             nil,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"OPENID_PROVIDER_CONFIGURATION_URL set",
			map[string]string{
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
			},
			&Settings{
				UpstreamTokenInfoURL:              nil,
				OpenIDProviderConfigurationURL:    nil,
				RevocationProviderUrl:             nil,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"OPENID_PROVIDER_CONFIGURATION_URL empty",
			map[string]string{
				"OPENID_PROVIDER_CONFIGURATION_URL": "",
			},
			nil,
			true,
		},
		{
			"UPSTREAM_TOKENINFO_URL and OPENID_PROVIDER_CONFIGURATION_URL empty",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "",
				"OPENID_PROVIDER_CONFIGURATION_URL": "",
			},
			nil,
			true,
		},
		{
			"UPSTREAM_TOKENINFO_URL empty, OPENID_PROVIDER_CONFIGURATION_URL set",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
			},
			nil,
			true,
		},
		{
			"UPSTREAM_TOKENINFO_URL set, OPENID_PROVIDER_CONFIGURATION_URL empty",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "",
			},
			nil,
			true,
		},
		{
			"REVOCATION_PROVIDER_URL empty",
			map[string]string{
				"REVOCATION_PROVIDER_URL": "",
			},
			nil,
			true,
		},
		{
			"REVOCATION_PROVIDER_URL set",
			map[string]string{
				"REVOCATION_PROVIDER_URL": "http://example.com",
			},
			nil,
			true,
		},
		{
			"UPSTREAM_TOKENINFO_URL, REVOCATION_PROVIDER_URL empty, OPENID_PROVIDER_CONFIGURATION_URL set",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "",
			},
			nil,
			true,
		},
		{
			"1",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "",
				"REVOCATION_PROVIDER_URL":           "",
			},
			nil,
			true,
		},
		// "2" was same as "4"
		{
			"3",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
			},
			nil,
			true,
		},
		{
			"4",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
			},
			&Settings{
				UpstreamTokenInfoURL:              nil,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"5",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"6",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"LISTEN_ADDRESS":                    ":80",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     ":80",
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"7",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"METRICS_LISTEN_ADDRESS":            ":80",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              ":80",
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"8",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"OPENID_PROVIDER_REFRESH_INTERVAL":  "1m",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     time.Minute,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"9",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"HTTP_CLIENT_TIMEOUT":               "1ms",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 time.Millisecond,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"10",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"UPSTREAM_CACHE_MAX_SIZE":           "123456789",
				"UPSTREAM_CACHE_TTL":                "17s",
				"UPSTREAM_TIMEOUT":                  "18s",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"HTTP_CLIENT_TLS_TIMEOUT":           "10ms",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              123456789,
				UpstreamCacheTTL:                  17 * time.Second,
				UpstreamTimeout:                   18 * time.Second,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              10 * time.Millisecond,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"11",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"UPSTREAM_CACHE_MAX_SIZE":           "0",
				"UPSTREAM_CACHE_TTL":                "0",
				"UPSTREAM_TIMEOUT":                  "0",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"HTTP_CLIENT_TLS_TIMEOUT":           "10ms",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              0,
				UpstreamCacheTTL:                  0,
				UpstreamTimeout:                   0,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              10 * time.Millisecond,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"12",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"REVOCATION_CACHE_TTL":              "10m0s",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                10 * time.Minute,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"13",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":               "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL":    "http://example.com",
				"REVOCATION_PROVIDER_URL":              "http://example.com",
				"REVOCATION_PROVIDER_REFRESH_INTERVAL": "30s",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: 30 * time.Second,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"14",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"REVOCATION_HASHING_SALT":           "TestSalt",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       "TestSalt",
				RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
		{
			"15",
			map[string]string{
				"UPSTREAM_TOKENINFO_URL":            "http://example.com",
				"OPENID_PROVIDER_CONFIGURATION_URL": "http://example.com",
				"REVOCATION_PROVIDER_URL":           "http://example.com",
				"REVOCATION_REFRESH_TOLERANCE":      "30s",
			},
			&Settings{
				UpstreamTokenInfoURL:              exampleCom,
				OpenIDProviderConfigurationURL:    exampleCom,
				RevocationProviderUrl:             exampleCom,
				UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
				UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
				UpstreamTimeout:                   defaultUpstreamTimeout,
				HTTPClientTimeout:                 defaultHTTPClientTimeout,
				HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
				OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
				ListenAddress:                     defaultListenAddress,
				MetricsListenAddress:              defaultMetricsListenAddress,
				RevocationCacheTTL:                defaultRevocationCacheTTL,
				RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
				HashingSalt:                       defaultHashingSalt,
				RevocationRefreshTolerance:        30 * time.Second,
				JwtProcessors:                     make(map[string]processor.JwtProcessor),
			},
			false,
		},
	} {
		os.Clearenv()
		for k, v := range test.env {
			os.Setenv(k, v)
		}
		err := LoadFromEnvironment()
		if test.wantFail {
			if err == nil {
				t.Errorf("TEST %s: Wanted failure to load settings but it seems that it succeeded: %+v", test.name, test)
			}
		} else {
			if !reflect.DeepEqual(AppSettings, test.want) {
				t.Errorf("TEST %s: Settings mismatch.\nWanted %+v\nGot %+v", test.name, test.want, AppSettings)
			}
		}
	}
}
