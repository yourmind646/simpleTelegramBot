package utils

import "github.com/jackc/pgx/v5/pgtype"

func ProcessRawUsername(rawUsername string) pgtype.Text {
	if rawUsername == "" {
		return pgtype.Text{Valid: false}
	} else {
		return pgtype.Text{String: "@" + rawUsername, Valid: true}
	}
}

func GetFullname(firstName string, lastName string) pgtype.Text {
	if firstName != "" && lastName != "" {
		return pgtype.Text{String: firstName + " " + lastName, Valid: true}
	} else if firstName != "" {
		return pgtype.Text{String: firstName, Valid: true}
	} else if lastName != "" {
		return pgtype.Text{String: lastName, Valid: true}
	} else {
		return pgtype.Text{String: "", Valid: true}
	}
}
