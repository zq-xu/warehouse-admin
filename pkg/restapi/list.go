package restapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zq-xu/warehouse-admin/pkg/restapi/list"
	"zq-xu/warehouse-admin/pkg/restapi/response"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
)

var (
	defaultBaseInfoColumnList = []string{"id", "name"}
)

type ListConf struct {
	//ModelObj struct should contain the TableName function or by obtaining the model struct
	//For example:
	//	type Supplier struct {
	//	  Model
	//  }
	//
	//  func (sr *Supplier) TableName() string {
	//	  return SupplierTableName
	//  }
	//
	//  type SupplierDetail struct {
	//    Supplier
	//    ......
	//  }
	//
	//ModelObj: &SupplierDetail{}
	AuthControl

	ModelObj              interface{}
	ModelObjList          interface{}
	TransObjToRespFunc    func(ac *auth.AccessControl) []interface{}
	FuzzySearchColumnList []string
	BaseInfoColumnList    []string

	LoadAssociationsDBFunc func(db, queryDB *gorm.DB, ac *auth.AccessControl) *gorm.DB
	GenerateQueryFunc      func(db *gorm.DB, reqParams *list.Params) *gorm.DB
	ResponseWriteFunc      func(ctx *gin.Context)
}

var ApiListInstance = &apiList{}

type apiList struct{}

func (l *apiList) List(ctx *gin.Context, conf *ListConf) {
	ei := conf.AuthControl.Validate(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	reqParams, ei := list.GetListParams(ctx)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	resp, ei := l.list(ctx, reqParams, conf)
	if ei != nil {
		ctx.JSON(ei.Status, ei)
		return
	}

	if conf.ResponseWriteFunc != nil {
		conf.ResponseWriteFunc(ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (l *apiList) list(ctx context.Context, reqParams *list.Params, conf *ListConf) (*list.PageResponse, *response.ErrorInfo) {
	db := store.DB(ctx)

	loadDB := store.GenerateDBForQuery(db, conf.FuzzySearchColumnList, reqParams.Queries[list.QueryValue])
	if conf.GenerateQueryFunc != nil {
		loadDB = conf.GenerateQueryFunc(loadDB, reqParams)
	}

	count, err := store.GetCount(loadDB, conf.ModelObj)
	if err != nil {
		return nil, response.NewStorageError(response.StorageErrorCode, fmt.Sprintf("get count failed. %v", err))
	}

	err = l.loadModel(db, loadDB, reqParams, conf)
	if err != nil {
		return nil, response.NewStorageError(response.StorageErrorCode, err)
	}

	var resp []interface{}
	if conf.TransObjToRespFunc != nil {
		resp = conf.TransObjToRespFunc(conf.AccessControl)
	}

	return list.NewPageResponse(int(count), reqParams.PageInfo, resp), nil
}

func (l *apiList) loadModel(db, queryDB *gorm.DB, reqParams *list.Params, conf *ListConf) error {
	baseInfo := reqParams.Queries[list.BaseInfo]
	b, _ := strconv.ParseBool(baseInfo)
	if b {
		return l.loadBaseInfo(queryDB, reqParams, conf)
	}
	return l.loadAllInfo(db, queryDB, reqParams, conf)
}

func (l *apiList) loadAllInfo(db, queryDB *gorm.DB, reqParams *list.Params, conf *ListConf) error {
	queryDB = store.OptPageDB(queryDB,
		reqParams.PageInfo.PageSize,
		reqParams.PageInfo.PageNum,
		reqParams.SortQuery,
		conf.ModelObj)

	if conf.LoadAssociationsDBFunc != nil {
		queryDB = conf.LoadAssociationsDBFunc(db, queryDB, conf.AuthControl.AccessControl)
	}

	return queryDB.Find(conf.ModelObjList).Error
}

func (conf *ListConf) getBaseInfoColumnList() []string {
	if len(conf.BaseInfoColumnList) > 0 {
		return conf.BaseInfoColumnList
	}

	return defaultBaseInfoColumnList
}

func (l *apiList) loadBaseInfo(queryDB *gorm.DB, reqParams *list.Params, conf *ListConf) error {
	db := store.OptPageDB(queryDB,
		reqParams.PageInfo.PageSize,
		reqParams.PageInfo.PageNum,
		reqParams.SortQuery,
		conf.ModelObj)

	db = db.Select(strings.Join(conf.getBaseInfoColumnList(), ","))
	return db.Find(conf.ModelObjList).Error
}
