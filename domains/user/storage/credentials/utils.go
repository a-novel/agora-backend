package credentials_storage

import (
	"fmt"
	"github.com/a-novel/agora-backend/models"
)

// WhereEmail returns arguments for a bun Where clause, to search for a precise email value.
//
//	db.NewSelect().Model(model).Where(WhereEmail("email", email))
func WhereEmail(source string, value models.Email) (string, string, string) {
	return fmt.Sprintf("%[1]s_user = ? AND %[1]s_domain = ?", source), value.User, value.Domain
}
