package helpers

import (
    "math"
    "strconv"
    "github.com/gin-gonic/gin"
)

type Pagination struct {
    TotalRecords int64
    PageSize     int
    CurrentPage  int
    TotalPages   int
}

func NewPagination(c *gin.Context, totalRecords int64) (*Pagination, error) {
    pageStr := c.DefaultQuery("page", "1")
    pageSizeStr := c.DefaultQuery("pageSize", "10")

    page, err := strconv.Atoi(pageStr)
    if err != nil {
        return nil, err
    }

    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil {
        return nil, err
    }

    totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

    return &Pagination{
        TotalRecords: totalRecords,
        PageSize:     pageSize,
        CurrentPage:  page,
        TotalPages:   totalPages,
    }, nil
}