package mapper

import (
	appuser "github.com/mvrduk/statum/services/user-service/internal/application/service/user"
	"github.com/mvrduk/statum/services/user-service/internal/transport/http/gen"
)

func ToUserResponse(u appuser.UserDTO) gen.UserResponse {
	return gen.UserResponse{Id: u.ID, Email: u.Email, CreatedAt: u.CreatedAt}
}
