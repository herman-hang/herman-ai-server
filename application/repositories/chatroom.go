package repositories

import (
	"fmt"
	"github.com/herman-hang/herman/application/constants"
	"github.com/herman-hang/herman/application/models"
	"github.com/herman-hang/herman/kernel/core"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// ChatroomRepository 聊天室表仓储层
type ChatroomRepository struct {
	BaseRepository
}

// Chatroom 实例化聊天室表仓储层
// @param *gorm.DB tx 事务
// @return AdminRepository 返回聊天室表仓储层
func Chatroom(tx ...*gorm.DB) *ChatroomRepository {
	if len(tx) > 0 && tx[0] != nil {
		return &ChatroomRepository{BaseRepository{Model: new(models.Chatroom), Db: tx[0]}}
	}

	return &ChatroomRepository{BaseRepository{Model: new(models.Chatroom), Db: core.Db()}}
}

// Get 查询聊天室表
// @param []uint ids 聊天室ID集合
// @param map[string]interface{} data 查询条件
// @return info err 返回数据和一个错误
func (base ChatroomRepository) Get(ids []uint, data map[string]interface{}) (info map[string]interface{}, err error) {
	var (
		page    PageInfo
		total   int64
		pageNum int64
		list    []map[string]interface{}
	)
	if len(data) > 0 {
		if err := mapstructure.WeakDecode(data, &page); err != nil {
			panic(constants.MapToStruct)
		}
	}
	// 总条数
	base.Db.Model(&models.Chatroom{}).Count(&total)
	// 计算总页数
	if page.PageSize != 0 && total%page.PageSize != 0 {
		pageNum = total / page.PageSize
		pageNum++
	}
	err = base.Db.Model(&models.Chatroom{}).
		Select([]string{"id", "name", "created_at"}).
		Where("id IN ?", ids).
		Where("ai_type", data["aiType"]).
		Where("name like ?", fmt.Sprintf("%%%s%%", data["keywords"])).
		Limit(int(page.PageSize)).
		Offset(int((page.Page - 1) * page.PageSize)).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	data = map[string]interface{}{
		"list":     list,          // 数据
		"total":    total,         // 总条数
		"pageNum":  pageNum,       // 总页数
		"pageSize": page.PageSize, // 每页大小
		"page":     page.Page,     // 当前页码
	}
	return data, nil
}
