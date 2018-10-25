package routers

import (
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sensu/sensu-go/backend/apid/actions"
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
	routes.Path("/{group}/{resource}/{name}", r.get).Methods(http.MethodGet)
}

func (r *RBAC) get(req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	group, version, resource, err := getGvrFromRoute(vars)
	if err != nil {
		return nil, actions.NewError(actions.InvalidArgument, err)
	}
	name, err := url.PathUnescape(vars["name"])
	if err != nil {
		return nil, actions.NewError(actions.InvalidArgument, err)
	}
	apiVersion := path.Join(group, version)

	// HACK: remove the last letter of the resource to get the kind
	// TODO: remove when https://github.com/sensu/sensu-go/pull/2214 is merged
	kind := strings.TrimSuffix(resource, "s")

	// TODO: use the key builder instead
	key := path.Join(kind, name)

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

	// Create a pointer of an unitialized sliceOfGvk
	ptr := reflect.New(typeOfGvk)

	// Assign our empty slice pointer to slicePtr so we can pass it to the store
	elemPtr := ptr.Interface()

	// List all elements of the given kind
	err = r.store.Get(req.Context(), key, elemPtr)
	if err != nil {
		if err == storev2.ErrNotFound {
			return nil, actions.NewErrorf(actions.NotFound, "no resource found")
		}
	}

	return elemPtr, err
}

func (r *RBAC) list(req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	group, version, resource, err := getGvrFromRoute(vars)
	if err != nil {
		return nil, actions.NewError(actions.InvalidArgument, err)
	}
	apiVersion := path.Join(group, version)

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
	if err != nil {
		return nil, err
	}

	// Make sure we retrieved at least one element
	// TODO: the store should return storage.ErrNotFound, just like the get method
	// does
	if reflect.ValueOf(slicePtr).Elem().Len() == 0 {
		return nil, actions.NewErrorf(actions.NotFound, "no resources found")
	}

	return slicePtr, nil
}
