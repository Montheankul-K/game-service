package oauth2Exception

type OAuth2Processing struct{}

func (e *OAuth2Processing) Error() string {
	return "failed to processing oauth2"
}
