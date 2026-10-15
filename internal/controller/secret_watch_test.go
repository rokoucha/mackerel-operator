package controller

import (
	"reflect"
	"testing"
)

func TestSecretCacheOptions(t *testing.T) {
	for _, tt := range []struct {
		name, namespaces  string
		disabled, invalid bool
		want              []string
	}{
		{name: "legacy all namespaces"},
		{name: "disabled", disabled: true},
		{name: "scoped", namespaces: "argocd, authentik,argocd", want: []string{"argocd", "authentik"}},
		{name: "conflicting", namespaces: "argocd", disabled: true, invalid: true},
		{name: "empty entry", namespaces: "argocd,", invalid: true},
		{name: "invalid namespace", namespaces: "UPPER", invalid: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := SecretCacheOptions(tt.namespaces, tt.disabled)
			if (err != nil) != tt.invalid {
				t.Fatalf("error=%v, invalid=%v", err, tt.invalid)
			}
			if err != nil {
				return
			}
			if opts.DefaultNamespaces != nil {
				t.Fatal("ExternalMonitor discovery must remain cluster-wide")
			}
			if len(tt.want) == 0 {
				if len(opts.ByObject) != 0 {
					t.Fatal("unexpected Secret cache override")
				}
				return
			}
			if len(opts.ByObject) != 1 {
				t.Fatalf("cache overrides=%d", len(opts.ByObject))
			}
			for _, obj := range opts.ByObject {
				want := map[string]bool{}
				for _, ns := range tt.want {
					want[ns] = true
				}
				got := map[string]bool{}
				for ns := range obj.Namespaces {
					got[ns] = true
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("namespaces=%v, want=%v", got, want)
				}
			}
		})
	}
}
