package routers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sensu/sensu-go/backend/apid/actions"
	storev2 "github.com/sensu/sensu-go/storage"
)

// Namespaced is a router that handles requests to namespaced resources.
type Namespaced struct {
	store storev2.Store
}

// NewNamespacedRouter instantiates a new Namespaced.
func NewNamespacedRouter(store storev2.Store) *Namespaced {
	return &Namespaced{
		store: store,
	}
}

// Mount mounts this router's routes onto the given parent router.
func (r *Namespaced) Mount(parent *mux.Router) {
	routes := ResourceRoute{Router: parent, PathPrefix: "/apis/"}

	// For now we have to ignore the version variable because our apis are not
	// frozen yet
	//routes.Path("/{group}/{version}/{resource}", r.list).Methods(http.MethodGet)
	routes.Path("/{group}/namespaces/{namespace}/{resource}", r.list).Methods(http.MethodGet)
	routes.Path("/{group}/namespaces/{namespace}/{resource}/{name}", r.get).Methods(http.MethodGet)
}

func (r *Namespaced) get(req *http.Request) (interface{}, error) {
	return nil, actions.NewError(actions.InternalErr, errors.New("not implemented"))
}

func (r *Namespaced) list(req *http.Request) (interface{}, error) {
	return nil, actions.NewError(actions.InternalErr, errors.New("not implemented"))
}
