package router

import "github.com/piupuer/go-helper/pkg/middleware"

func (rt Router) Base() {
	if rt.ops.jwt {
		router1 := rt.ops.Group.Group("/base")
		router2 := rt.Casbin("/base")
		router1.POST("/login", middleware.JwtLogin(rt.ops.jwtOps...))
		router1.POST("/logout", middleware.JwtLogout(rt.ops.jwtOps...))
		router1.POST("/refreshToken", middleware.JwtRefresh(rt.ops.jwtOps...))
		if rt.ops.idempotence {
			// need login
			router2.GET("/idempotenceToken", middleware.GetIdempotenceToken(rt.ops.idempotenceOps...))
		}
	}
}
