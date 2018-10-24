package routers

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/sensu/sensu-go/internal/apis/meta"
	"github.com/sensu/sensu-go/runtime/registry"
	storev2 "github.com/sensu/sensu-go/storage"
)

// Generic is a router that handles requests for any API group resources.
type Generic struct {
	store storev2.Store
}

// NewGenericRouter instantiates a new router that can handle any versioned,
// namespaced, API group routes.
func NewGenericRouter(store storev2.Store) *Generic {
	return &Generic{
		store: store,
	}
}

func (r *Generic) Mount(parent *mux.Router) {
	routes := ResourceRoute{Router: parent, PathPrefix: "/apis/"}
	routes.Path("/{group}/{version}/{kind}", r.list).Methods(http.MethodGet)
}

func (r *Generic) list(req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)

	kind := vars["kind"]
	apiVer := vars["group"] + "/" + vars["version"]

	// This probably won't work because Resolve() expects a "capitalized" kind.
	// I'm also not 100% sure of the format apiVer should have.
	v, err := registry.Resolve(meta.TypeMeta{Kind: kind, APIVersion: apiVer})
	slice := reflect.New(reflect.SliceOf(v)).Interface()

	err = r.store.List(req.Context(), vars["kind"], slice)
	return slice, err
}
