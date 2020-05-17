package feedly

// NewBool returns a pointer to the b bool.
func NewBool(b bool) *bool {
	return &b
}

// NewEmbeddedFilter returns a pointer to the e EmbeddedFilter.
func NewEmbeddedFilter(e EmbeddedFilter) *EmbeddedFilter {
	return &e
}

// NewEngagementFilter returns a pointer to the e EngagementFilter.
func NewEngagementFilter(e EngagementFilter) *EngagementFilter {
	return &e
}

// NewInt returns a pointer to the i int.
func NewInt(i int) *int {
	return &i
}

// NewMarkAction returns a pointer to the a MarkAction.
func NewMarkAction(m MarkAction) *MarkAction {
	return &m
}

// NewMarkType returns a pointer to the a MarkType.
func NewMarkType(m MarkType) *MarkType {
	return &m
}

// NewContentRank returns a pointer to the r ContentRank.
func NewContentRank(r ContentRank) *ContentRank {
	return &r
}

// NewString returns a pointer to the s string.
func NewString(s string) *string {
	return &s
}
