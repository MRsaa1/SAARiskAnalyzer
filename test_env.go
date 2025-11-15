package main

import (
	"fmt"
	"github.com/reserveone/saa-risk-analyzer/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Printf("DB Name: '%s'\n", cfg.DB.Name)
	fmt.Printf("DB Host: '%s'\n", cfg.DB.Host)
	fmt.Printf("DB User: '%s'\n", cfg.DB.User)
}
