package configs

import (
	"path/filepath"
	"os"
	"path"
	"log"
	"github.com/joho/godotenv"
	"strconv"
)

type Configs struct {
	WorkDir				string
	WebhookURL			string
	DiskWarning			float64
	SlackTitle 			string
	SlackTitleLink		string
	SlackText  			string
	SlackOKColor  		string
	SlackWarningColor	string
	SlackFoooter		string
}

func GetConfigs() Configs {
	work_dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	err := godotenv.Load(path.Join(work_dir, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	diskWarning, err := strconv.ParseFloat(os.Getenv("DISK_WARNING"), 64)
	if err != nil {
		log.Fatal("Error loading DISK_WARNING variable")
	}
	config := Configs {
		WorkDir: work_dir,
		WebhookURL: os.Getenv("WEBHOOK_URL"),
		DiskWarning: diskWarning,
		SlackTitle: os.Getenv("SLACK_TITLE"),
		SlackTitleLink: os.Getenv("SLACK_TITLE_LINK"),
		SlackText: os.Getenv("SLACK_TEXT"),
		SlackOKColor: os.Getenv("SLACK_OK_COLOR"),
		SlackWarningColor: os.Getenv("SLACK_WARING_COLOR"),
		SlackFoooter: os.Getenv("SLACK_FOOTER"),
	}
	return config
}