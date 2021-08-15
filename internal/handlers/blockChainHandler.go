package handlers

import (
	"net/http"

	"github.com/CharlesChou03/_git/block-chain.git/models"
	modelsReq "github.com/CharlesChou03/_git/block-chain.git/models/request"
	modelsRes "github.com/CharlesChou03/_git/block-chain.git/models/response"
	"github.com/CharlesChou03/_git/block-chain.git/services"
	"github.com/gin-gonic/gin"
)

// GetLatestNBlockHandler get latest n block information
// @Summary get latest n block information
// @Description get latest n block information
// @Tags Block Chain Information
// @Accept json
// @Produce json
// @Param limit query int false "limit"
// @Success 200 {object} modelsRes.GetLatestNBlockRes "ok"
// @Failure 204 "no content"
// @Failure 400 "bad request"
// @Failure 424 "failed dependency"
// @Failure 500 "internal error"
// @Router /blocks [get]
func GetLatestNBlockHandler(c *gin.Context) {
	req := modelsReq.GetLatestNBlockReq{}
	res := modelsRes.GetLatestNBlockRes{}
	c.Bind(&req)
	statusCode, err := services.GetBlocks(&req, &res)
	switch statusCode {
	case 200:
		c.JSON(http.StatusOK, res)
	case 204:
		c.JSON(http.StatusNoContent, err)
	case 400:
		c.JSON(http.StatusBadRequest, err)
	case 424:
		c.JSON(http.StatusFailedDependency, err)
	default:
		c.JSON(http.StatusInternalServerError, models.InternalServerError)
	}
}

// GetBlockByIdHandler get block information by block number
// @Summary get block information by block number
// @Description get block information by block number
// @Tags Block Chain Information
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} modelsRes.GetBlockByNumRes "ok"
// @Failure 204 "no content"
// @Failure 400 "bad request"
// @Failure 500 "internal error"
// @Router /blocks/{id} [get]
func GetBlockByNumHandler(c *gin.Context) {
	req := modelsReq.GetBlockByNumReq{}
	res := modelsRes.GetBlockByNumRes{}
	c.BindUri(&req)
	statusCode, err := services.GetBlock(&req, &res)
	switch statusCode {
	case 200:
		c.JSON(http.StatusOK, res)
	case 204:
		c.JSON(http.StatusNoContent, err)
	case 400:
		c.JSON(http.StatusBadRequest, err)
	default:
		c.JSON(http.StatusInternalServerError, models.InternalServerError)
	}
}

// GetTransactionByHashHandler get transaction information by hash
// @Summary get transaction information by hash
// @Description get transaction information by hash
// @Tags Block Chain Information
// @Accept json
// @Produce json
// @Param txHash path string true "txHash"
// @Success 200 {object} modelsRes.GetTransactionByHashRes "ok"
// @Failure 204 "no content"
// @Failure 400 "bad request"
// @Failure 500 "internal error"
// @Router /transaction/{txHash} [get]
func GetTransactionByHashHandler(c *gin.Context) {
	req := modelsReq.GetTransactionByHashReq{}
	res := modelsRes.GetTransactionByHashRes{}
	c.BindUri(&req)
	statusCode, err := services.GetTransaction(&req, &res)
	switch statusCode {
	case 200:
		c.JSON(http.StatusOK, res)
	case 204:
		c.JSON(http.StatusNoContent, err)
	case 400:
		c.JSON(http.StatusBadRequest, err)
	default:
		c.JSON(http.StatusInternalServerError, models.InternalServerError)
	}
}
