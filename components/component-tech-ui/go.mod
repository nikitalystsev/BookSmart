module github.com/nikitalystsev/BookSmart-tech-ui

go 1.22.5

require (
	github.com/google/uuid v1.6.0
	github.com/howeyc/gopass v0.0.0-20210920133722-c8aef6fb66ef
	github.com/jedib0t/go-pretty/v6 v6.5.9
	github.com/nikitalystsev/BookSmart-services v0.0.0-20240919123005-14b28ba85ee2
	github.com/nikitalystsev/BookSmart-web-api v0.0.0-20240921140007-b23de252cb67
)

require (
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/term v0.26.0 // indirect
)

replace (
	github.com/nikitalystsev/BookSmart-services => ../component-services
	github.com/nikitalystsev/BookSmart-web-api => ../component-web-api
)
