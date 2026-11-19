package schemas

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestOpeningToResponse(t *testing.T) {
	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := created.Add(time.Hour)

	o := Opening{
		Model:    gorm.Model{ID: 7, CreatedAt: created, UpdatedAt: updated},
		Role:     "Backend Developer",
		Company:  "Acme",
		Location: "Remote",
		Remote:   true,
		Link:     "https://acme.example/jobs/7",
		Salary:   120000,
	}

	got := o.ToResponse()
	want := OpeningResponse{
		ID: 7, CreatedAt: created, UpdatedAt: updated,
		Role: "Backend Developer", Company: "Acme", Location: "Remote",
		Remote: true, Link: "https://acme.example/jobs/7", Salary: 120000,
	}
	if got != want {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}

	body, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "deletedAt") {
		t.Errorf("deletedAt should be omitted for a live record: %s", body)
	}
}

func TestOpeningToResponseSoftDeleted(t *testing.T) {
	deleted := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	o := Opening{Model: gorm.Model{ID: 1, DeletedAt: gorm.DeletedAt{Time: deleted, Valid: true}}}

	got := o.ToResponse()

	if got.DeletedAt == nil || !got.DeletedAt.Equal(deleted) {
		t.Errorf("DeletedAt = %v, want %v", got.DeletedAt, deleted)
	}
}
