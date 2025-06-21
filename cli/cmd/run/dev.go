package run

import (
	"fmt"
	"os"
	"sync"

	"github.com/UchaBokeria/goyard/cli/utils"
)

var (
	Host string
	Port int
	WG   sync.WaitGroup
)

func Dev() {
	WG.Add(4)
	go Air()
	go Templ()
	go Tailwind()
	go Assets()
	WG.Wait()
}

func Air() {
	defer WG.Done()
	o, e := utils.Exec(`go run github.com/air-verse/air@latest \
		--build.cmd "go build -o ./bin/app ./cmd/app" \
		--build.bin "./bin/app" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go, templ" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true
	`)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Failed to run air: %v\n", e)
		os.Exit(1)
	}
	fmt.Println(o)

	//	if err := os.WriteFile(".air.toml", []byte(`root = "."
	//
	// tmp_dir = "tmp"
	// [build]
	//
	//	cmd = "go build -o ./tmp/main ."
	//	bin = "./tmp/main"
	//	delay = 1000
	//	exclude_dir = ["assets", "tmp", "vendor"]
	//	include_ext = ["go", "tpl", "tmpl", "templ", "html"]
	//	exclude_regex = ["_test\\.go"]
	//
	// [screen]
	//
	//	clear_on_rebuild = true
	//
	//	`), 0644); err != nil {
	//			fmt.Fprintf(os.Stderr, "Failed to create .air.toml: %v\n", err)
	//		}
}

func Templ() {
	defer WG.Done()
	o, e := utils.Exec(fmt.Sprintf(`go run github.com/a-h/templ/cmd/templ@latest \
	generate \
	--open-browser=false \
	--watch \
	--proxy="http://%s:%d" \
	--proxyport="%d" \
	--proxybind="%s"
	`, Host, Port, 7331, Host))
	if e != nil {
		fmt.Fprintf(os.Stderr, "Failed to run templ: %v\n", e)
		os.Exit(1)
	}
	fmt.Println(o)
}

func Tailwind() {
	defer WG.Done()
	o, e := utils.Exec(`bunx --yes tailwindcss -i ./public/assets/styles/tailwind.css -o ./public/assets/styles/style.css --watch`)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Failed to run tailwind: %v\n", e)
		os.Exit(1)
	}
	fmt.Println(o)
}

func Assets() {
	defer WG.Done()
}
