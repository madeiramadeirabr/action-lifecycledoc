package jsonc

type CommentValue struct {
	comment string
	value   interface{}
}

func NewCommentValue(comment string, value interface{}) *CommentValue {
	return &CommentValue{
		comment: comment,
		value:   value,
	}
}

func (c *CommentValue) GetComment() string {
	return c.comment
}

func (c *CommentValue) GetValue() interface{} {
	return c.value
}
