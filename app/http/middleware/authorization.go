package middleware

import (
	connection "booking/connection/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthorizationMiddleware struct {
	DB    connection.DBConnection
	Roles *[]Role
}

type Role struct {
	ID            string `gorm:"column:id" json:"id"`
	Name          string `gorm:"column:name" json:"name"`
	Privileges    string `gorm:"column:privileges" json:"privileges"`
	PrivilegesMap []Privilege
}

type Privilege struct {
	Resource     string   `json:"resource"`
	Descripttion string   `json:"description"`
	Scopes       []string `json:"scopes"`
}

func NewAuthorizationMiddleware(db connection.DBConnection) AuthorizationMiddleware {

	var roles []Role
	err := db.DB.Table("roles").Find(&roles).Error
	if err != nil {
		panic(err)
	}

	for i, each := range roles {
		var prvs []Privilege

		err := json.Unmarshal([]byte(each.Privileges), &prvs)

		if err != nil {
			panic(err)
		}

		roles[i].PrivilegesMap = prvs

	}

	return AuthorizationMiddleware{
		DB:    db,
		Roles: &roles,
	}
}

func (m AuthorizationMiddleware) Authorize(accessIdentifier string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*JWTCustomClaims)

			accessIdentifierMap := strings.Split(accessIdentifier, ".")
			resource := accessIdentifierMap[0]
			access := accessIdentifierMap[1]

			for _, each := range *m.Roles {
				fmt.Println(each)
				if claims.RoleID == each.ID {
					for _, elem := range each.PrivilegesMap {
						if elem.Resource == resource {
							for _, scp := range elem.Scopes {
								if access == scp {
									return next(c)
								}
							}
						}
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
		}
	}
}
