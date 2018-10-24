package routers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"

	"github.com/sensu/sensu-go/testing/mockstore/v2"
)

// "clusterrole" in this function should really be "clusterroles" (note the
// plural), but we're waiting for plural support from the type registry.
func TestGenericListGlobal(t *testing.T) {
	store := &mockstore.MockStore{}

	// I'm not sure why "context.Context" doesn't work for the type of the first
	// argument to List(). A type that works is "*context.valueCtx", infered
	// from the error I got, but that doesn't make sense to me and I must be
	// missing something here...
	// store.On("List", mock.AnythingOfType("context.Context"), mock.AnythingOfType("string"), mock.Anything).Return(nil)
	store.On("List", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.Anything).Return(nil)

	parentRouter := mux.NewRouter()
	genericRouter := NewGenericRouter(store)
	genericRouter.Mount(parentRouter)

	req, err := http.NewRequest("GET", "/apis/rbac/clusterrole", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	w := httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)

	fmt.Println(w.Result())
}
