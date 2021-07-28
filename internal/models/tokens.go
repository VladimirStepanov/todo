package models

type TokenDetails struct {
	AccessToken  string
	AccessET     int64
	AccessIAT    int64
	RefreshToken string
	RefreshET    int64
	RefreshIAT   int64
	UUID         string
}
