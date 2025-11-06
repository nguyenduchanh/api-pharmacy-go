package enum

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

func (m HttpMethod) IsValid() bool {
	switch m {
	case GET, POST, PUT, DELETE:
		return true
	}
	return false
}

var HttpMethodDescriptions = map[HttpMethod]string{
	GET:    "Get",
	POST:   "Post",
	PUT:    "Put",
	DELETE: "Delete",
}
