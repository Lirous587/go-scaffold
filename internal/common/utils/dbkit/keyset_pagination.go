package dbkit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type keysetCursor struct {
	CreatedAt time.Time
	ID        string
}

type CursorFields interface {
	GetCreatedAt() time.Time
	GetID() string
}

type keyset[T CursorFields] struct {
	IDCol        string
	CreatedAtCol string
	PrevCursor   string
	NextCursor   string
	PageSize     int
}

type paginationResult[T any] struct {
	Items      []T
	PrevCursor string
	NextCursor string
	HasPrev    bool
	HasNext    bool
}

func NewKeyset[T CursorFields](
	idCol string,
	createdAtCol string,
	prevCursor string,
	nextCursor string,
	pageSize int,
) *keyset[T] {
	if pageSize <= 0 {
		pageSize = 5
	}
	return &keyset[T]{
		IDCol:        idCol,
		CreatedAtCol: createdAtCol,
		PrevCursor:   prevCursor,
		NextCursor:   nextCursor,
		PageSize:     pageSize,
	}
}

func (k keyset[T]) encode(cursor *keysetCursor) string {
	if cursor == nil {
		return ""
	}

	data, _ := json.Marshal(cursor)
	return base64.StdEncoding.EncodeToString(data)
}

func (k keyset[T]) decodeCursor(cursorStr string) *keysetCursor {
	if cursorStr == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil
	}
	var payload keysetCursor
	if json.Unmarshal(data, &payload) != nil {
		return nil
	}
	return &payload
}

func (k keyset[T]) extractKeysetCursor(item T) *keysetCursor {
	return &keysetCursor{
		CreatedAt: item.GetCreatedAt(),
		ID:        item.GetID(),
	}
}

// ApplyKeysetMods 附加keyset查询的mods
// mods建议额外预留4
func (k keyset[T]) ApplyKeysetMods(base []qm.QueryMod) []qm.QueryMod {
	var cursor *keysetCursor
	var order string

	if k.NextCursor != "" {
		cursor = k.decodeCursor(k.NextCursor)
		if cursor != nil {
			base = append(base, qm.Where(fmt.Sprintf("(%s, %s) > (?, ?)", k.CreatedAtCol, k.IDCol), cursor.CreatedAt, cursor.ID))
		}
		order = fmt.Sprintf("%s ASC, %s ASC", k.CreatedAtCol, k.IDCol)
	} else if k.PrevCursor != "" {
		cursor = k.decodeCursor(k.PrevCursor)
		if cursor != nil {
			base = append(base, qm.Where(fmt.Sprintf("(%s, %s) < (?, ?)", k.CreatedAtCol, k.IDCol), cursor.CreatedAt, cursor.ID))
		}
		order = fmt.Sprintf("%s DESC, %s DESC", k.CreatedAtCol, k.IDCol)
	} else {
		order = fmt.Sprintf("%s ASC, %s ASC", k.CreatedAtCol, k.IDCol)
	}

	base = append(base, qm.OrderBy(order))
	base = append(base, qm.Limit(k.PageSize+1))

	return base
}

func (k keyset[T]) BuildPaginationResult(domainSlice []T) *paginationResult[T] {
	hasMore := len(domainSlice) > k.PageSize

	// 截取
	if hasMore {
		domainSlice = domainSlice[:k.PageSize]
	}

	// 判断分页方向
	isPrev := k.PrevCursor != ""
	isNext := k.NextCursor != ""

	// 如果是 Prev 分页，先按 DESC 取，结果需要反转为 ASC 展示
	if isPrev && len(domainSlice) > 0 {
		slices.Reverse(domainSlice)
	}

	// 生成新的游标
	var prevCursor, nextCursor string
	if len(domainSlice) > 0 {
		// 首条和末条数据生成游标
		first := domainSlice[0]
		last := domainSlice[len(domainSlice)-1]

		// 编码
		prevCursor = k.encode(k.extractKeysetCursor(first))
		nextCursor = k.encode(k.extractKeysetCursor(last))
	}

	// hasPrev/hasNext 判断
	var hasPrev, hasNext bool
	switch {
	case isNext:
		hasPrev = true
		hasNext = hasMore
	case isPrev:
		hasPrev = hasMore
		hasNext = true
	default:
		hasPrev = false
		hasNext = hasMore
	}

	// 没有下一页时 next_cursor置空
	if !hasNext {
		nextCursor = ""
	}
	// 没有上一页时 prev_cursor置空
	if !hasPrev {
		prevCursor = ""
	}

	return &paginationResult[T]{
		Items:      domainSlice,
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
	}
}
