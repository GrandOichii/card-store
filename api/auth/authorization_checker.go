package auth

import (
	"github.com/gin-gonic/gin"
	"store.api/model"
)

type PermitFunc func(*model.User) bool

type AuthorizationChecker interface {
	Check(c *gin.Context, user *model.User) (authorized bool, matches bool)
}

type ChainAuthorizationChecker struct {
	path     string
	checkers []*MethodChecker
}

func (ch ChainAuthorizationChecker) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	if ch.path != c.Request.URL.Path {
		return false, false
	}

	authorized = false
	for _, checker := range ch.checkers {
		if checker.Method == "*" || checker.Method == c.Request.Method {
			authorized = checker.PermitF(user)
		}
	}

	return authorized, true
}

type AuthorizationCheckerBuilder struct {
	path     string
	checkers []*MethodChecker
	current  *MethodChecker
}

type MethodChecker struct {
	Method  string
	PermitF PermitFunc
}

func NewAuthorizationCheckerBuilder(path string) *AuthorizationCheckerBuilder {
	result := new(AuthorizationCheckerBuilder)
	result.path = path
	result.checkers = make([]*MethodChecker, 0)
	return result
}

func (c *AuthorizationCheckerBuilder) addPermit(f PermitFunc) {
	c.current.PermitF = f
	c.checkers = append(c.checkers, c.current)
	c.current = nil
}

func (c *AuthorizationCheckerBuilder) addMethod(method string) {
	c.current = new(MethodChecker)
	c.current.Method = method
}

func (c *AuthorizationCheckerBuilder) ForAnyMethod() *AuthorizationCheckerBuilder {
	c.addMethod("*")

	return c
}

func (c *AuthorizationCheckerBuilder) PermitAll() *AuthorizationCheckerBuilder {
	c.addPermit(func(u *model.User) bool {
		return true
	})

	return c
}

func (c *AuthorizationCheckerBuilder) ForMethod(method string) *AuthorizationCheckerBuilder {
	c.addMethod(method)

	return c
}

func (c *AuthorizationCheckerBuilder) Permit(f PermitFunc) *AuthorizationCheckerBuilder {
	c.addPermit(f)

	return c
}

// .ForAnyMethod()
// .PermitAll()
// .ForMethod("POST")
// .Permit(func(user *model.User) bool {
// 	return user.IsAdmin && user.Verified
// })

func (b *AuthorizationCheckerBuilder) Build() *ChainAuthorizationChecker {
	result := new(ChainAuthorizationChecker)
	result.path = b.path
	result.checkers = b.checkers

	return result
}
