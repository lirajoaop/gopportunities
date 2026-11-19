package service

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/lirajoaop/gopportunities/schemas"
)

var errDB = errors.New("db failure")

// fakeRepo is an in-memory OpeningRepository whose failures can be injected per method.
type fakeRepo struct {
	openings map[string]*schemas.Opening

	createErr, findErr, findAllErr, saveErr, deleteErr error

	created, saved, deleted *schemas.Opening
}

func newFakeRepo(openings ...schemas.Opening) *fakeRepo {
	r := &fakeRepo{openings: map[string]*schemas.Opening{}}
	for i := range openings {
		o := openings[i]
		r.openings[idOf(&o)] = &o
	}
	return r
}

func idOf(o *schemas.Opening) string {
	return strconv.FormatUint(uint64(o.ID), 10)
}

func (r *fakeRepo) Create(o *schemas.Opening) error {
	r.created = o
	return r.createErr
}

func (r *fakeRepo) FindByID(id string) (*schemas.Opening, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	o, ok := r.openings[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	cp := *o
	return &cp, nil
}

func (r *fakeRepo) FindAll() ([]schemas.Opening, error) {
	if r.findAllErr != nil {
		return nil, r.findAllErr
	}
	out := make([]schemas.Opening, 0, len(r.openings))
	for _, o := range r.openings {
		out = append(out, *o)
	}
	return out, nil
}

func (r *fakeRepo) Save(o *schemas.Opening) error {
	r.saved = o
	return r.saveErr
}

func (r *fakeRepo) Delete(o *schemas.Opening) error {
	r.deleted = o
	return r.deleteErr
}

func ptr[T any](v T) *T { return &v }

func baseOpening() schemas.Opening {
	o := schemas.Opening{
		Role:     "Backend Developer",
		Company:  "Acme",
		Location: "Remote",
		Remote:   true,
		Link:     "https://acme.example/jobs/1",
		Salary:   100000,
	}
	o.ID = 1
	return o
}

func TestCreateOpening(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "persists the opening", repoErr: nil, wantErr: nil},
		{name: "propagates repository error", repoErr: errDB, wantErr: errDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			repo.createErr = tt.repoErr
			svc := NewOpeningService(repo)
			opening := baseOpening()

			err := svc.CreateOpening(&opening)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if repo.created != &opening {
				t.Errorf("repository received %p, want the same opening %p", repo.created, &opening)
			}
		})
	}
}

func TestGetOpeningByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		findErr error
		wantErr string
	}{
		{name: "returns existing opening", id: "1"},
		{name: "unknown id is reported as not found", id: "9", wantErr: "opening not found"},
		{name: "repository error is reported as not found", id: "1", findErr: errDB, wantErr: "opening not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo(baseOpening())
			repo.findErr = tt.findErr
			svc := NewOpeningService(repo)

			got, err := svc.GetOpeningByID(tt.id)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("got %+v, want nil opening on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if want := baseOpening(); !reflect.DeepEqual(*got, want) {
				t.Errorf("got %+v, want %+v", *got, want)
			}
		})
	}
}

func TestListOpenings(t *testing.T) {
	t.Run("returns all openings", func(t *testing.T) {
		second := baseOpening()
		second.ID = 2
		svc := NewOpeningService(newFakeRepo(baseOpening(), second))

		got, err := svc.ListOpenings()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("len = %d, want 2", len(got))
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repo := newFakeRepo()
		repo.findAllErr = errDB
		svc := NewOpeningService(repo)

		if _, err := svc.ListOpenings(); !errors.Is(err, errDB) {
			t.Fatalf("err = %v, want %v", err, errDB)
		}
	})
}

type updateInput struct {
	role, company, location, link string
	remote                        *bool
	salary                        int64
}

func TestUpdateOpening(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		in      updateInput
		findErr error
		saveErr error
		want    func(o *schemas.Opening) // mutates baseOpening into the expected result
		wantErr string
	}{
		{
			name: "updates every provided field",
			id:   "1",
			in: updateInput{
				role: "Staff Engineer", company: "Globex", location: "Lisbon",
				link: "https://globex.example/jobs/9", remote: ptr(false), salary: 200000,
			},
			want: func(o *schemas.Opening) {
				o.Role, o.Company, o.Location = "Staff Engineer", "Globex", "Lisbon"
				o.Link, o.Remote, o.Salary = "https://globex.example/jobs/9", false, 200000
			},
		},
		{
			name: "keeps fields that are empty or zero",
			id:   "1",
			in:   updateInput{role: "Senior Backend Developer"},
			want: func(o *schemas.Opening) { o.Role = "Senior Backend Developer" },
		},
		{
			name: "explicit remote=false overrides true",
			id:   "1",
			in:   updateInput{remote: ptr(false)},
			want: func(o *schemas.Opening) { o.Remote = false },
		},
		{
			name: "non-positive salary is ignored",
			id:   "1",
			in:   updateInput{company: "Initech", salary: -5},
			want: func(o *schemas.Opening) { o.Company = "Initech" },
		},
		{
			name:    "unknown id is reported as not found",
			id:      "9",
			in:      updateInput{role: "x"},
			wantErr: "opening not found",
		},
		{
			name:    "repository find error is reported as not found",
			id:      "1",
			in:      updateInput{role: "x"},
			findErr: errDB,
			wantErr: "opening not found",
		},
		{
			name:    "save error is wrapped",
			id:      "1",
			in:      updateInput{role: "x"},
			saveErr: errDB,
			wantErr: "error updating opening",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo(baseOpening())
			repo.findErr, repo.saveErr = tt.findErr, tt.saveErr
			svc := NewOpeningService(repo)

			got, err := svc.UpdateOpening(tt.id, tt.in.role, tt.in.company, tt.in.location, tt.in.link, tt.in.remote, tt.in.salary)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("got %+v, want nil opening on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := baseOpening()
			tt.want(&want)
			if !reflect.DeepEqual(*got, want) {
				t.Errorf("got  %+v\nwant %+v", *got, want)
			}
			if repo.saved != got {
				t.Errorf("saved opening %p is not the returned opening %p", repo.saved, got)
			}
		})
	}
}

func TestUpdateOpeningDoesNotSaveWhenNotFound(t *testing.T) {
	repo := newFakeRepo()
	svc := NewOpeningService(repo)

	_, _ = svc.UpdateOpening("1", "role", "", "", "", nil, 0)

	if repo.saved != nil {
		t.Errorf("Save was called for a missing opening")
	}
}

func TestDeleteOpening(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		findErr   error
		deleteErr error
		wantErr   string
	}{
		{name: "deletes and returns the opening", id: "1"},
		{name: "unknown id is reported as not found", id: "9", wantErr: "opening not found"},
		{name: "repository find error is reported as not found", id: "1", findErr: errDB, wantErr: "opening not found"},
		{name: "delete error is wrapped", id: "1", deleteErr: errDB, wantErr: "error deleting opening"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo(baseOpening())
			repo.findErr, repo.deleteErr = tt.findErr, tt.deleteErr
			svc := NewOpeningService(repo)

			got, err := svc.DeleteOpening(tt.id)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("got %+v, want nil opening on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.deleted != got {
				t.Errorf("deleted opening %p is not the returned opening %p", repo.deleted, got)
			}
			if got.ID != 1 {
				t.Errorf("ID = %d, want 1", got.ID)
			}
		})
	}
}
