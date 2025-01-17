package model

import "time"

type WorldLogo struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	LogoPath  string    `db:"logo_path"`
	SrcKey    string    `db:"src_key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type WorldLogoInput struct {
	Id       string
	SrcKey   string
	Name     string
	LogoPath string
}

type SaveWorldLogoInput struct {
	SrcKey        string
	Name          string
	LogoBase64Str string
	ContentType   string
	FileExtension string
}

type WorldLogosQueryOptions struct {
	Search *string
	Ids    []string
}
