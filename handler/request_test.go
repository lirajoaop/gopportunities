package handler

import "testing"

func TestCreateOpeningRequestValidate(t *testing.T) {
	remote := true
	valid := func() CreateOpeningRequest {
		return CreateOpeningRequest{
			Role: "Backend Developer", Company: "Acme", Location: "Remote",
			Remote: &remote, Link: "https://acme.example/jobs/1", Salary: 100000,
		}
	}

	tests := []struct {
		name    string
		mutate  func(r *CreateOpeningRequest)
		wantErr string
	}{
		{name: "valid request", mutate: func(*CreateOpeningRequest) {}},
		{name: "empty request", mutate: func(r *CreateOpeningRequest) { *r = CreateOpeningRequest{} }, wantErr: "request body is empty or invalid"},
		{name: "missing role", mutate: func(r *CreateOpeningRequest) { r.Role = "" }, wantErr: "param: role (type: string) is required"},
		{name: "missing company", mutate: func(r *CreateOpeningRequest) { r.Company = "" }, wantErr: "param: company (type: string) is required"},
		{name: "missing location", mutate: func(r *CreateOpeningRequest) { r.Location = "" }, wantErr: "param: location (type: string) is required"},
		{name: "missing link", mutate: func(r *CreateOpeningRequest) { r.Link = "" }, wantErr: "param: link (type: string) is required"},
		{name: "missing remote", mutate: func(r *CreateOpeningRequest) { r.Remote = nil }, wantErr: "param: remote (type: bool) is required"},
		{name: "zero salary", mutate: func(r *CreateOpeningRequest) { r.Salary = 0 }, wantErr: "param: salary (type: int64) is required"},
		{name: "negative salary", mutate: func(r *CreateOpeningRequest) { r.Salary = -1 }, wantErr: "param: salary (type: int64) is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := valid()
			tt.mutate(&req)
			assertErr(t, req.Validate(), tt.wantErr)
		})
	}
}

func TestUpdateOpeningRequestValidate(t *testing.T) {
	f := false
	tests := []struct {
		name    string
		req     UpdateOpeningRequest
		wantErr string
	}{
		{name: "empty request", req: UpdateOpeningRequest{}, wantErr: "at least one field must be provided to update"},
		{name: "non-positive salary only", req: UpdateOpeningRequest{Salary: -1}, wantErr: "at least one field must be provided to update"},
		{name: "role only", req: UpdateOpeningRequest{Role: "r"}},
		{name: "company only", req: UpdateOpeningRequest{Company: "c"}},
		{name: "location only", req: UpdateOpeningRequest{Location: "l"}},
		{name: "link only", req: UpdateOpeningRequest{Link: "k"}},
		{name: "remote=false only", req: UpdateOpeningRequest{Remote: &f}},
		{name: "salary only", req: UpdateOpeningRequest{Salary: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErr(t, tt.req.Validate(), tt.wantErr)
		})
	}
}

func assertErr(t *testing.T, err error, want string) {
	t.Helper()
	if want == "" {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	if err == nil || err.Error() != want {
		t.Fatalf("err = %v, want %q", err, want)
	}
}
