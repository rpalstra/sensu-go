package routers

import (
	"net/http"
	"path"
	"reflect"
	"strings"

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

// Mount ...
func (r *Generic) Mount(parent *mux.Router) {
	routes := ResourceRoute{Router: parent, PathPrefix: "/apis/"}

	// Global resources
	// For now we have to ignore the version variable because our apis are not
	// frozen yet
	//routes.Path("/{group}/{version}/{resource}", r.listGlobal).Methods(http.MethodGet)
	routes.Path("/{group}/{resource}", r.listGlobal).Methods(http.MethodGet)
}

func (r *Generic) listGlobal(req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	resource := strings.ToLower(vars["resource"])
	apiVersion := path.Join(vars["group"], vars["version"])

	// HACK: remove the last letter of the resource to get the kind
	// TODO: remove when https://github.com/sensu/sensu-go/pull/2214 is merged
	kind := strings.TrimSuffix(resource, "s")

	gvk, err := registry.Resolve(meta.TypeMeta{
		Kind:       kind,
		APIVersion: apiVersion,
	})
	if err != nil {
		return nil, err
	}

	// t is the concrete type of the meta.GroupVersionKind returned by
	// registry.Resolve()
	t := reflect.TypeOf(gvk)
	slice := reflect.SliceOf(t)
	ptr := reflect.New(slice)
	ptr.Elem().Set(reflect.MakeSlice(slice, 0, 0))
	s := ptr.Interface()
	err = r.store.List(req.Context(), kind, s)

	return s, err
}
