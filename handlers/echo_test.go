package handlers

import (
	"context"
	"encoding/json"
	errs "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/errors"
	"github.com/music-tribe/healthy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {

	e := echo.New()

	t.Run("When the services are healthy it should return a 200.", func(t *testing.T) {
		wantStatusCode := 200

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		hSvc, _ := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", healthy.NewMockChecker(nil)))
		h := Handler(hSvc)

		err := h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("OK"))

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

		erroringChecker := healthy.NewMockChecker(errors.NewCloudError(404, "")) // no matter what code we pass, it will return a 503 on error

		hSvc, err := healthy.New("some-service", "v1.1.3", healthy.NewChecker("mockStore", erroringChecker))
		require.NoError(t, err)

		h := Handler(hSvc)
		err = h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("Unavailable"))

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

		erroringChecker := healthy.NewMockChecker(errors.NewCloudError(404, "")) // no matter what code we pass, it will return a 503 on error

		hSvc, err := healthy.New(
			"some-service",
			"v1.1.3",
			healthy.NewChecker("mockStoreError", erroringChecker),
			healthy.NewChecker("mockStoreHealthy", healthy.NewMockChecker(nil)),
		)
		require.NoError(t, err)

		h := Handler(hSvc)
		err = h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("Unavailable"))

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

		erroringChecker := healthy.NewMockChecker(errors.NewCloudError(404, "")) // no matter what code we pass, it will return a 503 on error

		hSvc, err := healthy.New(
			"some-service",
			"v1.1.3",
			healthy.NewChecker("mockStoreHealthy", healthy.NewMockChecker(nil)),
			healthy.NewChecker("mockStoreError", erroringChecker),
		)
		require.NoError(t, err)

		h := Handler(hSvc)
		err = h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("Unavailable"))

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

	t.Run("When we don't pass a checkFunc to our checker, it should return a 503.", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		hSvc, _ := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", nil))
		h := Handler(hSvc)
		err := h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("Unavailable"))

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

	t.Run("When we pass a checker with timeout we should return a StatusTimeout after the timeout expires", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		slowChecker := health.CheckFunc(func(ctx context.Context) error {
			time.Sleep(time.Second * 2)
			return nil
		})
		hSvc, _ := healthy.New("some-service", "1.1.3", healthy.NewCheckerWithTimeout("mockStore", slowChecker, 1*time.Second))
		h := Handler(hSvc)
		err := h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("Unavailable"))
		assert.Equal(t, healthCheck.Failures["mockStore"], "Timeout during health check")

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

	t.Run("When we pass a checker and use default timeout we should return a StatusTimeout after the default timeout of 2 seconds expires", func(t *testing.T) {
		wantStatusCode := 503

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		slowChecker := health.CheckFunc(func(ctx context.Context) error {
			time.Sleep(time.Second * 3)
			return nil
		})
		hSvc, _ := healthy.New("some-service", "1.1.3", healthy.NewChecker("mockStore", slowChecker))
		h := Handler(hSvc)
		err := h(ctx)
		require.NoError(t, err)

		defer rec.Result().Body.Close()
		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		healthCheck := health.Check{}
		err = json.Unmarshal(body, &healthCheck)
		require.NoError(t, err)
		assert.Equal(t, healthCheck.Status, health.Status("Unavailable"))
		assert.Equal(t, healthCheck.Failures["mockStore"], "Timeout during health check")

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
