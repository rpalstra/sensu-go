package routers

import "net/url"

func getGvrFromRoute(vars map[string]string) (group, version, resource string, err error) {
	group, err = url.PathUnescape(vars["group"])
	if err != nil {
		return
	}

	version, err = url.PathUnescape(vars["version"])
	if err != nil {
		return
	}

	resource, err = url.PathUnescape(vars["resource"])
	if err != nil {
		return
	}

	return
}
