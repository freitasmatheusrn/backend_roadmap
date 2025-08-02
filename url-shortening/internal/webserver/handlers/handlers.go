package handlers

type Handlers struct {
	Url *UrlHandler
}

func NewHandler(url *UrlHandler) *Handlers {
	return &Handlers{
		Url: url,
	}
}
