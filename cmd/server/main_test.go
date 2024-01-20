package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
func TestMetricUpdateHandler(t *testing.T) {
	type want struct {
		code          int
		request       string
		response      string
		expectStorage *storage.MemoryStorage
	}

	testCases := []struct {
		name    string
		storage *storage.MemoryStorage
		want    want
	}{
		{
			name:    "empty type",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{}),
			want: want{
				code:     http.StatusBadRequest,
				request:  "/update/qwerty/qwerty/1451",
				response: "Incorrect metric type\\n",
			},
		},
		{
			name:    "empty name",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{}),
			want: want{
				code:     http.StatusNotFound,
				request:  "/update/gauge/",
				response: "404 page not found\n",
			},
		},
		{
			name:    "incorrect value",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{}),
			want: want{
				code:     http.StatusBadRequest,
				request:  "/update/gauge/test1/qwerty",
				response: "Incorrect value\\n",
			},
		},
		{
			name: "add exist gauge metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					MType: entity.Gauge,
					ID:    "test1",
					Value: float64Ptr(2.5),
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/gauge/test1/2",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						MType: entity.Gauge,
						ID:    "test1",
						Value: float64Ptr(2.0),
					},
				}),
			},
		},
		{
			name: "add not exist gauge metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					MType: entity.Gauge,
					ID:    "test1",
					Value: float64Ptr(2.5),
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/gauge/test2/2",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						MType: entity.Gauge,
						ID:    "test1",
						Value: float64Ptr(2.5),
					},
					"test2": {
						MType: entity.Gauge,
						ID:    "test2",
						Value: float64Ptr(2.0),
					},
				}),
			},
		},
		{
			name: "add not exist counter metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					MType: entity.Counter,
					ID:    "test1",
					Delta: int64Ptr(2),
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/counter/test2/3",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						MType: entity.Counter,
						ID:    "test1",
						Delta: int64Ptr(2),
					},
					"test2": {
						MType: entity.Counter,
						ID:    "test2",
						Delta: int64Ptr(3),
					},
				}),
			},
		},
		{
			name: "add exist counter metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					MType: entity.Counter,
					ID:    "test1",
					Delta: int64Ptr(2),
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/counter/test1/3",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						MType: entity.Counter,
						ID:    "test1",
						Delta: int64Ptr(5),
					},
				}),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", handler.MetricUpdateHandler(testCase.storage))
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = srv.URL + testCase.want.request

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, testCase.want.code, resp.StatusCode())

			if testCase.want.response != "" {
				require.Equal(t, testCase.want.response, string(resp.Body()))
			}
			if testCase.want.expectStorage != nil {
				require.Equal(t, testCase.want.expectStorage, testCase.storage)
			}
		})
	}
}
