package traefik

import "fmt"

func Enable() string {
	return "traefik.enable=true"
}

func Host(routerName, host string) string {
	return fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s`)", routerName, host)
}

func EnableTls(routerName string) string {
	return fmt.Sprintf("traefik.http.routers.%s.tls=true", routerName)
}

func TslResolver(routerName string) string {
	return fmt.Sprintf("traefik.http.routers.%s.tls.certresolver=letsencrypt", routerName)
}

func LoadBalancerPort(routerName string) string {
	return fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=80", routerName)
}
