module github.com/chuakid/cvwo-backend

// +heroku goVersion go1.16
go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e
	gorm.io/driver/postgres v1.2.3
	gorm.io/gorm v1.22.4
)
