package utils

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func ToPgTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

func PgTextToString(pgText pgtype.Text) string {
	if !pgText.Valid {
		return ""
	}
	return pgText.String
}

func PgUUIDToUUID(u pgtype.UUID) (string, error) {
	if !u.Valid {
		return "", errors.New("error converting uuid")
	}

	id, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return "", errors.New("error decoding from bytes")
	}

	return id.String(), nil
}
