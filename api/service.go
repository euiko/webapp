package api

type (
	Service interface {
		Route(router Router)
	}

	APIService interface {
		APIRoute(router Router)
	}
)
