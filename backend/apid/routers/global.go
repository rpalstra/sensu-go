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

// Global is a router that handles requests to globally defined resources.
type Global struct {
	store storev2.Store
}

// NewGlobalRouter instantiates a new Global.
func NewGlobalRouter(store storev2.Store) *Global {
	return &Global{
		store: store,
	}
}

// Mount mounts this router's routes onto the given parent router.
func (r *Global) Mount(parent *mux.Router) {
	routes := ResourceRoute{Router: parent, PathPrefix: "/apis/"}

	// For now we have to ignore the version variable because our apis are not
	// frozen yet
	//routes.Path("/{group}/{version}/{resource}", r.listGlobal).Methods(http.MethodGet)
	routes.Path("/{group}/{resource}", r.list).Methods(http.MethodGet)
	routes.Path("/{group}/{resource}/{name}", r.get).Methods(http.MethodGet)
}

func (r *Global) get(req *http.Request) (interface{}, error) {
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

	// Initialize a new struct of type typeOfGvk and fill all the required meta
	// fields required to obtain the key
	objectMeta := meta.ObjectMeta{Name: name}
	typeMeta := meta.TypeMeta{
		APIVersion: apiVersion,
		Kind:       kind,
	}
	keyablePtr := reflect.New(typeOfGvk)
	keyable := reflect.Indirect(keyablePtr)
	keyable.FieldByName("ObjectMeta").Set(reflect.ValueOf(objectMeta))
	keyable.FieldByName("TypeMeta").Set(reflect.ValueOf(typeMeta))
	storeKey := storev2.StorageKey{Keyable: keyablePtr.Interface().(storev2.Keyable)}

	// Create a pointer of an unitialized typeOfGvk
	ptr := reflect.New(typeOfGvk).Interface()

	// Get the requested object from the store
	err = r.store.Get(req.Context(), storeKey.Path(), ptr)
	if err != nil {
		if err == storev2.ErrNotFound {
			return nil, actions.NewErrorf(actions.NotFound, "no resource found")
		}
	}

	return ptr, err
}

func (r *Global) list(req *http.Request) (interface{}, error) {
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

	// Initialize a new struct of type typeOfGvk and fill all the required meta
	// fields required to obtain the key
	typeMeta := meta.TypeMeta{
		APIVersion: apiVersion,
		Kind:       kind,
	}
	keyablePtr := reflect.New(typeOfGvk)
	keyable := reflect.Indirect(keyablePtr)
	keyable.FieldByName("TypeMeta").Set(reflect.ValueOf(typeMeta))
	storeKey := storev2.StorageKey{Keyable: keyablePtr.Interface().(storev2.Keyable)}

	// Get the slice type with element of type typeOfGvk
	sliceOfGvk := reflect.SliceOf(typeOfGvk)

	// Create a pointer of an unitialized sliceOfGvk
	ptr := reflect.New(sliceOfGvk)

	// Set an empty slice of type sliceOfGvk to our pointer
	ptr.Elem().Set(reflect.MakeSlice(sliceOfGvk, 0, 0))

	// Assign our empty slice pointer to slicePtr so we can pass it to the store
	slicePtr := ptr.Interface()

	// List all elements of the given kind
	err = r.store.List(req.Context(), storeKey.PrefixPath(), slicePtr)
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
