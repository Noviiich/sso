package main

import (
	"fmt"

	"github.com/Noviiich/sso/internal/config"
)

func main() {
	cfd := config.MustLoad()
	fmt.Printf("Loaded config: %+v\n", cfd)
}
