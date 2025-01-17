package v1

import (
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

type CreateSyncInstanceTaskReqV1 struct {
	Source               string `json:"source" form:"source" validate:"required" example:"actiontech-dmp"`
	Version              string `json:"version" form:"version" validate:"required" example:"5.23.01.0"`
	URL                  string `json:"url" form:"url" validate:"required" example:"http://10.186.62.56:10000"`
	DbType               string `json:"db_type" form:"db_type" validate:"required" example:"mysql"`
	GlobalRuleTemplate   string `json:"global_rule_template" form:"global_rule_template" validate:"required" example:"default_mysql"`
	SyncInstanceInterval string `json:"sync_instance_interval" form:"sync_instance_interval" validate:"required" example:"0 0 * * *"`
}

// CreateSyncInstanceTask create sync instance task
// @Summary 创建同步实例任务
// @Description create sync instance task
// @Id createSyncInstanceTaskV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @Accept json
// @Param sync_task body v1.CreateSyncInstanceTaskReqV1 true "create sync instance task request"
// @Success 200 {object} controller.BaseRes
// @router /v1/sync_instances [post]
func CreateSyncInstanceTask(c echo.Context) error {
	return createSyncInstanceTask(c)
}

type UpdateSyncInstanceTaskReqV1 struct {
	Version              *string `json:"version" form:"version" example:"5.23.01.0"`
	URL                  *string `json:"url" form:"url" example:"http://10.186.62.56:10000"`
	GlobalRuleTemplate   *string `json:"global_rule_template" form:"global_rule_template" example:"default_mysql"`
	SyncInstanceInterval *string `json:"sync_instance_interval" form:"sync_instance_interval" example:"0 0 * * *"`
}

// UpdateSyncInstanceTask update sync instance task
// @Summary 更新同步实例任务
// @Description update sync instance task
// @Id updateSyncInstanceTaskV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @Param task_id path string true "sync task id"
// @param sync_task body v1.UpdateSyncInstanceTaskReqV1 true "update sync instance request"
// @Success 200 {object} controller.BaseRes
// @router /v1/sync_instances/{task_id}/ [patch]
func UpdateSyncInstanceTask(c echo.Context) error {
	return updateSyncInstanceTask(c)
}

// DeleteSyncInstanceTask delete sync instance task
// @Summary 删除同步实例任务
// @Description delete sync instance task
// @Id deleteSyncInstanceTaskV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @param task_id path string true "sync task id"
// @Success 200 {object} controller.BaseRes
// @router /v1/sync_instances/{task_id}/ [delete]
func DeleteSyncInstanceTask(c echo.Context) error {
	return deleteSyncInstanceTask(c)
}

type SyncInstanceResV1 struct {
	IsSyncInstanceSuccess bool   `json:"is_sync_instance_success" example:"true"`
	SyncErrorMessage      string `json:"sync_error_message"`
}

type TriggerSyncInstanceResV1 struct {
	controller.BaseRes
	Data SyncInstanceResV1 `json:"data"`
}

// TriggerSyncInstance trigger sync instance
// @Summary 触发同步实例
// @Description trigger sync instance
// @Id triggerSyncInstanceV1
// @Tags sync_instance
// @Security ApiKeyAuth
// @param task_id path string true "sync task id"
// @Success 200 {object} v1.TriggerSyncInstanceResV1
// @router /v1/sync_instances/{task_id}/trigger [post]
func TriggerSyncInstance(c echo.Context) error {
	return triggerSyncInstance(c)
}

type GetSyncInstanceTaskListResV1 struct {
	controller.BaseRes
	Data []InstanceTaskResV1 `json:"data"`
}

type InstanceTaskResV1 struct {
	ID                  int        `json:"id" example:"1"`
	Source              string     `json:"source" example:"actiontech-dmp"`
	Version             string     `json:"version" example:"1.23.1"`
	URL                 string     `json:"url" example:"http://10.186.62.56:10000"`
	DbType              string     `json:"db_type" example:"mysql"`
	LastSyncStatus      string     `json:"last_sync_status" enums:"succeeded,failed" example:"succeeded"`
	LastSyncSuccessTime *time.Time `json:"last_sync_success_time" example:"2021-08-12 12:00:00"`
}

// GetSyncInstanceTaskList get sync instance task list
// @Summary 获取同步实例任务列表
// @Description get sync instance task list
// @Id GetSyncInstanceTaskList
// @Tags sync_instance
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSyncInstanceTaskListResV1
// @router /v1/sync_instances [get]
func GetSyncInstanceTaskList(c echo.Context) error {
	return getSyncInstanceTaskList(c)
}

type SyncTaskTipsResV1 struct {
	Source  string   `json:"source" example:"actiontech-dmp"`
	DbTypes []string `json:"db_types"`
}

type GetSyncTaskSourceTipsResV1 struct {
	controller.BaseRes
	Data []SyncTaskTipsResV1 `json:"data"`
}

// GetSyncTaskSourceTips get sync instance tips
// @Summary 获取同步任务提示
// @Description get sync task tips
// @Id GetSyncTaskSourceTips
// @Tags sync_instance
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSyncTaskSourceTipsResV1
// @router /v1/sync_instances/source_tips [get]
func GetSyncTaskSourceTips(c echo.Context) error {
	return getSyncTaskSourceTips(c)
}

type InstanceTaskDetailResV1 struct {
	ID                   int    `json:"id" example:"1"`
	Source               string `json:"source" example:"actiontech-dmp"`
	Version              string `json:"version" example:"1.23.1"`
	URL                  string `json:"url" example:"http://10.186.62.56:10000"`
	DbType               string `json:"db_type" example:"mysql"`
	RuleTemplate         string `json:"rule_template" example:"default_mysql"`
	SyncInstanceInterval string `json:"sync_instance_interval" example:"0 0 * * *"`
}

type GetSyncInstanceTaskResV1 struct {
	controller.BaseRes
	Data InstanceTaskDetailResV1 `json:"data"`
}

// GetSyncInstanceTask get sync task detail
// @Summary 获取同步任务详情
// @Description get sync task detail
// @Id GetSyncInstanceTask
// @Tags sync_instance
// @Security ApiKeyAuth
// @param task_id path string true "sync task id"
// @Success 200 {object} v1.GetSyncInstanceTaskResV1
// @router /v1/sync_instances/{task_id}/ [get]
func GetSyncInstanceTask(c echo.Context) error {
	return getSyncInstanceTask(c)
}
