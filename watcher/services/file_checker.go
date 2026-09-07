package services

import (
	"log/slog"
	"os"
	"strings"
	"time"

	repo "watcher/db"
	"watcher/models"
)

func VerifyFileNotChanged(repo *repo.Sqlyte, fileType string, path string, originalContent string, revertChanges bool) {

	isError := true
	description := ""

	content, err := os.ReadFile("")
	if err != nil {
		description = "Failed to read the file: "
		slog.Error(description, "error", err)
		return
	}

	currentText := string(content)
	currentText = strings.ReplaceAll(currentText, " ", "")
	currentText = strings.ReplaceAll(currentText, "\n", "")
	currentText = strings.ReplaceAll(currentText, "\r", "")

	if currentText != originalContent {
		description = "Content has changed."

		if revertChanges {
			err := os.WriteFile(path, []byte(originalContent), 0600)
			if err != nil {
				slog.Error(err.Error())
			}

			s := strings.Split(path, "/")

			file, err := os.Create("backup-" + s[len(s)-1])
			if err != nil {
				slog.Error(err.Error())
			} else {
				defer file.Close()

				_, err = file.WriteString(originalContent)
				if err != nil {
					slog.Error(err.Error())
				}
			}
		}

	} else {
		isError = false
	}

	repo.SaveRecord(models.Record{
		App:         fileType,
		Message:     "",
		Description: description,
		IsError:     isError,
		CreatedAt:   time.Now(),
	})

}
