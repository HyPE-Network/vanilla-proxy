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

replace github.com/sandertv/go-raknet => github.com/hashimthearab/go-raknet v1.14.2-0.20240712204703-9b99c862e9db

replace github.com/sandertv/gophertunnel => github.com/smell-of-curry/gophertunnel v1.39.1-0.20240820202257-325c4167b8cd
