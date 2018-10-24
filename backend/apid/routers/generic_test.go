package routers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"

	"github.com/sensu/sensu-go/internal/apis/rbac"
	"github.com/sensu/sensu-go/testing/mockstore/v2"
)

// "clusterrole" in this function should really be "clusterroles" (note the
// plural), but we're waiting for plural support from the type registry.
func TestGenericListGlobal(t *testing.T) {
	store := &mockstore.MockStore{}

	store.On("List",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("string"),
		mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			clusterRoles := args.Get(2).(*[]rbac.ClusterRole)
			*clusterRoles = append(*clusterRoles, rbac.ClusterRole{
				Rules: []rbac.Rule{
					rbac.Rule{Verbs: []string{"*"}},
				},
			})
		})

	parentRouter := mux.NewRouter()
	genericRouter := NewGenericRouter(store)
	genericRouter.Mount(parentRouter)

	req, err := http.NewRequest("GET", "/apis/rbac/clusterroles", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	w := httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
}
