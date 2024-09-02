module github.com/HyPE-Network/vanilla-proxy

go 1.22

toolchain go1.22.4

require (
	github.com/go-gl/mathgl v1.1.0
	github.com/gofrs/flock v0.12.1
	github.com/google/uuid v1.6.0
	github.com/pelletier/go-toml v1.9.5
	github.com/sandertv/go-raknet v1.14.1
	github.com/sandertv/gophertunnel v1.40.0
	github.com/sirupsen/logrus v1.9.3
	github.com/tailscale/hujson v0.0.0-20221223112325-20486734a56a
)

require (
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/muhammadmuzzammil1998/jsonc v1.0.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/image v0.19.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/oauth2 v0.22.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)

replace github.com/sandertv/go-raknet => github.com/smell-of-curry/go-raknet v0.0.0-20240831142756-af42da4d0d54

replace github.com/sandertv/gophertunnel => github.com/smell-of-curry/gophertunnel v1.39.1-0.20240902142224-951735410358
