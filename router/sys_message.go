package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/piupuer/go-helper/api/v1"
)

func (rt Router) Message() gin.IRoutes {
	router1 := rt.Casbin("/message")
	router2 := rt.CasbinAndIdempotence("/message")
	router1.GET("/all", v1.FindMessage(rt.ops.v1Ops...))
	router1.GET("/unRead/count", v1.GetUnReadMessageCount(rt.ops.v1Ops...))
	router2.POST("/push", v1.PushMessage(rt.ops.v1Ops...))
	router1.PATCH("/read/batch", v1.BatchUpdateMessageRead(rt.ops.v1Ops...))
	router1.PATCH("/deleted/batch", v1.BatchUpdateMessageDeleted(rt.ops.v1Ops...))
	router1.PATCH("/read/all", v1.UpdateAllMessageRead(rt.ops.v1Ops...))
	router1.PATCH("/deleted/all", v1.UpdateAllMessageDeleted(rt.ops.v1Ops...))
	return rt.ops.group
}
