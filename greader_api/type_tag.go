package greader_api

// struct: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPITag.swift

/*
{
  "tags": [
    {
      "id": "user/-/state/com.google/starred"
    },
    {
      "id": "user/-/label/未分类",
      "type": "folder"
    }
  ]
}
*/

type Tag struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"` // folder
}

type tagList struct {
	Tags []*Tag `json:"tags,omitempty"`
}

const (
	tagDefaultStarred = "user/-/state/com.google/starred"
	tagTypeFolder     = "folder"
)

func buildTads(tagNames []string) []*Tag {
	var tags []*Tag
	tags = append(tags, &Tag{
		ID: tagDefaultStarred,
	})

	for _, tagName := range tagNames {
		tags = append(tags, &Tag{
			ID:   "user/-/label/" + tagName,
			Type: tagTypeFolder,
		})
	}
	return tags
}
