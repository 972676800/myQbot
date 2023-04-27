package db

import remember2 "Qbot_gocode/Mods/remember"

func GetRememberFromDatabase() error {
	err := DB.Table(remember2.R.TableName()).First(&remember2.R).Error
	if err != nil {
		return err
	}
	return nil
}
