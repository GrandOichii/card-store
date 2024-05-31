package auth

import (
	"github.com/IGLOU-EU/go-wildcard"
	"github.com/gin-gonic/gin"
	"store.api/model"
)

type PermitFunc func(*model.User) bool

type AuthorizationChecker interface {
	Check(c *gin.Context, user *model.User) (authorized bool, matches bool)
}

type ChainAuthorizationChecker struct {
	checkers []*MethodChecker
}

func (ch ChainAuthorizationChecker) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	authorized = false
	matches = false
	for _, checker := range ch.checkers {
		if !wildcard.Match(checker.Path, c.Request.URL.Path) {
			continue
		}
		if checker.Method == "*" || checker.Method == c.Request.Method {
			authorized = checker.PermitF(user)
			matches = true
		}
	}

	return authorized, matches
}

type AuthorizationCheckerBuilder struct {
	checkers []*MethodChecker
	current  *MethodChecker
}

type MethodChecker struct {
	Method  string
	Path    string
	PermitF PermitFunc
}

func NewAuthorizationCheckerBuilder() *AuthorizationCheckerBuilder {
	result := new(AuthorizationCheckerBuilder)
	result.checkers = make([]*MethodChecker, 0)
	return result
}

func (c *AuthorizationCheckerBuilder) addPermit(f PermitFunc) {
	c.current.PermitF = f
	c.checkers = append(c.checkers, c.current)
	c.current = nil
}

func (c *AuthorizationCheckerBuilder) addMethod(method string) {
	if c.current == nil {
		c.current = new(MethodChecker)
	}
	c.current.Method = method
}

func (c *AuthorizationCheckerBuilder) addPath(path string) {
	if c.current == nil {
		c.current = new(MethodChecker)
	}
	c.current.Path = path
}

func (c *AuthorizationCheckerBuilder) ForPath(path string) *AuthorizationCheckerBuilder {
	c.addPath(path)

	return c
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

func (b *AuthorizationCheckerBuilder) Build() *ChainAuthorizationChecker {
	result := new(ChainAuthorizationChecker)
	result.checkers = b.checkers

	return result
}
