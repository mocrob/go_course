package main

import (
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
				response: "Некорректный тип\n",
			},
		},
		{
			name:    "empty name",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{}),
			want: want{
				code:     http.StatusNotFound,
				request:  "/update/gauge/",
				response: "Некорректный URL\n",
			},
		},
		{
			name:    "incorrect value",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{}),
			want: want{
				code:     http.StatusBadRequest,
				request:  "/update/gauge/test1/qwerty",
				response: "Некорректное значение\n",
			},
		},
		{
			name: "add exist gauge metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					Type:  entity.Gauge,
					Name:  "test1",
					Value: 2.5,
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/gauge/test1/2",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						Type:  entity.Gauge,
						Name:  "test1",
						Value: 2.0,
					},
				}),
			},
		},
		{
			name: "add not exist gauge metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					Type:  entity.Gauge,
					Name:  "test1",
					Value: 2.5,
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/gauge/test2/2",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						Type:  entity.Gauge,
						Name:  "test1",
						Value: 2.5,
					},
					"test2": {
						Type:  entity.Gauge,
						Name:  "test2",
						Value: 2.0,
					},
				}),
			},
		},
		{
			name: "add not exist counter metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					Type:  entity.Counter,
					Name:  "test1",
					Value: int64(2),
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/counter/test2/3",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						Type:  entity.Counter,
						Name:  "test1",
						Value: int64(2),
					},
					"test2": {
						Type:  entity.Counter,
						Name:  "test2",
						Value: int64(3),
					},
				}),
			},
		},
		{
			name: "add exist counter metric",
			storage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
				"test1": {
					Type:  entity.Counter,
					Name:  "test1",
					Value: int64(2),
				},
			}),
			want: want{
				code:     http.StatusOK,
				request:  "/update/counter/test1/3",
				response: "",
				expectStorage: storage.NewMemoryStoragePlusMetrics(map[string]entity.Metric{
					"test1": {
						Type:  entity.Counter,
						Name:  "test1",
						Value: int64(5),
					},
				}),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, testCase.want.request, nil)
			w := httptest.NewRecorder()
			h := handler.MetricUpdateHandler(testCase.storage)
			h(w, request)

			res := w.Result()

			assert.Equal(t, testCase.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, testCase.want.response, string(resBody))

			if testCase.want.expectStorage != nil {
				require.Equal(t, testCase.want.expectStorage, testCase.storage)
			}
		})
	}
}
