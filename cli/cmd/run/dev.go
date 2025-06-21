package run

import (
	"strconv"
	"sync"

	"github.com/UchaBokeria/goyard/cli/utils"
)

var (
	Host string
	Port int
)

func Dev() {
	return
	var wg sync.WaitGroup
	commands := [][]string{
		{"go", "run", "github.com/air-verse/air@latest",
			"--build.cmd", "go build -o ./bin/app ./cmd/app",
			"--build.bin", "./bin/app",
			"--build.delay", "100",
			"--build.exclude_dir", "node_modules",
			"--build.include_ext", "go, templ",
			"--build.stop_on_error", "false",
			"--misc.clean_on_exit", "true",
		},
		{"go", "run", "github.com/a-h/templ/cmd/templ@latest",
			"generate",
			"--open-browser=false",
			"--watch",
			"--proxy=" + Host + ":" + strconv.Itoa(Port),
			"--proxyport=7331",
			"--proxybind=" + Host,
		},
		{"bunx", "--yes", "tailwindcss", "-i", "./public/assets/styles/tailwind.css", "-o", "./public/assets/styles/style.css", "--watch"},
	}

	wg.Add(len(commands))
	for _, cmd := range commands {
		go func(c []string) {
			defer wg.Done()
			utils.Exec(c[0], c[1:]...)
		}(cmd)
	}

	wg.Wait()
}
