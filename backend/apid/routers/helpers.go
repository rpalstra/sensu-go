package routers

import "net/url"

func getGvrFromRoute(vars map[string]string) (string, string, string, error) {
	group, err := url.PathUnescape(vars["group"])
	if err != nil {
		return "", "", "", err
	}
	version, err := url.PathUnescape(vars["version"])
	if err != nil {
		return group, "", "", err
	}
	resource, err := url.PathUnescape(vars["resource"])
	if err != nil {
		return group, version, "", err
	}

	return group, version, resource, nil
}
