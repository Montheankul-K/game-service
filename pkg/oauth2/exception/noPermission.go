package oauth2Exception

type NoPermission struct{}

func (e *NoPermission) Error() string {
	return "no permission"
}
