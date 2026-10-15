package controller

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretCacheOptions limits only Secret informers, leaving ExternalMonitor discovery unchanged.
func SecretCacheOptions(namespaces string, disabled bool) (cache.Options, error) {
	if disabled && namespaces != "" {
		return cache.Options{}, fmt.Errorf("secret-watch-namespaces and disable-secret-watch are mutually exclusive")
	}
	if namespaces == "" {
		return cache.Options{}, nil
	}
	allowed := map[string]cache.Config{}
	for _, value := range strings.Split(namespaces, ",") {
		ns := strings.TrimSpace(value)
		if errors := validation.IsDNS1123Label(ns); len(errors) != 0 {
			return cache.Options{}, fmt.Errorf("invalid Secret watch namespace %q: %s", ns, strings.Join(errors, "; "))
		}
		allowed[ns] = cache.Config{}
	}
	return cache.Options{ByObject: map[client.Object]cache.ByObject{&corev1.Secret{}: {Namespaces: allowed}}}, nil
}
