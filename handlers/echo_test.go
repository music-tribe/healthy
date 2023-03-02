package handlers

import (
	errs "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/errors"
	"github.com/music-tribe/healthy"
)

func TestHandler(t *testing.T) {

	e := echo.New()

	t.Run("When the services are healthy it should return a 200.", func(t *testing.T) {
		wantErr := false
		wantStatusCode := 200

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		hSvc := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", &healthy.MockPinger{}))
		h := Handler(hSvc)

		err := h(ctx)
		if (err != nil) != wantErr {
			t.Errorf("wanted error to be %v but got %v\n", wantErr, err)
			return
		}

		gotStatusCode := rec.Result().StatusCode
		if err != nil {
			ce := new(errors.CloudError)
			if errs.As(err, &ce) {
				gotStatusCode = ce.StatusCode
			}
		}

		if wantStatusCode != gotStatusCode {
			t.Errorf("wanted status code to be %d but got %d\n", wantStatusCode, gotStatusCode)
		}
	})

	t.Run("When the services are unhealthy it should return an error.", func(t *testing.T) {
		wantErr := false
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		erroringPinger := &healthy.MockPinger{
			Err: errors.NewCloudError(404, ""), // no matter what code we pass, it will return a 503 on error
		}

		hSvc := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", erroringPinger))
		h := Handler(hSvc)
		err := h(ctx)
		if (err != nil) != wantErr {
			t.Errorf("wanted error to be %v but got %v\n", wantErr, err)
			return
		}

		gotStatusCode := rec.Result().StatusCode
		if err != nil {
			ce := new(errors.CloudError)
			if errs.As(err, &ce) {
				gotStatusCode = ce.StatusCode
			}
		}

		if wantStatusCode != gotStatusCode {
			t.Errorf("wanted status code to be %d but got %d\n", wantStatusCode, gotStatusCode)
		}
	})

	t.Run("When we don't pass a pinger to our checker, it should return an error.", func(t *testing.T) {
		wantErr := false
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		hSvc := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", nil))
		h := Handler(hSvc)
		err := h(ctx)
		if (err != nil) != wantErr {
			t.Errorf("wanted error to be %v but got %v\n", wantErr, err)
			return
		}

		gotStatusCode := rec.Result().StatusCode
		if err != nil {
			ce := new(errors.CloudError)
			if errs.As(err, &ce) {
				gotStatusCode = ce.StatusCode
			}
		}

		if wantStatusCode != gotStatusCode {
			t.Errorf("wanted status code to be %d but got %d\n", wantStatusCode, gotStatusCode)
		}
	})
}
