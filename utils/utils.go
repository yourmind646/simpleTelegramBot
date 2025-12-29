package utils

import (
	"fmt"
	"html"

	"github.com/jackc/pgx/v5/pgtype"
)

var AvailableFileKeys []string = []string{
	// UI
	"profilePhoto", "inventoryPhoto", "sortiePhoto",
	// locations
	"locYardsPhoto", "locSupermarketPhoto", "locMetroPhoto", "locArmyBasePhoto",
}

func ProcessRawUsername(rawUsername string) pgtype.Text {
	if rawUsername == "" {
		return pgtype.Text{Valid: false}
	} else {
		return pgtype.Text{String: "@" + rawUsername, Valid: true}
	}
}

func GetFullname(firstName string, lastName string) pgtype.Text {
	var fullName string

	if firstName != "" && lastName != "" {
		fullName = html.EscapeString(firstName + " " + lastName)
	} else if firstName != "" {
		fullName = html.EscapeString(firstName)
	} else if lastName != "" {
		fullName = html.EscapeString(lastName)
	} else {
		fullName = ""
	}

	return pgtype.Text{String: fullName, Valid: true}
}

func IsFileKeyExists(fileKey string) bool {
	for _, afileKey := range AvailableFileKeys {
		if afileKey == fileKey {
			return true
		}
	}
	return false
}

func FormatSecondsToString(seconds int) string {
	minutes := 0
	hours := 0

	if seconds/3600 > 0 {
		hours += seconds / 3600
		seconds -= hours * 3600
	}
	if seconds/60 > 0 {
		minutes += seconds / 60
		seconds -= minutes * 60
	}

	return fmt.Sprintf("%dч %dм %dс", hours, minutes, seconds)
}

func FormatDropList(drops map[string]map[string]int) string {
	result := ""

	for itemKey, dropOpts := range drops {
		result += fmt.Sprintf("• %s (до %d шт.) - %d%%\n", itemKey, dropOpts["maxN"], dropOpts["pct"])
	}

	return result
}
