package enums

// Operation defines the operation of a diff item.
type Operation int8

//go:generate stringer -type=Operation -trimprefix=Diff

const (
	//删除
	DiffDelete Operation = -1
	//新增
	DiffInsert Operation = 1
	//相等
	DiffEqual Operation = 0
	//修改
	DiffModify Operation = 2
)

type IndexType string

const (
	//主键
	IndexPK IndexType = "pk"
	//唯一
	IndexUnq IndexType = "unq"
	//普通
	IndexNor IndexType = "idx"
)
