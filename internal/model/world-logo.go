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
	Id            *string
	Name          string
	LogoPath      *string
	LogoBase64Str *string
	SrcKey        string
}

type WorldLogosQueryOptions struct {
	Search *string
	Ids    []string
}
