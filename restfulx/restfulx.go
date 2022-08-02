package restfulx

import (
	"encoding/json"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/logger"
	"github.com/XM-GO/PandaKit/model"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"strconv"
)

// 绑定并校验请求结构体参数  结构体添加 例如： binding:"required" 或binding:"required,gt=10"
func BindJsonAndValid(request *restful.Request, data any) {
	if err := request.ReadEntity(data); err != nil {
		panic(any(biz.NewBizErr("传参格式错误：" + err.Error())))
	}
}

// 绑定查询字符串到
func BindQuery(request *restful.Request, data any) {
	if err := request.ReadEntity(data); err != nil {
		panic(any(biz.NewBizErr(err.Error())))
	}
}
func ParamsToAny(request *restful.Request, in any) {
	vars := make(map[string]any)
	for k, v := range request.PathParameters() {
		vars[k] = v
	}
	marshal, _ := json.Marshal(vars)
	err := json.Unmarshal(marshal, in)
	biz.ErrIsNil(err, "error get path value encoding unmarshal")
	return
}

// 获取分页参数
func GetPageParam(request *restful.Request) *model.PageParam {
	return &model.PageParam{PageNum: QueryInt(request, "pageNum", 1), PageSize: QueryInt(request, "pageSize", 10)}
}

// 获取查询参数中指定参数值，并转为int
func QueryInt(request *restful.Request, qm string, defaultInt int) int {
	qv := request.QueryParameter(qm)
	if qv == "" {
		return defaultInt
	}
	qvi, err := strconv.Atoi(qv)
	biz.ErrIsNil(err, "query param not int")
	return qvi
}

// 获取路径参数
func PathParamInt(request *restful.Request, pm string) int {
	value, _ := strconv.Atoi(request.PathParameter(pm))
	return value
}

// 文件下载
func Download(req *restful.Request, resp *restful.Response, filename string) {
	resp.Header().Add("success", "true")
	resp.Header().Set("Content-Length", "-1")
	resp.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(
		resp.ResponseWriter,
		req.Request, filename)
}

// 返回统一成功结果
func SuccessRes(response *restful.Response, data any) {
	response.WriteEntity(model.Success(data))
}

// 返回失败结果集
func ErrorRes(response *restful.Response, err any) {
	switch t := err.(type) {
	case *biz.BizError:
		response.WriteEntity(model.Error(t))
		break
	case error:
		response.WriteEntity(model.ServerError())
		logger.Log.Error(t)
		break
	case string:
		response.WriteEntity(model.ServerError())
		logger.Log.Error(t)
		break
	default:
		logger.Log.Error(t)
	}
}
