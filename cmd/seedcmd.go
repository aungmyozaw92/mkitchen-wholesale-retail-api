package cmd

import (
	"fmt"
	"os"

	"github.com/myanmarmarathon/mkitchen-distribution-backend/seeder"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Your Application Description",
}

var seedingCommand = &cobra.Command{
    Use:   "seed",
    Short: "Seed the database with mock data",
    Run: func(cmd *cobra.Command, args []string) {

		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

		DB, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{})
        if err != nil {
            fmt.Println("Failed to connect to the database")
            return
        }

        seeder.SeedDatabase(DB)
        fmt.Println("Database seeded successfully")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        if err != nil {
            fmt.Println("Execute error")
            return
        }
    }
}

func init() {
    rootCmd.AddCommand(seedingCommand)
}