package ginx

import (
	"encoding/json"
	"github.com/PandaXGO/PandaKit/biz"
	"github.com/PandaXGO/PandaKit/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 绑定并校验请求结构体参数  结构体添加 例如： binding:"required" 或binding:"required,gt=10"
func BindJsonAndValid(rc *ReqCtx, data any) {
	if err := rc.GinCtx.ShouldBindJSON(data); err != nil {
		panic(any(biz.NewBizErr("传参格式错误：" + err.Error())))
	}
	if err := rc.Validate.Struct(data); err != nil {
		panic(any(biz.CodeInvalidParameter))
	}
}

// 绑定查询字符串到
func BindQuery(rc *ReqCtx, data any) {
	if err := rc.GinCtx.ShouldBindQuery(data); err != nil {
		panic(any(biz.NewBizErr(err.Error())))
	}
}

func ParamsToAny(rc *ReqCtx, in any) {
	vars := make(map[string]any)
	for _, v := range rc.GinCtx.Params {
		vars[v.Key] = v.Value
	}
	marshal, _ := json.Marshal(vars)
	err := json.Unmarshal(marshal, in)
	biz.ErrIsNil(err, "error get path value encoding unmarshal")
	return
}

// 获取分页参数
func GetPageParam(rc *ReqCtx) *model.PageParam {
	return &model.PageParam{PageNum: QueryInt(rc, "pageNum", 1), PageSize: QueryInt(rc, "pageSize", 10)}
}

// 获取查询参数中指定参数值，并转为int
func QueryInt(rc *ReqCtx, qm string, defaultInt int) int {
	qv := rc.GinCtx.Query(qm)
	if qv == "" {
		return defaultInt
	}
	qvi, err := strconv.Atoi(qv)
	biz.ErrIsNil(err, "query param not int")
	return qvi
}

func QueryParam(rc *ReqCtx, key string) string {
	return rc.GinCtx.Query(key)
}

// 获取路径参数
func PathParamInt(rc *ReqCtx, pm string) int {
	value, _ := strconv.Atoi(rc.GinCtx.Param(pm))
	return value
}

func PathParam(rc *ReqCtx, pm string) string {
	return rc.GinCtx.Param(pm)
}

// 文件下载
func Download(rc *ReqCtx, filename string) {
	rc.GinCtx.Writer.Header().Add("success", "true")
	rc.GinCtx.Writer.Header().Set("Content-Length", "-1")
	rc.GinCtx.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	rc.GinCtx.File(filename)
	http.ServeFile(
		rc.GinCtx.Writer,
		rc.GinCtx.Request, filename)
}

// 返回统一成功结果
func SuccessRes(g *gin.Context, data any) {
	g.JSON(http.StatusOK, model.Success(data))
}

// 返回失败结果集
func ErrorRes(g *gin.Context, err any) {
	if err != nil {

	}
	switch t := err.(type) {
	case *biz.BizError:
		g.JSON(http.StatusOK, model.Error(t))
		break
	case error:
		g.JSON(http.StatusOK, model.ServerError())
		// panic(err)
		break
	case string:
		g.JSON(http.StatusOK, model.ServerError())
		// panic(err)
		break
	default:
	}
}
