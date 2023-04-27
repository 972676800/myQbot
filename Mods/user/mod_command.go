package user

import (
	"Qbot_gocode/db"
	"fmt"
	"gorm.io/gorm"
)

var Commands map[string]string
var commands_slice []string

type Command struct {
	gorm.Model
	Command     string `gorm:"uniqueIndex"`
	Description string
}

func (Command) TableName() string {
	return "t_command"
}

func init() {
	db.DB.AutoMigrate(&Command{})
	Commands = make(map[string]string)
	GetCommandFromDatabase()
}

func GetCommandFromDatabase() {
	var commandsSlice []Command
	if err := db.DB.Find(&commandsSlice).Error; err != nil {
		fmt.Printf("Failed to get command records: %v\n", err)
	}
	commands_slice = make([]string, len(commandsSlice))
	for k, v := range commandsSlice {
		Commands[v.Command] = v.Description
		commands_slice[k] = v.Command + v.Description
	}
}
