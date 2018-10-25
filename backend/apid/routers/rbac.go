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

// RBAC is a router that handles requests for the RBAC API group resources.
type RBAC struct {
	store storev2.Store
}

// NewRBACRouter instantiates a router that handles requests to the RBAC group
func NewRBACRouter(store storev2.Store) *RBAC {
	return &RBAC{
		store: store,
	}
}

// Mount ...
func (r *RBAC) Mount(parent *mux.Router) {
	routes := ResourceRoute{Router: parent, PathPrefix: "/apis/"}

	// For now we have to ignore the version variable because our apis are not
	// frozen yet
	//routes.Path("/{group}/{version}/{resource}", r.listGlobal).Methods(http.MethodGet)
	routes.Path("/{group}/{resource}", r.list).Methods(http.MethodGet)
}

func (r *RBAC) list(req *http.Request) (interface{}, error) {
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

	// Get the concrete type of the meta.GroupVersionKind returned by
	// registry.Resolve()
	typeOfGvk := reflect.TypeOf(gvk)

	// Get the slice type with element of type typeOfGvk
	sliceOfGvk := reflect.SliceOf(typeOfGvk)

	// Create a pointer of an unitialized sliceOfGvk
	ptr := reflect.New(sliceOfGvk)

	// Set an empty slice of type sliceOfGvk to our pointer
	ptr.Elem().Set(reflect.MakeSlice(sliceOfGvk, 0, 0))

	// Assign our empty slice pointer to slicePtr so we can pass it to the store
	slicePtr := ptr.Interface()

	// List all elements of the given kind
	err = r.store.List(req.Context(), kind, slicePtr)

	return slicePtr, err
}
