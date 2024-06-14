package oauth2Exception

type Logout struct{}

func (e *Logout) Error() string {
	return "failed to logout"
}
