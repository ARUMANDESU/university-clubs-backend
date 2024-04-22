package club

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
	"net/http"
	"reflect"
)

func (h *Handler) CreateRoleHandler(c *gin.Context) {
	const op = "ClubHandler.CreateRoleHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		Name  string `json:"name"`
		Color int32  `json:"color"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.clubClient.CreateRole(c, &clubv1.CreateRoleRequest{
		ClubId: clubID,
		UserId: userID,
		Name:   input.Name,
		Color:  input.Color,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": domain.ProtoToRole(role)})

}

func (h *Handler) DeleteRoleHandler(c *gin.Context) {
	const op = "ClubHandler.DeleteRoleHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roleID, err := utils.GetIntFromParams(c.Params, "role_id")
	if err != nil {
		log.Warn("failed to get role_id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	_, err = h.clubClient.DeleteRole(c, &clubv1.DeleteRoleRequest{
		ClubId: clubID,
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateRoleHandler(c *gin.Context) {
	const op = "ClubHandler.UpdateRoleHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roleID, err := utils.GetIntFromParams(c.Params, "role_id")
	if err != nil {
		log.Warn("failed to get role_id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		Name        string  `json:"name,omitempty"`
		Color       *int32  `json:"color,omitempty"`
		Permissions *uint64 `json:"permissions,omitempty"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// name -> change name
	// name, permissions -> name, permissions

	var paths []string
	if input.Name != "" {
		paths = append(paths, "name")
	}
	// TODO: fix this  does not update if frontend sends 0
	if !reflect.ValueOf(input.Color).IsZero() { // is not default
		paths = append(paths, "color")
	}
	if input.Permissions != nil {
		paths = append(paths, "permissions")
	}

	role, err := h.clubClient.UpdateRole(c, &clubv1.UpdateRoleRequest{
		ClubId:      clubID,
		UserId:      userID,
		RoleId:      roleID,
		Name:        input.Name,
		Permissions: *input.Permissions,
		Color:       *input.Color,
		UpdateMask:  &fieldmaskpb.FieldMask{Paths: paths},
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": domain.ProtoToRole(role)})
}

func (h *Handler) UpdateRolesPositionHandler(c *gin.Context) {
	const op = "ClubHandler.UpdateRolesPositionHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		Roles []*clubv1.ChangeRolesPositionItems `json:"roles"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.clubClient.ChangeRolesPosition(c, &clubv1.ChangeRolesPositionRequest{
		ClubId: clubID,
		UserId: userID,
		Roles:  input.Roles,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": domain.MapProtoToRoleArr(result.Roles)})
}

func (h *Handler) AddRoleMembersHandler(c *gin.Context) {
	const op = "ClubHandler.AddRoleMembersHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roleID, err := utils.GetIntFromParams(c.Params, "role_id")
	if err != nil {
		log.Warn("failed to get role_id from params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		Members []int64 `json:"members"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.clubClient.AddRoleMembers(c, &clubv1.AddRoleMembersRequest{
		ClubId:  clubID,
		RoleId:  roleID,
		UserId:  userID,
		UsersId: input.Members,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
