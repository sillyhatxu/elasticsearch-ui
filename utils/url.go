package utils

import "net/url"

func FormatURL(Scheme, URL string) string {
	u := url.URL{
		Scheme: Scheme,
		Path:   URL,
	}
	return u.String()
}
