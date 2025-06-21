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
	processes := []func() (string, error){
		Air,
		Templ,
		Tailwind,
	}
	WG.Add(len(processes))
	for _, process := range processes {
		go func(process func() (string, error)) {
			defer WG.Done()
			o, e := process()
			if e != nil {
				fmt.Fprintf(os.Stderr, "Failed to run %s: %v\n", process, e)
				os.Exit(1)
			}
			fmt.Println(o)
		}(process)
	}
	WG.Wait()
}

func Air() (string, error) {
	return utils.Exec(`go run github.com/air-verse/air@latest \
		--build.cmd "go build -o ./bin/app ./cmd/app" \
		--build.bin "./bin/app" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go, templ" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true
	`)
}

func Templ() (string, error) {
	return utils.Exec(fmt.Sprintf(`go run github.com/a-h/templ/cmd/templ@latest \
	generate \
	--open-browser=false \
	--watch \
	--proxy="http://%s:%d" \
	--proxyport="%d" \
	--proxybind="%s"
	`, Host, Port, 7331, Host))
}

func Tailwind() (string, error) {
	return utils.Exec(`bunx --yes tailwindcss -i ./public/assets/styles/tailwind.css -o ./public/assets/styles/style.css --watch`)
}
