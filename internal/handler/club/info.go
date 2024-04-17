package club

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) GetClubHandler(c *gin.Context) {
	const op = "ClubHandler.GetClubHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clubClient.GetClub(c, &clubv1.GetClubRequest{ClubId: clubID})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"club": domain.ClubObjectToClub(res)})
}

func (h *Handler) ListClubsHandler(c *gin.Context) {
	const op = "ClubHandler.ListClubsHandler"
	log := h.log.With(slog.String("op", op))

	query := c.Query("query")
	//todo: make in another way
	clubTypeStr := c.Query("club_types")
	var clubType []string
	if clubTypeStr != "" {
		clubType = strings.Split(clubTypeStr, ",")
	}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clubClient.ListClubs(c, &clubv1.ListClubRequest{
		Query:      query,
		ClubType:   clubType,
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"clubs": domain.MapClubObjArrToClubArr(res.Clubs), "metadata": res.Metadata})
}

func (h *Handler) ListClubMembersHandler(c *gin.Context) {
	const op = "ClubHandler.NewClubHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clubClient.ListClubMembers(c, &clubv1.ListClubMembersRequest{
		ClubId:     clubID,
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": domain.MapUserObjArrToMemberArr(res.GetUsers()), "metadata": res.GetMetadata()})

}

func (h *Handler) ListJoinRequestsHandler(c *gin.Context) {
	const op = "ClubHandler.ListJoinRequestsHandler"
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

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clubClient.ListMembershipRequests(c, &clubv1.ListMembershipRequestsRequest{
		UserId:     userID,
		ClubId:     clubID,
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": domain.MapUserObjArrToMemberArr(res.GetUsers()), "metadata": res.GetMetadata()})
}

func (h *Handler) ListNewClubRequestsHandler(c *gin.Context) {
	const op = "ClubHandler.ListNewClubRequestsHandler"
	log := h.log.With(slog.String("op", op))

	query := c.Query("query")
	//todo: make in another way
	clubTypeStr := c.Query("club_types")
	var clubType []string
	if clubTypeStr != "" {
		clubType = strings.Split(clubTypeStr, ",")
	}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clubClient.ListNotApprovedClubs(c, &clubv1.ListNotApprovedClubsRequest{
		Query:      query,
		ClubType:   clubType,
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": domain.MapClubUserArrToClubList(res.GetList()), "metadata": res.Metadata})
}

func (h *Handler) GetUserClubsHandler(c *gin.Context) {
	const op = "ClubHandler.GetUserClubsHandler"
	log := h.log.With(slog.String("op", op))

	userID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clubClient.GetUserClubs(c, &clubv1.GetUserClubsRequest{UserId: userID})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"clubs": domain.MapClubObjArrToClubArr(res.Clubs)})

}
