module github.com/AndersKaae/go_virk_api

go 1.25.3

require (
	github.com/AndersKaae/go_virk_updater v0.0.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/joho/godotenv v1.5.1
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/AndersKaae/go_virk_updater => ../go_virk_updater
