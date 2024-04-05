package handlers

import (
	errs "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/errors"
	"github.com/music-tribe/healthy"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {

	e := echo.New()

	t.Run("When the services are healthy it should return a 200.", func(t *testing.T) {
		wantStatusCode := 200

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		hSvc, _ := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", &healthy.MockPinger{}))
		h := Handler(hSvc)

		err := h(ctx)
		require.NoError(t, err)

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

	t.Run("When the services are unhealthy it should return a 503.", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		erroringPinger := &healthy.MockPinger{
			Err: errors.NewCloudError(404, ""), // no matter what code we pass, it will return a 503 on error
		}

		hSvc, err := healthy.New("some-service", "v1.1.3", healthy.NewChecker("mockStore", erroringPinger))
		require.NoError(t, err)

		h := Handler(hSvc)
		err = h(ctx)
		require.NoError(t, err)

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

	t.Run("When one of many services are unhealthy it should return a 503.", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		erroringPinger := &healthy.MockPinger{
			Err: errors.NewCloudError(404, ""), // no matter what code we pass, it will return a 503 on error
		}

		hSvc, err := healthy.New(
			"some-service",
			"v1.1.3",
			healthy.NewChecker("mockStoreError", erroringPinger),
			healthy.NewChecker("mockStoreHealthy", &healthy.MockPinger{}),
		)
		require.NoError(t, err)

		h := Handler(hSvc)
		err = h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()

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

	t.Run("When one of many services, added in a different order, are unhealthy it should return a 503.", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		erroringPinger := &healthy.MockPinger{
			Err: errors.NewCloudError(404, ""), // no matter what code we pass, it will return a 503 on error
		}

		hSvc, err := healthy.New(
			"some-service",
			"v1.1.3",
			healthy.NewChecker("mockStoreHealthy", &healthy.MockPinger{}),
			healthy.NewChecker("mockStoreError", erroringPinger),
		)
		require.NoError(t, err)

		h := Handler(hSvc)
		err = h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()

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

	t.Run("When we don't pass a pinger to our checker, it should return a 503.", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		hSvc, _ := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", nil))
		h := Handler(hSvc)
		err := h(ctx)
		require.NoError(t, err)

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
