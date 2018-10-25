package routers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sensu/sensu-go/internal/apis/rbac"
	storev2 "github.com/sensu/sensu-go/storage"
	"github.com/sensu/sensu-go/testing/mockstore/v2"
	"github.com/stretchr/testify/mock"
)

func TestRBACGet(t *testing.T) {
	type storeFunc func(*mockstore.MockStore)
	tests := []struct {
		name      string
		url       string
		storeFunc storeFunc
		wantCode  int
		wantLen   int
	}{
		{
			name:     "kind not found in the registry",
			url:      "/apis/rbac/missingkind/admin",
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "store error",
			url:  "/apis/rbac/clusterroles/admin",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(errors.New("error"))
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "empty response from store",
			url:  "/apis/rbac/clusterroles/admin",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(storev2.ErrNotFound)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "successful retrieval of the resource",
			url:  "/apis/rbac/clusterroles/admin",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(nil).
					Run(func(args mock.Arguments) {
						clusterRoles := args.Get(2).(*rbac.ClusterRole)
						*clusterRoles = rbac.ClusterRole{
							Rules: []rbac.Rule{
								rbac.Rule{Verbs: []string{"*"}},
							},
						}
					})
			},
			wantCode: http.StatusOK,
			wantLen:  122,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockstore.MockStore{}
			if tt.storeFunc != nil {
				tt.storeFunc(store)
			}
			parentRouter := mux.NewRouter()
			genericRouter := NewRBACRouter(store)
			genericRouter.Mount(parentRouter)

			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("http.NewRequest() error = %v", err)
			}

			w := httptest.NewRecorder()
			parentRouter.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("RBAC.list() = %d %s, want %v %s",
					w.Code, http.StatusText(w.Code), tt.wantCode, http.StatusText(tt.wantCode))
			}
			if w.Body.Len() < tt.wantLen {
				t.Errorf("RBAC.list() ContentLength = %d bytes, want %d bytes", w.Body.Len(), tt.wantLen)
			}
		})
	}
}

func TestRBACList(t *testing.T) {
	type storeFunc func(*mockstore.MockStore)
	tests := []struct {
		name      string
		url       string
		storeFunc storeFunc
		wantCode  int
		wantLen   int
	}{
		{
			name:     "kind not found in the registry",
			url:      "/apis/rbac/missingkind",
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "store error",
			url:  "/apis/rbac/clusterroles",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(errors.New("error"))
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "empty response from store",
			url:  "/apis/rbac/clusterroles",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(nil)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "successful retrieval of all resources",
			url:  "/apis/rbac/clusterroles",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(nil).
					Run(func(args mock.Arguments) {
						clusterRoles := args.Get(2).(*[]rbac.ClusterRole)
						*clusterRoles = append(*clusterRoles, rbac.ClusterRole{
							Rules: []rbac.Rule{
								rbac.Rule{Verbs: []string{"*"}},
							},
						})
					})
			},
			wantCode: http.StatusOK,
			wantLen:  124,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockstore.MockStore{}
			if tt.storeFunc != nil {
				tt.storeFunc(store)
			}
			parentRouter := mux.NewRouter()
			genericRouter := NewRBACRouter(store)
			genericRouter.Mount(parentRouter)

			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("http.NewRequest() error = %v", err)
			}

			w := httptest.NewRecorder()
			parentRouter.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("RBAC.list() = %d %s, want %v %s",
					w.Code, http.StatusText(w.Code), tt.wantCode, http.StatusText(tt.wantCode))
			}
			if w.Body.Len() < tt.wantLen {
				t.Errorf("RBAC.list() ContentLength = %d bytes, want %d bytes", w.Body.Len(), tt.wantLen)
			}
		})
	}
}
