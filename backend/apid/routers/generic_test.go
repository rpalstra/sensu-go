package routers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sensu/sensu-go/testing/mockstore/v2"
)

func TestGenericListGlobal(t *testing.T) {
	store := &mockstore.MockStore{}

	parentRouter := mux.NewRouter()
	genericRouter := NewGenericRouter(store)
	genericRouter.Mount(parentRouter)

	req, err := http.NewRequest("GET", "/apis/rbac/clusterroles", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	w := httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)

	fmt.Println(w.Code)
}
