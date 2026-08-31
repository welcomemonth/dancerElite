package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// MenuService 菜单管理服务
type MenuService struct {
	st *store.Store
}

// NewMenuService 创建菜单服务
func NewMenuService(st *store.Store) *MenuService {
	return &MenuService{st: st}
}

// List 获取所有菜单
func (s *MenuService) List(ctx context.Context) ([]model.Menu, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Order("sort, id")
	}
	return s.st.MenuRepo.FindAll(ctx, scope)
}

// Tree 获取菜单树形结构
func (s *MenuService) Tree(ctx context.Context) ([]*model.Menu, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("status = 1").Order("sort, id")
	}
	menus, err := s.st.MenuRepo.FindAll(ctx, scope)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

// Get 获取单个菜单
func (s *MenuService) Get(ctx context.Context, id int64) (*model.Menu, error) {
	return s.st.MenuRepo.GetByID(ctx, id)
}

// Create 创建菜单
func (s *MenuService) Create(ctx context.Context, menu *model.Menu) error {
	if menu.Type == 0 {
		menu.Type = 1
	}
	menu.Status = 1
	return s.st.MenuRepo.Create(ctx, menu)
}

// Update 更新菜单
func (s *MenuService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.MenuRepo.Update(ctx, id, updates)
}

// UpdateStatus 更新菜单状态
func (s *MenuService) UpdateStatus(ctx context.Context, id int64, status int) error {
	return s.st.MenuRepo.Update(ctx, id, map[string]any{"status": status})
}

// Delete 删除菜单
func (s *MenuService) Delete(ctx context.Context, id int64) error {
	// 检查是否有子菜单
	exists, _ := s.st.MenuRepo.Exists(ctx, "parent_id = ?", id)
	if exists {
		return errcode.ErrMenuHasChildren
	}
	// 删除角色菜单关联
	s.st.DB().WithContext(ctx).Where("menu_id = ?", id).Delete(&model.RoleMenu{})
	return s.st.MenuRepo.Delete(ctx, id)
}

// InitDefaultMenus 初始化默认菜单
func (s *MenuService) InitDefaultMenus(ctx context.Context) error {
	var count int64
	s.st.DB().WithContext(ctx).Model(&model.Menu{}).Count(&count)
	if count > 0 {
		return s.ensureDefaultMenuPermissions(ctx)
	}

	// 内容管理目录
	contentDir := model.Menu{ParentID: 0, Name: "content", Title: "内容管理", Path: "", Component: "", Icon: "Document", Sort: 1, Type: 1, Status: 1}
	s.st.DB().WithContext(ctx).Create(&contentDir)

	contentMenus := []model.Menu{
		{ParentID: contentDir.ID, Name: "articles", Title: "文章管理", Path: "/admin/articles", Component: "content/ArticleList", Icon: "Document", Sort: 1, Type: 2, Status: 1},
		{ParentID: contentDir.ID, Name: "columns", Title: "栏目管理", Path: "/admin/columns", Component: "content/ColumnList", Icon: "Menu", Sort: 2, Type: 2, Status: 1},
	}
	for i := range contentMenus {
		s.st.DB().WithContext(ctx).Create(&contentMenus[i])
	}

	// 业务管理目录
	bizDir := model.Menu{ParentID: 0, Name: "business", Title: "业务管理", Path: "", Component: "", Icon: "ShoppingCart", Sort: 2, Type: 1, Status: 1}
	s.st.DB().WithContext(ctx).Create(&bizDir)

	bizMenus := []model.Menu{
		{ParentID: bizDir.ID, Name: "activities", Title: "活动管理", Path: "/admin/activities", Component: "business/ActivityList", Icon: "Calendar", Sort: 1, Type: 2, Status: 1},
		{ParentID: bizDir.ID, Name: "registrations", Title: "报名管理", Path: "/admin/registrations", Component: "business/RegistrationList", Icon: "List", Sort: 2, Type: 2, Status: 1},
		{ParentID: bizDir.ID, Name: "payments", Title: "支付管理", Path: "/admin/payments", Component: "business/PaymentList", Icon: "Money", Sort: 3, Type: 2, Status: 1},
	}
	for i := range bizMenus {
		s.st.DB().WithContext(ctx).Create(&bizMenus[i])
	}

	// 系统管理目录
	sysDir := model.Menu{ParentID: 0, Name: "system", Title: "系统管理", Path: "", Component: "", Icon: "Setting", Sort: 3, Type: 1, Status: 1}
	s.st.DB().WithContext(ctx).Create(&sysDir)

	sysMenus := []model.Menu{
		{ParentID: sysDir.ID, Name: "mp-users", Title: "小程序用户", Path: "/admin/users", Component: "system/UserList", Icon: "User", Sort: 1, Type: 2, Status: 1},
		{ParentID: sysDir.ID, Name: "backend-users", Title: "用户管理", Path: "/admin/backend-users", Component: "system/BackendUserList", Icon: "UserFilled", Sort: 2, Type: 2, Status: 1},
		{ParentID: sysDir.ID, Name: "roles", Title: "角色管理", Path: "/admin/roles", Component: "system/RoleList", Icon: "Key", Sort: 3, Type: 2, Status: 1},
		{ParentID: sysDir.ID, Name: "menus", Title: "菜单管理", Path: "/admin/menus", Component: "system/MenuList", Icon: "Grid", Sort: 4, Type: 2, Status: 1},
		{ParentID: sysDir.ID, Name: "operation-logs", Title: "操作日志", Path: "/admin/operation-logs", Component: "system/OperationLogList", Icon: "Notebook", Sort: 5, Type: 2, Status: 1},
		{ParentID: sysDir.ID, Name: "system-configs", Title: "系统配置", Path: "/admin/system-configs", Component: "system/SystemConfigList", Icon: "Tools", Sort: 6, Type: 2, Status: 1},
	}
	for i := range sysMenus {
		s.st.DB().WithContext(ctx).Create(&sysMenus[i])
	}

	// 代码生成器
	codegenDir := model.Menu{ParentID: 0, Name: "codegen", Title: "代码生成", Path: "", Component: "", Icon: "Cpu", Sort: 4, Type: 1, Status: 1}
	s.st.DB().WithContext(ctx).Create(&codegenDir)

	codegenMenu := model.Menu{ParentID: codegenDir.ID, Name: "codegen-config", Title: "生成配置", Path: "/admin/codegen", Component: "codegen/CodegenList", Icon: "Cpu", Sort: 1, Type: 2, Status: 1}
	s.st.DB().WithContext(ctx).Create(&codegenMenu)

	return s.ensureDefaultMenuPermissions(ctx)
}

type defaultPermission struct {
	Action string
	Title  string
	Method string
	Sort   int
}

func (s *MenuService) ensureDefaultMenuPermissions(ctx context.Context) error {
	permissions := map[string][]defaultPermission{
		"articles": {
			{Action: "list", Title: "查看文章列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看文章详情", Method: "GET", Sort: 2},
			{Action: "create", Title: "新增文章", Method: "POST", Sort: 3},
			{Action: "update", Title: "编辑文章", Method: "PUT", Sort: 4},
			{Action: "update_status", Title: "更新文章状态", Method: "PUT", Sort: 5},
			{Action: "delete", Title: "删除文章", Method: "DELETE", Sort: 6},
		},
		"columns": {
			{Action: "list", Title: "查看栏目列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看栏目详情", Method: "GET", Sort: 2},
			{Action: "create", Title: "新增栏目", Method: "POST", Sort: 3},
			{Action: "update", Title: "编辑栏目", Method: "PUT", Sort: 4},
			{Action: "delete", Title: "删除栏目", Method: "DELETE", Sort: 5},
		},
		"mp-users": {
			{Action: "list", Title: "查看用户列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看用户详情", Method: "GET", Sort: 2},
			{Action: "update", Title: "编辑用户", Method: "PUT", Sort: 3},
			{Action: "delete", Title: "删除用户", Method: "DELETE", Sort: 4},
		},
		"backend-users": {
			{Action: "list", Title: "查看后台用户列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看后台用户详情", Method: "GET", Sort: 2},
			{Action: "create", Title: "新增后台用户", Method: "POST", Sort: 3},
			{Action: "update", Title: "编辑后台用户", Method: "PUT", Sort: 4},
			{Action: "update_status", Title: "更新后台用户状态", Method: "PUT", Sort: 5},
			{Action: "update_reset-password", Title: "重置后台用户密码", Method: "PUT", Sort: 6},
			{Action: "delete", Title: "删除后台用户", Method: "DELETE", Sort: 7},
		},
		"roles": {
			{Action: "list", Title: "查看角色列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看角色详情", Method: "GET", Sort: 2},
			{Action: "menus", Title: "查看角色权限", Method: "GET", Sort: 3},
			{Action: "create", Title: "新增角色", Method: "POST", Sort: 4},
			{Action: "update", Title: "编辑角色", Method: "PUT", Sort: 5},
			{Action: "update_status", Title: "更新角色状态", Method: "PUT", Sort: 6},
			{Action: "update_menus", Title: "编辑角色权限", Method: "PUT", Sort: 7},
			{Action: "delete", Title: "删除角色", Method: "DELETE", Sort: 8},
		},
		"menus": {
			{Action: "list", Title: "查看菜单列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看菜单详情", Method: "GET", Sort: 2},
			{Action: "create", Title: "新增菜单", Method: "POST", Sort: 3},
			{Action: "update", Title: "编辑菜单", Method: "PUT", Sort: 4},
			{Action: "update_status", Title: "更新菜单状态", Method: "PUT", Sort: 5},
			{Action: "delete", Title: "删除菜单", Method: "DELETE", Sort: 6},
		},
		"operation-logs": {
			{Action: "list", Title: "查看操作日志", Method: "GET", Sort: 1},
		},
		"system-configs": {
			{Action: "list", Title: "查看系统配置", Method: "GET", Sort: 1},
			{Action: "create", Title: "新增系统配置", Method: "POST", Sort: 2},
			{Action: "delete", Title: "删除系统配置", Method: "DELETE", Sort: 3},
		},
		"activities": {
			{Action: "list", Title: "查看活动列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看活动详情", Method: "GET", Sort: 2},
			{Action: "create", Title: "新增活动", Method: "POST", Sort: 3},
			{Action: "update", Title: "编辑活动", Method: "PUT", Sort: 4},
			{Action: "update_status", Title: "更新活动状态", Method: "PUT", Sort: 5},
			{Action: "delete", Title: "删除活动", Method: "DELETE", Sort: 6},
		},
		"registrations": {
			{Action: "list", Title: "查看报名列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看报名详情", Method: "GET", Sort: 2},
		},
		"payments": {
			{Action: "list", Title: "查看支付列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看支付详情", Method: "GET", Sort: 2},
			{Action: "update_refund", Title: "退款", Method: "PUT", Sort: 3},
		},
		"codegen-config": {
			{Action: "list", Title: "查看生成配置列表", Method: "GET", Sort: 1},
			{Action: "get", Title: "查看生成配置详情", Method: "GET", Sort: 2},
			{Action: "preview", Title: "预览生成配置", Method: "GET", Sort: 3},
			{Action: "create", Title: "新建生成配置", Method: "POST", Sort: 4},
			{Action: "update", Title: "编辑生成配置", Method: "PUT", Sort: 5},
			{Action: "update_generate", Title: "生成代码", Method: "POST", Sort: 6},
			{Action: "delete", Title: "删除生成配置", Method: "DELETE", Sort: 7},
		},
	}

	for menuName, items := range permissions {
		var parent model.Menu
		if err := s.st.DB().WithContext(ctx).Where("name = ? AND type = 2", menuName).First(&parent).Error; err != nil {
			continue
		}

		module := permissionModule(menuName)
		for _, item := range items {
			permission := module + ":" + item.Action
			var count int64
			s.st.DB().WithContext(ctx).Model(&model.Menu{}).
				Where("parent_id = ? AND permission = ? AND type = 3", parent.ID, permission).
				Count(&count)
			if count > 0 {
				continue
			}

			menu := model.Menu{
				ParentID:   parent.ID,
				Name:       parent.Name + "-" + item.Action,
				Title:      item.Title,
				Sort:       item.Sort,
				Type:       3,
				Permission: permission,
				Method:     item.Method,
				Status:     1,
			}
			if err := s.st.DB().WithContext(ctx).Create(&menu).Error; err != nil {
				return err
			}
		}
	}

	return s.assignAllMenusToAdmin(ctx)
}

func permissionModule(menuName string) string {
	modules := map[string]string{
		"mp-users":       "user",
		"backend-users":  "backend_user",
		"operation-logs": "operation_log",
		"system-configs": "system_config",
		"codegen-config": "codegen",
	}
	if module, ok := modules[menuName]; ok {
		return module
	}
	if len(menuName) > 3 && menuName[len(menuName)-3:] == "ies" {
		return menuName[:len(menuName)-3] + "y"
	}
	if len(menuName) > 1 && menuName[len(menuName)-1:] == "s" {
		return menuName[:len(menuName)-1]
	}
	return menuName
}

func (s *MenuService) assignAllMenusToAdmin(ctx context.Context) error {
	var adminRole model.Role
	if err := s.st.DB().WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return nil
	}

	var allMenus []model.Menu
	if err := s.st.DB().WithContext(ctx).Find(&allMenus).Error; err != nil {
		return err
	}

	for _, menu := range allMenus {
		var count int64
		s.st.DB().WithContext(ctx).Model(&model.RoleMenu{}).
			Where("role_id = ? AND menu_id = ?", adminRole.ID, menu.ID).
			Count(&count)
		if count > 0 {
			continue
		}
		if err := s.st.DB().WithContext(ctx).Create(&model.RoleMenu{RoleID: adminRole.ID, MenuID: menu.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}
