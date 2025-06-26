package constant

type CacheKey string

const (
	DepartmentTreeSfgKey         CacheKey = "department_cache:tree"
	DepartmentTreeWithCrewSfgKey CacheKey = "department_cache:tree_with_crew"
	DepartmentTreeRefreshSfgKey  CacheKey = "department_cache:tree:refresh"
	DepartmentTreeCleanSfgKey    CacheKey = "department_cache:tree:clean"
)
