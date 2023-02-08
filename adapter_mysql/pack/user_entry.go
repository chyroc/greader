package pack

import "github.com/chyroc/greader/adapter_mysql/dal"

func UserEntryToRelation(userIDs []int64, entryPOs []*dal.ModelEntry) []*dal.ModeUserEntryRelation {
	pos := []*dal.ModeUserEntryRelation{}
	for _, v := range entryPOs {
		for _, userID := range userIDs {
			pos = append(pos, &dal.ModeUserEntryRelation{
				BaseModel: dal.BaseModel{},
				UserID:    userID,
				FeedID:    v.FeedID,
				EntryID:   v.ID,
				Readed:    false,
				Starred:   false,
			})
		}
	}
	return pos
}
