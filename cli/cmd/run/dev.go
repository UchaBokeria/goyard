package run

import (
	"fmt"
	"os"
	"sync"

	"github.com/UchaBokeria/goyard/cli/utils"
)

func Dev(Host string, Port int) {
	processes := []struct {
		Name string
		Cmd  string
	}{
		{"AIR", `go run github.com/air-verse/air@latest --build.cmd "go build -o ./bin/app ./cmd/app" --build.bin "./bin/app" --build.delay "100" --build.exclude_dir "node_modules" --build.include_ext "go, templ" --build.stop_on_error "false" --misc.clean_on_exit true`},
		{"TEMPL", fmt.Sprintf(`go run github.com/a-h/templ/cmd/templ@latest generate --open-browser=false --watch --proxy="http://%s:%d" --proxyport="%d" --proxybind="%s"`, Host, Port, 7331, Host)},
		{"TAILWIND", `bunx --yes tailwindcss -i ./public/assets/styles/tailwind.css -o ./public/assets/styles/style.css --watch`},
	}

	var wg sync.WaitGroup
	for _, p := range processes {
		wg.Add(1)
		go func(name, cmd string) {
			defer wg.Done()
			outChan, errChan := utils.ExecLive(name, cmd)

			for line := range outChan {
				fmt.Println(line)
			}

			if err := <-errChan; err != nil {
				fmt.Fprintf(os.Stderr, "[%s] exited with error: %v\n", name, err)
			}
		}(p.Name, p.Cmd)
	}

	wg.Wait()
}
