package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lirajoaop/gopportunities/config"
	"github.com/lirajoaop/gopportunities/schemas"
)

// fakeService is an OpeningService whose results are set per test.
type fakeService struct {
	opening  *schemas.Opening
	openings []schemas.Opening
	err      error

	createdWith  *schemas.Opening
	calledWithID string
	updateArgs   *updateArgs
}

type updateArgs struct {
	role, company, location, link string
	remote                        *bool
	salary                        int64
}

func (f *fakeService) CreateOpening(o *schemas.Opening) error {
	f.createdWith = o
	if f.err == nil {
		o.ID = 42
	}
	return f.err
}

func (f *fakeService) GetOpeningByID(id string) (*schemas.Opening, error) {
	f.calledWithID = id
	return f.opening, f.err
}

func (f *fakeService) ListOpenings() ([]schemas.Opening, error) {
	return f.openings, f.err
}

func (f *fakeService) UpdateOpening(id string, role, company, location, link string, remote *bool, salary int64) (*schemas.Opening, error) {
	f.calledWithID = id
	f.updateArgs = &updateArgs{role, company, location, link, remote, salary}
	return f.opening, f.err
}

func (f *fakeService) DeleteOpening(id string) (*schemas.Opening, error) {
	f.calledWithID = id
	return f.opening, f.err
}

// envelope mirrors the JSON shape produced by sendSuccess and sendError.
type envelope struct {
	Message   string          `json:"message"`
	ErrorCode int             `json:"errorCode"`
	Data      json.RawMessage `json:"data"`
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// newRouter wires svc into the package-level handler state and returns a router
// with the same routes as router/routes.go. The previous state is restored on cleanup.
func newRouter(t *testing.T, svc OpeningService) *gin.Engine {
	t.Helper()
	prevSvc, prevLogger := openingService, logger
	openingService, logger = svc, config.NewLogger("handler-test")
	t.Cleanup(func() { openingService, logger = prevSvc, prevLogger })

	r := gin.New()
	v1 := r.Group("/api/v1")
	v1.GET("/openings", ListOpeningsHandler)
	v1.GET("/openings/:id", ShowOpeningHandler)
	v1.POST("/openings", CreateOpeningHandler)
	v1.PUT("/openings/:id", UpdateOpeningHandler)
	v1.DELETE("/openings/:id", DeleteOpeningHandler)
	return r
}

func do(t *testing.T, r http.Handler, method, path, body string) (*httptest.ResponseRecorder, envelope) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	var env envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("response is not valid JSON: %v\nbody: %s", err, rec.Body.String())
	}
	return rec, env
}

func sampleOpening() *schemas.Opening {
	o := &schemas.Opening{
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

func decodeOpening(t *testing.T, raw json.RawMessage) schemas.OpeningResponse {
	t.Helper()
	var got schemas.OpeningResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("data is not an opening: %v\ndata: %s", err, raw)
	}
	return got
}

const validCreateBody = `{"role":"Backend Developer","company":"Acme","location":"Remote","remote":false,"link":"https://acme.example/jobs/1","salary":100000}`

func TestCreateOpeningHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
		wantMsg    string
	}{
		{name: "creates opening", body: validCreateBody, wantStatus: http.StatusOK, wantMsg: "operation create-opening successful"},
		{name: "empty body", body: `{}`, wantStatus: http.StatusBadRequest, wantMsg: "request body is empty or invalid"},
		{name: "malformed JSON", body: `{"role":`, wantStatus: http.StatusBadRequest, wantMsg: "request body is empty or invalid"},
		{name: "missing role", body: `{"company":"Acme","location":"Remote","remote":true,"link":"l","salary":1}`, wantStatus: http.StatusBadRequest, wantMsg: "param: role (type: string) is required"},
		{name: "missing remote", body: `{"role":"r","company":"Acme","location":"Remote","link":"l","salary":1}`, wantStatus: http.StatusBadRequest, wantMsg: "param: remote (type: bool) is required"},
		{name: "service error", body: validCreateBody, svcErr: errors.New("db down"), wantStatus: http.StatusInternalServerError, wantMsg: "error creating opening on database"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{err: tt.svcErr}
			rec, env := do(t, newRouter(t, svc), http.MethodPost, "/api/v1/openings", tt.body)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if env.Message != tt.wantMsg {
				t.Errorf("message = %q, want %q", env.Message, tt.wantMsg)
			}
			if tt.wantStatus != http.StatusOK {
				if env.ErrorCode != tt.wantStatus {
					t.Errorf("errorCode = %d, want %d", env.ErrorCode, tt.wantStatus)
				}
				if tt.wantStatus == http.StatusBadRequest && svc.createdWith != nil {
					t.Errorf("service was called for an invalid request")
				}
				return
			}

			got := decodeOpening(t, env.Data)
			if got.ID != 42 || got.Role != "Backend Developer" || got.Remote != false || got.Salary != 100000 {
				t.Errorf("unexpected opening in response: %+v", got)
			}
		})
	}
}

func TestShowOpeningHandler(t *testing.T) {
	tests := []struct {
		name       string
		opening    *schemas.Opening
		svcErr     error
		wantStatus int
		wantMsg    string
	}{
		{name: "returns opening", opening: sampleOpening(), wantStatus: http.StatusOK, wantMsg: "operation show-opening successful"},
		{name: "not found", svcErr: errors.New("opening not found"), wantStatus: http.StatusNotFound, wantMsg: "opening not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{opening: tt.opening, err: tt.svcErr}
			rec, env := do(t, newRouter(t, svc), http.MethodGet, "/api/v1/openings/1", "")

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if env.Message != tt.wantMsg {
				t.Errorf("message = %q, want %q", env.Message, tt.wantMsg)
			}
			if svc.calledWithID != "1" {
				t.Errorf("service called with id %q, want %q", svc.calledWithID, "1")
			}
			if tt.wantStatus == http.StatusOK {
				if got := decodeOpening(t, env.Data); got.ID != 1 || got.Company != "Acme" {
					t.Errorf("unexpected opening in response: %+v", got)
				}
			}
		})
	}
}

func TestListOpeningsHandler(t *testing.T) {
	t.Run("returns openings", func(t *testing.T) {
		second := *sampleOpening()
		second.ID = 2
		svc := &fakeService{openings: []schemas.Opening{*sampleOpening(), second}}
		rec, env := do(t, newRouter(t, svc), http.MethodGet, "/api/v1/openings", "")

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body)
		}
		var got []schemas.OpeningResponse
		if err := json.Unmarshal(env.Data, &got); err != nil {
			t.Fatalf("data is not a list of openings: %v", err)
		}
		if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
			t.Errorf("unexpected openings: %+v", got)
		}
	})

	t.Run("empty list is an empty JSON array", func(t *testing.T) {
		rec, env := do(t, newRouter(t, &fakeService{}), http.MethodGet, "/api/v1/openings", "")

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if string(env.Data) != "[]" {
			t.Errorf("data = %s, want []", env.Data)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeService{err: errors.New("db down")}
		rec, env := do(t, newRouter(t, svc), http.MethodGet, "/api/v1/openings", "")

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", rec.Code)
		}
		if env.Message != "error listing openings" {
			t.Errorf("message = %q", env.Message)
		}
	})
}

func TestUpdateOpeningHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		opening    *schemas.Opening
		svcErr     error
		wantStatus int
		wantMsg    string
	}{
		{name: "updates opening", body: `{"role":"Senior Backend Developer","remote":false}`, opening: sampleOpening(), wantStatus: http.StatusOK, wantMsg: "operation update-opening successful"},
		{name: "empty body", body: `{}`, wantStatus: http.StatusBadRequest, wantMsg: "at least one field must be provided to update"},
		{name: "malformed JSON", body: `not json`, wantStatus: http.StatusBadRequest, wantMsg: "at least one field must be provided to update"},
		{name: "not found", body: `{"role":"x"}`, svcErr: errors.New("opening not found"), wantStatus: http.StatusNotFound, wantMsg: "opening not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{opening: tt.opening, err: tt.svcErr}
			rec, env := do(t, newRouter(t, svc), http.MethodPut, "/api/v1/openings/1", tt.body)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if env.Message != tt.wantMsg {
				t.Errorf("message = %q, want %q", env.Message, tt.wantMsg)
			}
			if tt.wantStatus == http.StatusBadRequest && svc.updateArgs != nil {
				t.Errorf("service was called for an invalid request")
			}
		})
	}
}

func TestUpdateOpeningHandlerForwardsFields(t *testing.T) {
	svc := &fakeService{opening: sampleOpening()}
	body := `{"role":"r","company":"c","location":"l","link":"k","remote":false,"salary":7}`
	do(t, newRouter(t, svc), http.MethodPut, "/api/v1/openings/5", body)

	if svc.calledWithID != "5" {
		t.Errorf("id = %q, want 5", svc.calledWithID)
	}
	a := svc.updateArgs
	if a == nil {
		t.Fatal("service was not called")
	}
	if a.role != "r" || a.company != "c" || a.location != "l" || a.link != "k" || a.salary != 7 {
		t.Errorf("unexpected args: %+v", *a)
	}
	if a.remote == nil || *a.remote != false {
		t.Errorf("remote = %v, want pointer to false", a.remote)
	}
}

func TestDeleteOpeningHandler(t *testing.T) {
	tests := []struct {
		name       string
		opening    *schemas.Opening
		svcErr     error
		wantStatus int
		wantMsg    string
	}{
		{name: "deletes opening", opening: sampleOpening(), wantStatus: http.StatusOK, wantMsg: "operation delete-opening successful"},
		{name: "not found", svcErr: errors.New("opening not found"), wantStatus: http.StatusNotFound, wantMsg: "opening not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{opening: tt.opening, err: tt.svcErr}
			rec, env := do(t, newRouter(t, svc), http.MethodDelete, "/api/v1/openings/1", "")

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if env.Message != tt.wantMsg {
				t.Errorf("message = %q, want %q", env.Message, tt.wantMsg)
			}
			if svc.calledWithID != "1" {
				t.Errorf("service called with id %q, want 1", svc.calledWithID)
			}
		})
	}
}

// The handlers also guard against an empty id, which the /openings/:id routes
// never produce; mount them on a route without the param to exercise that branch.
func TestHandlersRejectMissingID(t *testing.T) {
	handlers := map[string]gin.HandlerFunc{
		http.MethodGet:    ShowOpeningHandler,
		http.MethodPut:    UpdateOpeningHandler,
		http.MethodDelete: DeleteOpeningHandler,
	}

	for method, h := range handlers {
		t.Run(method, func(t *testing.T) {
			svc := &fakeService{}
			r := newRouter(t, svc)
			r.Handle(method, "/no-id", h)

			rec, env := do(t, r, method, "/no-id", `{"role":"x"}`)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
			if want := "param: id (type: pathParameter) is required"; env.Message != want {
				t.Errorf("message = %q, want %q", env.Message, want)
			}
			if svc.calledWithID != "" || svc.updateArgs != nil {
				t.Errorf("service was called without an id")
			}
		})
	}
}
