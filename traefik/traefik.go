package traefik

import (
	"fmt"
	"strings"
)

const (
	True   = "true"
	Enable = "traefik.enable"
)

func Rule(routerName string) string {
	return fmt.Sprintf("traefik.http.routers.%s.rule", routerName)
}

func Host(host string) string {
	return fmt.Sprintf("Host(`%s`)", strings.ToLower(host))
}

func TLS(routerName string) string {
	return fmt.Sprintf("traefik.http.routers.%s.tls", routerName)
}

func TLSResolver(routerName string) string {
	return fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", routerName)
}

func LoadBalancerPort(routerName string) string {
	return fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", routerName)
}
