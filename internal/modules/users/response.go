package users

type CreatePATResponse struct {
	Token *PersonalAccessToken `json:"token"`
	Plain string               `json:"plain"`
}

type AppearanceResponse struct {
	Appearance *AppearancePreferences `json:"appearance"`
}

type SetupStatusResponse struct {
	SetupRequired bool   `json:"setupRequired"`
	Cwd           string `json:"cwd"`
}
