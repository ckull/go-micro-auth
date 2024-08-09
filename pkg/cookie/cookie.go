

type (
	Cookie interface {
	}


)

func newCookie(name string, value string, expiredAt Time) http.Cookie {
	return &http.Cookie{
		Name: name,
		Value: value,
		Expires: expiredAt,
		HttpOnly: true,
		Path: path
	}
}