module github.com/noxworld-dev/opennox-lib

go 1.21

require (
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b
	github.com/go-gl/mathgl v1.0.0
	github.com/icza/bitio v1.1.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/noxworld-dev/noxcrypt v0.0.0-20230831140413-02623e75408e
	github.com/noxworld-dev/noxscript/eud/v171 v171.3.1
	github.com/noxworld-dev/noxscript/ns v1.0.2
	github.com/noxworld-dev/noxscript/ns/v3 v3.4.2
	github.com/noxworld-dev/noxscript/ns/v4 v4.14.0
	github.com/shoenig/test v1.7.0
	github.com/spf13/cobra v1.7.0
	github.com/traefik/yaegi v0.15.1
	github.com/veandco/go-sdl2 v0.5.0-alpha.4
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	golang.org/x/image v0.11.0
	golang.org/x/sys v0.14.0
	golang.org/x/text v0.12.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.12.0 // indirect
)

replace github.com/traefik/yaegi => github.com/dennwc/yaegi v0.14.4-0.20230902105020-22e112bbf6e8
