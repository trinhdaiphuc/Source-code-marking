package classes

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *ClassHandler) ClassStatistic(c echo.Context) (err error) {
	statistic := &models.StatisticQueryParam{}
	if err := c.Bind(statistic); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	begin, _ := time.Parse(time.RFC3339, statistic.BeginDate)
	end, _ := time.Parse(time.RFC3339, statistic.EndDate)

	h.Logger.Debug("StatisticQueryParam: begin: ", begin, ", end: ", end)

	match := bson.D{
		primitive.E{Key: "created_at", Value: bson.D{
			primitive.E{Key: "$gte", Value: begin},
			primitive.E{Key: "$lte", Value: end},
		}},
	}

	group := bson.D{
		primitive.E{Key: "_id", Value: bson.D{
			primitive.E{Key: "$dateToString", Value: bson.D{
				primitive.E{Key: "format", Value: "%d-%m-%Y"},
				primitive.E{Key: "date", Value: "$created_at"},
			}},
		}},
		primitive.E{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}},
	}

	matchStage := bson.D{primitive.E{Key: "$match", Value: match}}
	groupStage := bson.D{primitive.E{Key: "$group", Value: group}}

	classCollection := models.GetClassCollection(h.DB)
	ctx := context.TODO()
	totalInfoCursor, err := classCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		h.Logger.Error("Error when group ", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal server error",
			Internal: err,
		}
	}

	totalWithInfo := []models.ListUserStatistic{}
	if err = totalInfoCursor.All(ctx, &totalWithInfo); err != nil {
		h.Logger.Error("Error when decode cursor ", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal server error",
			Internal: err,
		}
	}

	return c.JSON(http.StatusOK, totalWithInfo)
}
