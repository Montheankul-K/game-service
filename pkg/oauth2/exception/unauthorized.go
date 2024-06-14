package oauth2Exception

type Unauthorized struct{}

func (e *Unauthorized) Error() string {
	return "unauthorized"
}
